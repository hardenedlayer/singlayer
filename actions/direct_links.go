package actions

import (
	"errors"

	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

type DirectLinksResource struct {
	buffalo.Resource
}

type Reply struct {
	Reply string
}

func (v DirectLinksResource) List(c buffalo.Context) error {
	dlinks := &models.DirectLinks{}
	if c.Session().Get("is_admin").(bool) {
		tx := c.Value("tx").(*pop.Connection)
		err := tx.Order("created_at desc").All(dlinks)
		if err != nil {
			return err
		}
	} else {
		ticks := &models.DirectLinks{}
		single := getCurrentSingle(c)
		actor := c.Value("actor").(string)
		if actor == "All" {
			c.Logger().Debugf("multi mode for single!")
			ticks = single.MyDirectLinks()
		} else {
			user := single.UserByUsername(c.Value("actor"))
			if user == nil {
				c.Logger().Errorf("SECURITY: cannot found user for %v", actor)
				c.Flash().Add("warning", "Oops! Who are you?")
			} else {
				c.Logger().Debugf("single mode for %v.", actor)
				ticks = user.DirectLinks()
			}
		}
		if ticks == nil {
			c.Flash().Add("danger", "Oops! cannot search on directlinks!")
		} else {
			dlinks = ticks
		}
	}
	c.Set("dlinks", dlinks)
	return c.Render(200, r.HTML("direct_links/index.html"))
}

func (v DirectLinksResource) Show(c buffalo.Context) error {
	dlink, err := setDirectLink(c)
	if err != nil {
		return err
	}

	user := getCurrentSingle(c).UserByAccount(dlink.AccountId)
	if c.Session().Get("is_admin").(bool) {
		user, _ = models.FindUser(dlink.UserId)
	}
	models.SyncTickets(user)
	ticket, err := models.FindTicket(dlink.TicketId)
	if err == nil {
		ticket.SyncTicketUpdates(user)
	} else {
		c.Logger().Errorf("cannot pick related ticket: %v", err)
	}

	c.Set("statuses", []string{
		"note",
		"accepted",
		"configured",
		"confirmed",
		"canceled",
	})

	c.Set("vlan", models.VLAN(dlink.VlanId))
	c.Set("dlink", dlink)
	c.Set("progresses", dlink.Progresses())
	c.Set("updates", dlink.Updates())
	return c.Render(200, r.HTML("direct_links/show.html"))
}

func (v DirectLinksResource) Add(c buffalo.Context) error {
	plink, err := setDirectLink(c)
	if err != nil {
		c.Logger().Errorf("oops! cannot get previous link: %v", err)
	}
	c.Logger().Debugf("add redundancy link for %v", plink)
	c.Set("plink", plink)
	return v.New(c)
}

func (v DirectLinksResource) New(c buffalo.Context) error {
	actor := c.Value("actor").(string)
	if actor == "All" {
		c.Flash().Add("danger", "Please select actor(account) before order.")
		return c.Redirect(302, "/directlinks")
	}
	user := getCurrentSingle(c).UserByUsername(c.Value("actor"))
	c.Set("order_for", user.CompanyName)
	account := user.Account()

	plink, _ := c.Value("plink").(*models.DirectLink)
	c.Logger().Debugf("previous link: %v", plink)

	dlink := &models.DirectLink{}
	dlink.UserId = user.ID
	dlink.AccountId = user.AccountId
	dlink.Location = "SEO01"

	// default values
	dlink.Type = "Cloud"
	dlink.Speed = 1
	dlink.RoutingOption = "Local"
	dlink.Prefix = 31
	dlink.Migration = "ANYDAY 00:00 ~ 06:00 KST"

	// auto assignment
	dlink.LineNumber = len(*account.DirectLinks()) + 1
	if plink != nil {
		dlink.MultiPath = true
		dlink.SiblingID = plink.ID
		dlink.VlanId = plink.VlanId
		dlink.Router = (plink.Router + 2) % 2 + 1
		c.Set("plink", plink)
	} else {
		dlink.VlanId = models.NextVLAN(dlink.AccountId).ID
		dlink.Router = models.NextRouter()
		c.Set("plink", "No")
	}
	dlink.ASN = 4204200000 + dlink.VlanId
	dlink.Port = "N/A"

	c.Set("types", []string{"Cloud", "NSP"})
	c.Set("speeds", []int{1, 10})
	c.Set("routing_options", []string{"Local", "Global"})
	c.Set("prefixes", []int{31, 30})

	c.Set("dlink", dlink)
	return c.Render(200, r.HTML("direct_links/new.html"))
}

func (v DirectLinksResource) Create(c buffalo.Context) error {
	dlink := &models.DirectLink{}
	err := c.Bind(dlink)
	if err != nil {
		return err
	}
	single := getCurrentSingle(c)
	user := single.UserByUsername(c.Value("actor"))
	if user.AccountId != dlink.AccountId || user.ID != dlink.UserId {
		c.Logger().Errorf("SECURITY: incorrect account/user %v!=%v or %v!=%v",
			user.AccountId, dlink.AccountId, user.ID, dlink.UserId)
		return err
	}
	if dlink.Type == "Cloud" {
		dlink.Port = "N/A"
	}
	dlink.SingleID = single.ID
	dlink.Status = "draft"
	dlink.Signature = dlink.Hash()
	c.Logger().Debugf("binded dlink: %v", dlink.Marshal())

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndCreate(dlink)
	c.Logger().Debugf("saved dlink: %v", dlink.Marshal())
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("dlink", dlink)
		c.Set("errors", verrs)
		c.Logger().Printf("errors %v", verrs)
		return c.Render(422, r.HTML("direct_links/new.html"))
	}

	vlan := models.VLAN(dlink.VlanId)
	if vlan.AccountId != dlink.AccountId {
		return c.Error(404, errors.New("DirectLink Not Found"))
	}
	vlan.AccountId = dlink.AccountId
	vlan.Booked = false
	switch dlink.Router {
	case 1:
		vlan.R1LinkID = dlink.ID
	case 2:
		vlan.R2LinkID = dlink.ID
	default:
		c.Logger().Errorf("oops! why am I here???")
	}
	err = tx.Save(vlan)
	if err != nil {
		c.Logger().Errorf("oops! cannot save vlan: %v", err)
		return err
	}

	plink := models.PickDirectLink(dlink.SiblingID)
	if plink == nil {
		c.Logger().Warnf("no pair link? %v", dlink.SiblingID)
	} else {
		plink.SiblingID = dlink.ID
		plink.MultiPath = true
		err = tx.Save(plink)
		if err != nil {
			c.Logger().Errorf("oops! cannot save pair link: %v %v", plink, err)
		}
	}

	c.Flash().Add("success", "DirectLink was created successfully")
	return c.Redirect(302, "/directlinks/%s", dlink.ID)
}

func (v DirectLinksResource) Edit(c buffalo.Context) error {
	dlink, err := setDirectLink(c)
	if err != nil {
		return err
	}
	c.Set("dlink", dlink)
	return c.Render(200, r.HTML("direct_links/edit.html"))
}

func (v DirectLinksResource) Update(c buffalo.Context) error {
	dlink, err := setDirectLink(c)
	if err != nil {
		return err
	}
	err = c.Bind(dlink)
	if err != nil {
		return err
	}
	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndUpdate(dlink)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("dlink", dlink)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("direct_links/edit.html"))
	}
	c.Flash().Add("success", "DirectLink was updated successfully")
	return c.Redirect(302, "/directlinks/%s", dlink.ID)
}

func (v DirectLinksResource) Order(c buffalo.Context) error {
	dlink, err := setDirectLink(c)
	if err != nil {
		return err
	}

	single := getCurrentSingle(c)
	user := single.UserByUsername(c.Value("actor"))
	if c.Session().Get("is_admin").(bool) {
		user, _ = models.FindUser(dlink.UserId)
	}
	ticket_id, err := models.CreateDirectLinkTicket(user, dlink)
	if err != nil {
		c.Logger().Errorf("ticket creation error: %v", err)
		return err
	}
	progress := models.NewProgress(dlink.ID, "order")
	progress.SingleID = getCurrentSingle(c).ID
	if ticket, err := models.FindTicket(ticket_id); err == nil {
		progress.UpdateId = ticket.FirstUpdate().ID
	}
	progress.Save()
	c.Logger().Debugf("progress: %v", progress)

	dlink.TicketId = ticket_id
	verrs, err := c.Value("tx").(*pop.Connection).ValidateAndUpdate(dlink)
	if err != nil {
		c.Logger().Errorf("database error: %v", err)
		return err
	}
	if verrs.HasAny() {
		c.Logger().Errorf("validation error: %v", verrs)
		return verrs
	}
	c.Flash().Add("success", "DirectLink was ordered successfully")
	return c.Redirect(302, "/directlinks/%s", dlink.ID)
}

func (v DirectLinksResource) Proceed(c buffalo.Context) error {
	single := getCurrentSingle(c)

	dlink, err := setDirectLink(c)
	if err != nil {
		return err
	}
	progress := models.NewProgress(dlink.ID, "")
	err = c.Bind(progress)
	if err != nil {
		return err
	}
	progress.SingleID = single.ID
	c.Logger().Debugf("progress: %v", progress)
	progress.Save()

	dlink.Status = progress.Action
	verrs, err := c.Value("tx").(*pop.Connection).ValidateAndUpdate(dlink)
	if err != nil {
		c.Logger().Errorf("database error: %v", err)
		return err
	}
	if verrs.HasAny() {
		c.Logger().Errorf("validation error: %v", verrs)
		return verrs
	}

	reply := &Reply{}
	err = c.Bind(reply)
	if err == nil && len(reply.Reply) > 0 {
		user := single.UserByUsername(c.Value("actor"))
		ticket := dlink.Ticket()
		if user == nil {
			return c.Error(412, errors.New("Actor/User Not Found"))
		}
		if ticket == nil {
			return c.Error(412, errors.New("Associated Ticket Not Found"))
		}
		u, err := ticket.AddUpdate(user, reply.Reply)
		c.Logger().Debugf("new update %v on %v created!", u.ID, u.TicketId)
		if err != nil {
			c.Logger().Errorf("cannot add an update: %v", err)
			return err
		}
	} else {
		c.Logger().Debugf("blank reply, no-reply mode: '%v'", reply.Reply)
	}
	c.Flash().Add("success", "DirectLink was ordered successfully")
	return c.Redirect(302, "/directlinks/%s", dlink.ID)
}

// ADMIN PROTECTED
func (v DirectLinksResource) Destroy(c buffalo.Context) error {
	dlink, err := setDirectLink(c)
	if err != nil {
		return err
	}
	tx := c.Value("tx").(*pop.Connection)
	err = tx.Destroy(dlink)
	if err != nil {
		return err
	}
	c.Flash().Add("success", "DirectLink was destroyed successfully")
	return c.Redirect(302, "/directlinks")
}

//// action helpers

// Find target user based on the context and permission.
func setDirectLink(c buffalo.Context) (dlink *models.DirectLink, err error) {
	tx := c.Value("tx").(*pop.Connection)
	dlink = &models.DirectLink{}
	if c.Session().Get("is_admin").(bool) {
		err = tx.Find(dlink, c.Param("directlink_id"))
	} else {
		single := getCurrentSingle(c)
		err = tx.Where("single_id=?", single.ID).
			Find(dlink, c.Param("directlink_id"))
		if err != nil {
			err = c.Error(404, errors.New("DirectLink Not Found"))
		}
	}
	c.Logger().Debugf("setDirectLink() returns dlink: %v, err: %v", dlink, err)
	return
}

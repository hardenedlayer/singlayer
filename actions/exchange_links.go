package actions

import (
	"errors"

	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

// PROTECTED BY GROUP PROTECTOR
func ExchangeLinksList(c buffalo.Context) error {
	dlinks := &models.DirectLinks{}
	q := pop.Q(c.Value("tx").(*pop.Connection)).Where("type = 'cloud'")
	for_all := c.Param("all") == "true"
	if !for_all {
		q = q.Where("status = 'configured'")
	}
	err := q.Order("created_at desc").All(dlinks)
	if err != nil {
		return err
	}
	c.Set("dlinks", dlinks)
	c.Set("theme", "dark")
	return c.Render(200, r.HTML("exchange_links/index.html"))
}

// PROTECTED BY GROUP PROTECTOR
func ExchangeLinksShow(c buffalo.Context) error {
	dlink, err := setExchangeLink(c)
	if err != nil {
		return err
	}

	c.Set("statuses", []string{
		"exchanger note",
	})

	c.Set("dlink", dlink)
	c.Set("progresses", dlink.Progresses())
	c.Set("updates", dlink.Updates())
	c.Set("theme", "dark")
	return c.Render(200, r.HTML("exchange_links/show.html"))
}

// PROTECTED BY GROUP PROTECTOR
func ExchangeLinksProceed(c buffalo.Context) error {
	dlink, err := setExchangeLink(c)
	if err != nil {
		return err
	}

	progress := models.NewProgress(dlink.ID, "")
	err = c.Bind(progress)
	if err != nil {
		return err
	}
	progress.SingleID = getCurrentSingle(c).ID
	progress.Save()
	c.Logger().Infof("add progress: %v %v", dlink.ID, progress.Action)

	c.Flash().Add("success", "Note added successfully")
	return c.Redirect(302, "/exchange/links/%s", dlink.ID)
}

// PROTECTED BY GROUP PROTECTOR
func ExchangeLinksConfirm(c buffalo.Context) error {
	dlink, err := setExchangeLink(c)
	if err != nil {
		return err
	}

	dlink.Status = "confirmed" // update directlink status first
	verrs, err := c.Value("tx").(*pop.Connection).ValidateAndUpdate(dlink)
	if err != nil {
		c.Logger().Errorf("database error: %v", err)
		return err
	}
	if verrs.HasAny() {
		c.Logger().Errorf("validation error: %v", verrs)
		return verrs
	}

	user, err := models.FindUser(dlink.UserId) // ticket update. so long...
	if err != nil {
		return c.Error(412, errors.New("Actor/User Not Found"))
	}
	ticket := dlink.Ticket()
	if ticket == nil {
		return c.Error(412, errors.New("Associated Ticket Not Found"))
	}
	u, err := ticket.AddUpdate(user, `Automated Update:
Connection configured and confirmed by network engineer.`)
	if err != nil {
		c.Logger().Errorf("cannot add an update: %v", err)
		return err
	}
	c.Logger().Infof("new update %v on %v created!", u.ID, u.TicketId)

	single := getCurrentSingle(c)
	progress := models.NewProgress(dlink.ID, "confirmed") // add a progress
	progress.SingleID = single.ID
	//progress.UpdateId = u.ID
	progress.Note = "The link was configured and confirmed by network engineer."
	progress.Save()
	c.Logger().Infof("add progress: %v %v", dlink.ID, progress.Action)

	s, err := models.FindSingle(dlink.SingleID) // finally shot a mail
	if err == nil {
		single.AdminMail(*dlink, "DLink Completed", s.Mail(), "admin", "exman")
	}

	c.Flash().Add("success", "DirectLink was confirmed successfully")
	return c.Redirect(302, "/exchange/links")
}

//// action helpers

// Find target user based on the context and permission.
func setExchangeLink(c buffalo.Context) (dlink *models.DirectLink, err error) {
	tx := c.Value("tx").(*pop.Connection)
	dlink = &models.DirectLink{}
	err = tx.Find(dlink, c.Param("directlink_id"))
	return
}

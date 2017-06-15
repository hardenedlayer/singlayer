package actions

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"

	"github.com/hardenedlayer/singlayer/models"
)

type TicketsResource struct {
	buffalo.Resource
}

func (v TicketsResource) List(c buffalo.Context) error {
	tickets := &models.Tickets{}

	pager := &pop.Paginator{}
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		page = 1
	}
	pp, err := strconv.Atoi(c.Param("pp"))
	if err != nil {
		pp = 20
	}
	if pp > 100 {
		pp = 100
	}

	if c.Session().Get("is_admin").(bool) {
		tx := c.Value("tx").(*pop.Connection)
		q := tx.Paginate(page, pp)
		err := q.Order("last_edit_date desc").All(tickets)
		pager = q.Paginator
		if err != nil {
			return err
		}
	} else {
		ticks := &models.Tickets{}
		single := getCurrentSingle(c)
		actor := c.Value("actor").(string)
		if actor == "All" {
			c.Logger().Debugf("multi mode for single!")
			ticks, pager = single.MyTickets(page, pp)
		} else {
			user := single.UserByUsername(c.Value("actor"))
			if user == nil {
				c.Logger().Errorf("SECURITY: cannot found user for %v", actor)
				c.Flash().Add("warning", "Oops! Who are you?")
			} else {
				c.Logger().Debugf("single mode for %v. full-sync", actor)
				count, err := models.SyncTickets(user)
				if err != nil {
					return err
				}
				c.Logger().Debugf("%v new tickets synced", count)
				ticks, pager = user.Tickets(page, pp)
			}
		}
		if ticks == nil {
			c.Flash().Add("danger", "Oops! cannot search on tickets!")
		} else {
			tickets = ticks
		}
	}
	addTicketHelpers(c)
	c.Set("pager", pager)
	c.Set("tickets", tickets)
	return c.Render(200, r.HTML("tickets/index.html"))
}

func (v TicketsResource) Show(c buffalo.Context) error {
	ticket, err := setTicket(c)
	if err != nil {
		return err
	}

	updates := ticket.Updates()
	if c.Session().Get("is_admin").(bool) == false {
		// do not update updates with admin priviledge.
		if ticket.TotalUpdateCount != len(*updates) {
			c.Logger().Debugf("updates mismatch. resync... has %v of %v",
				len(*updates), ticket.TotalUpdateCount)
			user := getCurrentSingle(c).UserByAccount(ticket.AccountId)
			if user == nil {
				return err
			}
			count, err := ticket.SyncTicketUpdates(user)
			if err != nil {
				return err
			}
			c.Logger().Debugf("%v new updates synced.", count)
			updates = ticket.Updates()
		}
	}

	addTicketHelpers(c)
	c.Set("ticket", ticket)
	c.Set("updates", updates)
	return c.Render(200, r.HTML("tickets/show.html"))
}

func (v TicketsResource) New(c buffalo.Context) error {
	c.Set("ticket", &models.Ticket{})
	return c.Render(200, r.HTML("tickets/new.html"))
}

func (v TicketsResource) Create(c buffalo.Context) error {
	ticket := &models.Ticket{}
	err := c.Bind(ticket)
	if err != nil {
		return err
	}
	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndCreate(ticket)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("ticket", ticket)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("tickets/new.html"))
	}
	c.Flash().Add("success", "Ticket was created successfully")
	return c.Redirect(302, "/tickets/%d", ticket.ID)
}

func (v TicketsResource) Edit(c buffalo.Context) error {
	ticket, err := setTicket(c)
	if err != nil {
		return err
	}
	c.Set("ticket", ticket)
	return c.Render(200, r.HTML("tickets/edit.html"))
}

func (v TicketsResource) Update(c buffalo.Context) error {
	ticket, err := setTicket(c)
	if err != nil {
		return err
	}
	err = c.Bind(ticket)
	if err != nil {
		return err
	}
	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndUpdate(ticket)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("ticket", ticket)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("tickets/edit.html"))
	}
	c.Flash().Add("success", "Ticket was updated successfully")
	return c.Redirect(302, "/tickets/%d", ticket.ID)
}

// Custom functions:

// Find target user based on the context and permission. (SECURITY)
func setTicket(c buffalo.Context) (ticket *models.Ticket, err error) {
	tx := c.Value("tx").(*pop.Connection)
	ticket = &models.Ticket{}
	if c.Session().Get("is_admin").(bool) {
		err = tx.Find(ticket, c.Param("ticket_id"))
	} else {
		single := getCurrentSingle(c)
		ticket = single.Ticket(c.Param("ticket_id"))
		if ticket == nil {
			err = c.Error(404, errors.New("Ticket Not Found"))
		}
	}
	c.Logger().Debugf("setTicket() returns ticket: %v, err: %v", ticket, err)
	return
}

// Custom helpers
func addTicketHelpers(c buffalo.Context) {
	c.Set("user", func(id int) interface{} {
		if id == 0 {
			return "Employee"
		}
		model := &models.User{}
		err := c.Value("tx").(*pop.Connection).Find(model, id)
		if err != nil {
			return id
		}
		return model.Username
	})
	c.Set("trimSubject", func(s string, id int) interface{} {
		model := &models.TicketSubject{}
		err := c.Value("tx").(*pop.Connection).Find(model, id)
		if err != nil {
			return s
		}
		if ns := strings.TrimPrefix(s, model.Name+" - "); len(ns) > 0 {
			return "... " + ns
		} else {
			return s
		}
	})
	c.Set("subject", func(id int) interface{} {
		model := &models.TicketSubject{}
		err := c.Value("tx").(*pop.Connection).Find(model, id)
		if err != nil {
			return id
		}
		return model.Name
	})
	c.Set("group", func(id int) interface{} {
		model := &models.TicketGroup{}
		err := c.Value("tx").(*pop.Connection).Find(model, id)
		if err != nil {
			return id
		}
		return model.Name
	})
	c.Set("status", func(id int) interface{} {
		model := &models.TicketStatus{}
		err := c.Value("tx").(*pop.Connection).Find(model, id)
		if err != nil {
			return id
		}
		return model.Name
	})
	c.Set("ticketTag", func(t models.Ticket) interface{} {
		tag := "ticket-item"
		if time.Now().Sub(t.ModifyDate).Hours() < (14 * 24) {
			tag += " ticket-new"
		}
		return tag
	})
}

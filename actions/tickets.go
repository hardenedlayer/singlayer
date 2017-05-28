package actions

import (
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"

	"github.com/hardenedlayer/singlayer/models"
)

type TicketsResource struct {
	buffalo.Resource
}

func (v TicketsResource) List(c buffalo.Context) error {
	tickets := &models.Tickets{}
	if c.Session().Get("is_admin").(bool) {
		tx := c.Value("tx").(*pop.Connection)
		err := tx.All(tickets)
		if err != nil {
			return err
		}
	} else {
		ticks := &models.Tickets{}
		single := getCurrentSingle(c)
		actor := c.Value("actor").(string)
		if actor == "All" {
			c.Logger().Debugf("multi mode for single!")
			ticks = single.MyTickets()
		} else {
			user := single.UserByUsername(c.Value("actor"))
			if user == nil {
				c.Logger().Errorf("SECURITY: cannot found user for %v", actor)
				c.Flash().Add("warning", "Oops! Who are you?")
			} else {
				c.Logger().Debugf("single mode for %v, %v.", actor, user)
				models.SyncTickets(user)
				ticks = user.Tickets()
			}
			//c.Logger().Debugf("ticks: %v --", ticks)
		}
		if ticks == nil {
			c.Flash().Add("danger", "Oops! cannot search on tickets!")
		} else {
			tickets = ticks
		}
	}
	c.Set("user", func(id int) interface{} {
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
		if ns := strings.TrimPrefix(s, model.Name + " - "); len(ns) > 0 {
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
	c.Set("tickets", tickets)
	return c.Render(200, r.HTML("tickets/index.html"))
}

func (v TicketsResource) Show(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticket := &models.Ticket{}
	err := tx.Find(ticket, c.Param("ticket_id"))
	if err != nil {
		return err
	}
	c.Set("ticket", ticket)
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
	tx := c.Value("tx").(*pop.Connection)
	ticket := &models.Ticket{}
	err := tx.Find(ticket, c.Param("ticket_id"))
	if err != nil {
		return err
	}
	c.Set("ticket", ticket)
	return c.Render(200, r.HTML("tickets/edit.html"))
}

func (v TicketsResource) Update(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticket := &models.Ticket{}
	err := tx.Find(ticket, c.Param("ticket_id"))
	if err != nil {
		return err
	}
	err = c.Bind(ticket)
	if err != nil {
		return err
	}
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

func (v TicketsResource) Destroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticket := &models.Ticket{}
	err := tx.Find(ticket, c.Param("ticket_id"))
	if err != nil {
		return err
	}
	err = tx.Destroy(ticket)
	if err != nil {
		return err
	}
	c.Flash().Add("success", "Ticket was destroyed successfully")
	return c.Redirect(302, "/tickets")
}

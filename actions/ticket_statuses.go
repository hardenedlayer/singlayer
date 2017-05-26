package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

type TicketStatusesResource struct {
	buffalo.Resource
}

// Sync all TicketStatuses from original site.
func SyncTicketStatuses(c buffalo.Context) error {
	user := &models.User{}
	tx := c.Value("tx").(*pop.Connection)
	err := tx.Where("single_id=?", c.Session().Get("user_id")).First(user)
	if err != nil {
		c.Logger().Errorf("TX error: %v", err)
	}
	c.Logger().Infof("using %v:%v", user.Username, user.APIKey)
	err = models.SyncTicketStatuses(user)
	if err == nil {
		c.Flash().Add("success", "TicketStatuses were synced successfully")
	} else {
		c.Flash().Add("danger", "Cannot sync TicketStatuses")
	}
	return c.Redirect(302, "/n/meta/tickets")
}

// List gets all TicketStatuses. This function is mapped to the the path
// GET /ticket_statuses
func (v TicketStatusesResource) List(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketStatuses := &models.TicketStatuses{}
	err := tx.All(ticketStatuses)
	if err != nil {
		return err
	}
	c.Set("ticketStatuses", ticketStatuses)
	return c.Render(200, r.HTML("ticket_statuses/index.html"))
}

// Show gets the data for one TicketStatuse. This function is mapped to
// the path GET /ticket_statuses/{ticket_status_id}
func (v TicketStatusesResource) Show(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketStatus := &models.TicketStatus{}
	err := tx.Find(ticketStatus, c.Param("ticket_status_id"))
	if err != nil {
		return err
	}
	c.Set("ticketStatus", ticketStatus)
	return c.Render(200, r.HTML("ticket_statuses/show.html"))
}

// New renders the formular for creating a new ticket_status.
// This function is mapped to the path GET /ticket_statuses/new
func (v TicketStatusesResource) New(c buffalo.Context) error {
	c.Set("ticketStatus", &models.TicketStatus{})
	return c.Render(200, r.HTML("ticket_statuses/new.html"))
}

// Create adds a ticket_status to the DB. This function is mapped to the
// path POST /ticket_statuses
func (v TicketStatusesResource) Create(c buffalo.Context) error {
	ticketStatus := &models.TicketStatus{}
	err := c.Bind(ticketStatus)
	if err != nil {
		return err
	}
	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndCreate(ticketStatus)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("ticketStatus", ticketStatus)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("ticket_statuses/new.html"))
	}
	c.Flash().Add("success", "TicketStatus was created successfully")
	return c.Redirect(302, "/ticket_statuses/%d", ticketStatus.ID)
}

// Edit renders a edit formular for a ticket_status. This function is
// mapped to the path GET /ticket_statuses/{ticket_status_id}/edit
func (v TicketStatusesResource) Edit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketStatus := &models.TicketStatus{}
	err := tx.Find(ticketStatus, c.Param("ticket_status_id"))
	if err != nil {
		return err
	}
	c.Set("ticketStatus", ticketStatus)
	return c.Render(200, r.HTML("ticket_statuses/edit.html"))
}

// Update changes a ticket_status in the DB. This function is mapped to
// the path PUT /ticket_statuses/{ticket_status_id}
func (v TicketStatusesResource) Update(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketStatus := &models.TicketStatus{}
	err := tx.Find(ticketStatus, c.Param("ticket_status_id"))
	if err != nil {
		return err
	}
	err = c.Bind(ticketStatus)
	if err != nil {
		return err
	}
	verrs, err := tx.ValidateAndUpdate(ticketStatus)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("ticketStatus", ticketStatus)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("ticket_statuses/edit.html"))
	}
	c.Flash().Add("success", "TicketStatus was updated successfully")
	return c.Redirect(302, "/ticket_statuses/%d", ticketStatus.ID)
}

// Destroy deletes a ticket_status from the DB. This function is mapped
// to the path DELETE /ticket_statuses/{ticket_status_id}
func (v TicketStatusesResource) Destroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketStatus := &models.TicketStatus{}
	err := tx.Find(ticketStatus, c.Param("ticket_status_id"))
	if err != nil {
		return err
	}
	err = tx.Destroy(ticketStatus)
	if err != nil {
		return err
	}
	c.Flash().Add("success", "TicketStatus was destroyed successfully")
	return c.Redirect(302, "/ticket_statuses")
}

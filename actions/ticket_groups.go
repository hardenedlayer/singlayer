package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

type TicketGroupsResource struct {
	buffalo.Resource
}

// Sync all TicketGroups from original site.
func SyncTicketGroups(c buffalo.Context) error {
	user := &models.User{}
	tx := c.Value("tx").(*pop.Connection)
	err := tx.Where("single_id=?", c.Session().Get("user_id")).First(user)
	if err != nil {
		c.Logger().Errorf("TX error: %v", err)
	}
	c.Logger().Infof("using %v:%v", user.Username, user.APIKey)
	err = models.SyncTicketGroups(user)
	if err == nil {
		c.Flash().Add("success", "TicketGroups were synced successfully")
	} else {
		c.Flash().Add("danger", "Cannot sync TicketGroups")
	}
	return c.Redirect(302, "/n/meta/tickets")
}

// List gets all TicketGroups. This function is mapped to the the path
// GET /ticket_groups
func (v TicketGroupsResource) List(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketGroups := &models.TicketGroups{}
	err := tx.All(ticketGroups)
	if err != nil {
		return err
	}
	c.Set("ticketGroups", ticketGroups)
	return c.Render(200, r.HTML("ticket_groups/index.html"))
}

// Show gets the data for one TicketGroup. This function is mapped to
// the path GET /ticket_groups/{ticket_group_id}
func (v TicketGroupsResource) Show(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketGroup := &models.TicketGroup{}
	err := tx.Find(ticketGroup, c.Param("ticket_group_id"))
	if err != nil {
		return err
	}
	c.Set("ticketGroup", ticketGroup)
	return c.Render(200, r.HTML("ticket_groups/show.html"))
}

// New renders the formular for creating a new ticket_group.
// This function is mapped to the path GET /ticket_groups/new
func (v TicketGroupsResource) New(c buffalo.Context) error {
	c.Set("ticketGroup", &models.TicketGroup{})
	return c.Render(200, r.HTML("ticket_groups/new.html"))
}

// Create adds a ticket_group to the DB. This function is mapped to the
// path POST /ticket_groups
func (v TicketGroupsResource) Create(c buffalo.Context) error {
	ticketGroup := &models.TicketGroup{}
	err := c.Bind(ticketGroup)
	if err != nil {
		return err
	}
	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndCreate(ticketGroup)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("ticketGroup", ticketGroup)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("ticket_groups/new.html"))
	}
	c.Flash().Add("success", "TicketGroup was created successfully")
	return c.Redirect(302, "/ticket_groups/%d", ticketGroup.ID)
}

// Edit renders a edit formular for a ticket_group. This function is
// mapped to the path GET /ticket_groups/{ticket_group_id}/edit
func (v TicketGroupsResource) Edit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketGroup := &models.TicketGroup{}
	err := tx.Find(ticketGroup, c.Param("ticket_group_id"))
	if err != nil {
		return err
	}
	c.Set("ticketGroup", ticketGroup)
	return c.Render(200, r.HTML("ticket_groups/edit.html"))
}

// Update changes a ticket_group in the DB. This function is mapped to
// the path PUT /ticket_groups/{ticket_group_id}
func (v TicketGroupsResource) Update(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketGroup := &models.TicketGroup{}
	err := tx.Find(ticketGroup, c.Param("ticket_group_id"))
	if err != nil {
		return err
	}
	err = c.Bind(ticketGroup)
	if err != nil {
		return err
	}
	verrs, err := tx.ValidateAndUpdate(ticketGroup)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("ticketGroup", ticketGroup)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("ticket_groups/edit.html"))
	}
	c.Flash().Add("success", "TicketGroup was updated successfully")
	return c.Redirect(302, "/ticket_groups/%d", ticketGroup.ID)
}

// Destroy deletes a ticket_group from the DB. This function is mapped
// to the path DELETE /ticket_groups/{ticket_group_id}
func (v TicketGroupsResource) Destroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketGroup := &models.TicketGroup{}
	err := tx.Find(ticketGroup, c.Param("ticket_group_id"))
	if err != nil {
		return err
	}
	err = tx.Destroy(ticketGroup)
	if err != nil {
		return err
	}
	c.Flash().Add("success", "TicketGroup was destroyed successfully")
	return c.Redirect(302, "/ticket_groups")
}

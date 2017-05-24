package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Ticket)
// DB Table: Plural (Tickets)
// Resource: Plural (Tickets)
// Path: Plural (/tickets)
// View Template Folder: Plural (/templates/tickets/)

// TicketsResource is the resource for the ticket model
type TicketsResource struct {
	buffalo.Resource
}

// List gets all Tickets. This function is mapped to the the path
// GET /tickets
func (v TicketsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	tickets := &models.Tickets{}
	// You can order your list here. Just change
	err := tx.All(tickets)
	// to:
	// err := tx.Order("(case when completed then 1 else 2 end) desc, lower([sort_parameter]) asc").All(tickets)
	// Don't forget to change [sort_parameter] to the parameter of
	// your model, which should be used for sorting.
	if err != nil {
		return err
	}
	// Make tickets available inside the html template
	c.Set("tickets", tickets)
	return c.Render(200, r.HTML("tickets/index.html"))
}

// Show gets the data for one Ticket. This function is mapped to
// the path GET /tickets/{ticket_id}
func (v TicketsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Ticket
	ticket := &models.Ticket{}
	// To find the Ticket the parameter ticket_id is used.
	err := tx.Find(ticket, c.Param("ticket_id"))
	if err != nil {
		return err
	}
	// Make ticket available inside the html template
	c.Set("ticket", ticket)
	return c.Render(200, r.HTML("tickets/show.html"))
}

// New renders the formular for creating a new ticket.
// This function is mapped to the path GET /tickets/new
func (v TicketsResource) New(c buffalo.Context) error {
	// Make ticket available inside the html template
	c.Set("ticket", &models.Ticket{})
	return c.Render(200, r.HTML("tickets/new.html"))
}

// Create adds a ticket to the DB. This function is mapped to the
// path POST /tickets
func (v TicketsResource) Create(c buffalo.Context) error {
	// Allocate an empty Ticket
	ticket := &models.Ticket{}
	// Bind ticket to the html form elements
	err := c.Bind(ticket)
	if err != nil {
		return err
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(ticket)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		// Make ticket available inside the html template
		c.Set("ticket", ticket)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("tickets/new.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Ticket was created successfully")
	// and redirect to the tickets index page
	return c.Redirect(302, "/tickets/%s", ticket.ID)
}

// Edit renders a edit formular for a ticket. This function is
// mapped to the path GET /tickets/{ticket_id}/edit
func (v TicketsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Ticket
	ticket := &models.Ticket{}
	err := tx.Find(ticket, c.Param("ticket_id"))
	if err != nil {
		return err
	}
	// Make ticket available inside the html template
	c.Set("ticket", ticket)
	return c.Render(200, r.HTML("tickets/edit.html"))
}

// Update changes a ticket in the DB. This function is mapped to
// the path PUT /tickets/{ticket_id}
func (v TicketsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Ticket
	ticket := &models.Ticket{}
	err := tx.Find(ticket, c.Param("ticket_id"))
	if err != nil {
		return err
	}
	// Bind ticket to the html form elements
	err = c.Bind(ticket)
	if err != nil {
		return err
	}
	verrs, err := tx.ValidateAndUpdate(ticket)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		// Make ticket available inside the html template
		c.Set("ticket", ticket)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("tickets/edit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Ticket was updated successfully")
	// and redirect to the tickets index page
	return c.Redirect(302, "/tickets/%s", ticket.ID)
}

// Destroy deletes a ticket from the DB. This function is mapped
// to the path DELETE /tickets/{ticket_id}
func (v TicketsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Ticket
	ticket := &models.Ticket{}
	// To find the Ticket the parameter ticket_id is used.
	err := tx.Find(ticket, c.Param("ticket_id"))
	if err != nil {
		return err
	}
	err = tx.Destroy(ticket)
	if err != nil {
		return err
	}
	// If there are no errors set a flash message
	c.Flash().Add("success", "Ticket was destroyed successfully")
	// Redirect to the tickets index page
	return c.Redirect(302, "/tickets")
}

package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"

	"github.com/hardenedlayer/singlayer/models"
)

// HomeHandler is a default handler to serve up
// a home page.
func TicketsMetaHandler(c buffalo.Context) (err error) {
	tx := c.Value("tx").(*pop.Connection)
	ticket_subjects := &models.TicketSubjects{}
	err = tx.All(ticket_subjects)
	if err != nil {
		return err
	}
	c.Set("ticket_subjects", ticket_subjects)

	ticket_groups := &models.TicketGroups{}
	err = tx.All(ticket_groups)
	if err != nil {
		return err
	}
	c.Set("ticket_groups", ticket_groups)

	ticket_statuses := &models.TicketStatuses{}
	err = tx.All(ticket_statuses)
	if err != nil {
		return err
	}
	c.Set("ticket_statuses", ticket_statuses)

	c.Set("theme", "default")
	return c.Render(200, r.HTML("n.meta.tickets.html"))
}

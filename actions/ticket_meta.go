package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

func SyncTicketGroups(c buffalo.Context) error {
	user := &models.User{}
	tx := c.Value("tx").(*pop.Connection)
	err := tx.First(user)
	if err != nil {
		c.Logger().Errorf("TX error: %v", err)
	}
	c.Logger().Infof("using %v", user.Username)
	err = models.SyncTicketGroups(user)
	if err == nil {
		c.Flash().Add("success", "TicketGroups were synced successfully")
	} else {
		c.Flash().Add("danger", "Cannot sync TicketGroups")
	}
	return c.Redirect(302, "/n/meta/tickets")
}

func SyncTicketStatuses(c buffalo.Context) error {
	user := &models.User{}
	tx := c.Value("tx").(*pop.Connection)
	err := tx.First(user)
	if err != nil {
		c.Logger().Errorf("TX error: %v", err)
	}
	c.Logger().Infof("using %v", user.Username)
	err = models.SyncTicketStatuses(user)
	if err == nil {
		c.Flash().Add("success", "TicketStatuses were synced successfully")
	} else {
		c.Flash().Add("danger", "Cannot sync TicketStatuses")
	}
	return c.Redirect(302, "/n/meta/tickets")
}

func SyncTicketSubjects(c buffalo.Context) error {
	user := &models.User{}
	tx := c.Value("tx").(*pop.Connection)
	err := tx.First(user)
	if err != nil {
		c.Logger().Errorf("TX error: %v", err)
	}
	c.Logger().Infof("using %v", user.Username)
	err = models.SyncTicketSubjects(user)
	if err == nil {
		c.Flash().Add("success", "TicketSubjects were synced successfully")
	} else {
		c.Flash().Add("danger", "Cannot sync TicketSubjects")
	}
	return c.Redirect(302, "/n/meta/tickets")
}

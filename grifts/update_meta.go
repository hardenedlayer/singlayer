package grifts

import (
	"github.com/markbates/grift/grift"
	"github.com/markbates/pop"

	"github.com/hardenedlayer/singlayer/models"
)

var _ = grift.Add("db:meta:ticket", func(c *grift.Context) error {
	return models.DB.Transaction(func(tx *pop.Connection) error {
		p := models.Logger.Printf
		user := &models.User{}
		err := tx.First(user)
		p("using %v:%v", user.Username, user.APIKey)

		p("perform SyncTicketSubjects...")
		err = models.SyncTicketSubjects(user)
		if err != nil {
			p("error while running SyncTicketSubjects: %v", err)
		} else {
			p("SyncTicketSubjects finished successfully.")
		}

		p("perform SyncTicketGroups...")
		err = models.SyncTicketGroups(user)
		if err != nil {
			p("error while running SyncTicketGroups: %v", err)
		} else {
			p("SyncTicketGroups finished successfully.")
		}

		p("perform SyncTicketStatuses...")
		err = models.SyncTicketStatuses(user)
		if err != nil {
			p("error while running SyncTicketStatuses: %v", err)
		} else {
			p("SyncTicketStatuses finished successfully.")
		}
		return nil
	})
})

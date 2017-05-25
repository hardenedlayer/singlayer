package grifts

import (
	"github.com/jinzhu/copier"
	"github.com/markbates/grift/grift"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

var _ = grift.Add("db:update:meta", func(c *grift.Context) error {
	return models.DB.Transaction(func(tx *pop.Connection) error {
		p := models.Logger.Printf
		user := &models.User{}
		err := tx.First(user)
		p("using %v:%v", user.Username, user.APIKey)
		p("")
		sess := session.New(user.Username, user.APIKey)
		sess.Endpoint = "https://api.softlayer.com/rest/v3.1"

		p("update ticket_status...")
		p("")
		ticket_service := services.GetTicketService(sess)
		ticket_statuses, err := ticket_service.GetAllTicketStatuses()
		if err != nil {
			return err
		}
		for _,els := range ticket_statuses {
			ts := &models.TicketStatus{}
			copier.Copy(ts, els)
			ts.ID = *els.Id
			if ok, _ := pop.Q(tx).Where("id=?", ts.ID).Exists(ts); ok {
				p("ticket_group %v already exists!", ts.ID)
			} else {
				p("about to create ticket_group %v...", ts)
				err = tx.Create(ts)
				if err != nil {
					p("cannot create ticket_group:%v", err)
				}
			}
		}

		p("update ticket_group...")
		p("")
		ticket_groups, err := ticket_service.GetAllTicketGroups()
		if err != nil {
			return err
		}
		for _,elg := range ticket_groups {
			tg := &models.TicketGroup{}
			copier.Copy(tg, elg)
			tg.ID = *elg.Id
			if ok, _ := pop.Q(tx).Where("id=?", tg.ID).Exists(tg); ok {
				p("ticket_group %v already exists!", tg.ID)
			} else {
				p("about to create ticket_group %v...", tg)
				err = tx.Create(tg)
				if err != nil {
					p("cannot create ticket_group:%v", err)
				}
			}
		}

		p("update ticket_subject...")
		p("")
		subject_service := services.GetTicketSubjectService(sess)
		ticket_subjects, err := subject_service.GetAllObjects()
		if err != nil {
			return err
		}
		for _,els := range ticket_subjects {
			ts := &models.TicketSubject{}
			copier.Copy(ts, els)
			ts.ID = *els.Id
			if ok, _ := pop.Q(tx).Where("id=?", ts.ID).Exists(ts); ok {
				p("ticket_subject %v already exists!", ts.ID)
			} else {
				p("about to create ticket_subject %v...", ts)
				err = tx.Create(ts)
				if err != nil {
					p("cannot create ticket_subject:%v", err)
				}
			}
		}
		return nil
	})
})

package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/copier"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type TicketGroup struct {
	ID                    int       `json:"id" db:"id"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
	TicketGroupCategoryId int       `json:"ticket_group_category_id" db:"ticket_group_category_id"`
	Name                  string    `json:"name" db:"name"`
}

func (t TicketGroup) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

type TicketGroups []TicketGroup

func (t TicketGroups) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (t *TicketGroup) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: t.ID, Name: "ID"},
		&validators.IntIsPresent{Field: t.TicketGroupCategoryId, Name: "TicketGroupCategoryId"},
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (t *TicketGroup) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (t *TicketGroup) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (t *TicketGroup) IsNew() bool {
	return time.Now().Sub(t.CreatedAt) < time.Duration(24 * 14 * time.Hour)
}

//// Backend API Calls:

// SyncTicketGroups() creates and updates all Ticket Groups.
func SyncTicketGroups(user *User) error {
	log.Infof("sync ticket groups...(use %v)", user.Username)
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"

	service := services.GetTicketService(sess)
	data, err := service.GetAllTicketGroups()
	if err != nil {
		return err
	}
	for _, el := range data {
		ts := &TicketGroup{}
		copier.Copy(ts, el)
		ts.ID = *el.Id
		if ok, _ := DB.Where("id=?", ts.ID).Exists(ts); ok {
			log.Debugf("ticket_group %v already exists!", ts.ID)
		} else {
			err = DB.Create(ts)
			if err != nil {
				log.Errorf("cannot create ticket_group:%v", err)
			} else {
				log.Debugf("ticket_group %v created.", ts)
			}
		}
	}
	return nil
}

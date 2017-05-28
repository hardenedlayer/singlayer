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

type TicketStatus struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name      string    `json:"name" db:"name"`
}

func (t TicketStatus) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

type TicketStatuses []TicketStatus

func (t TicketStatuses) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (t *TicketStatus) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: t.ID, Name: "ID"},
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (t *TicketStatus) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (t *TicketStatus) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (t *TicketStatus) IsNew() bool {
	return time.Now().Sub(t.CreatedAt) < time.Duration(24 * 14 * time.Hour)
}

//// Backend API Calls:

// SyncTicketStatuses() creates and updates all Ticket Statuses.
func SyncTicketStatuses(user *User) error {
	Logger.Printf("sync ticket statuses... (use %v)", user.Username)
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"

	service := services.GetTicketService(sess)
	data, err := service.GetAllTicketStatuses()
	if err != nil {
		return err
	}
	for _, el := range data {
		ts := &TicketStatus{}
		copier.Copy(ts, el)
		ts.ID = *el.Id
		if ok, _ := DB.Where("id=?", ts.ID).Exists(ts); ok {
			Logger.Printf("ticket_status %v already exists!", ts.ID)
		} else {
			err = DB.Create(ts)
			if err != nil {
				Logger.Printf("cannot create ticket_status:%v", err)
			} else {
				Logger.Printf("ticket_status %v created.", ts)
			}
		}
	}
	return nil
}

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

type TicketSubject struct {
	ID         int       `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	CategoryId int       `json:"category_id" db:"category_id"`
	Name       string    `json:"name" db:"name"`
}

func (t TicketSubject) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

type TicketSubjects []TicketSubject

func (t TicketSubjects) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (t *TicketSubject) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: t.ID, Name: "ID"},
		&validators.IntIsPresent{Field: t.CategoryId, Name: "CategoryId"},
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (t *TicketSubject) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (t *TicketSubject) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (t *TicketSubject) IsNew() bool {
	return time.Now().Sub(t.CreatedAt) < time.Duration(24 * 14 * time.Hour)
}

func SyncTicketSubjects(user *User) error {
	Logger.Printf("sync ticket subjects... (use %v)", user.Username)
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"

	service := services.GetTicketSubjectService(sess)
	data, err := service.GetAllObjects()
	if err != nil {
		return err
	}
	for _, el := range data {
		ts := &TicketSubject{}
		copier.Copy(ts, el)
		ts.ID = *el.Id
		if ok, _ := DB.Where("id=?", ts.ID).Exists(ts); ok {
			Logger.Printf("ticket_subject %v already exists!", ts.ID)
		} else {
			err = DB.Create(ts)
			if err != nil {
				Logger.Printf("cannot create ticket_subject:%v", err)
			} else {
				Logger.Printf("ticket_subject %v created.", ts)
			}
		}
	}
	return nil
}

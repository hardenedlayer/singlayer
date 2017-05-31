package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
)

type TicketUpdate struct {
	ID         int       `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	TicketId   int       `json:"ticket_id" db:"ticket_id"`
	EditorId   int       `json:"editor_id" db:"editor_id"`
	EditorType string    `json:"editor_type" db:"editor_type"`
	Entry      string    `json:"entry" db:"entry"`
	CreateDate time.Time `json:"create_date" db:"create_date"`
}

// String is not required by pop and may be deleted
func (t TicketUpdate) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TicketUpdates is not required by pop and may be deleted
type TicketUpdates []TicketUpdate

// String is not required by pop and may be deleted
func (t TicketUpdates) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (t *TicketUpdate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: t.ID, Name: "ID"},
		&validators.IntIsPresent{Field: t.TicketId, Name: "TicketId"},
		&validators.IntIsPresent{Field: t.EditorId, Name: "EditorId"},
		&validators.StringIsPresent{Field: t.EditorType, Name: "EditorType"},
		&validators.StringIsPresent{Field: t.Entry, Name: "Entry"},
		&validators.TimeIsPresent{Field: t.CreateDate, Name: "CreateDate"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (t *TicketUpdate) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (t *TicketUpdate) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

//// display helpers:

// EditorName
func (t TicketUpdate) EditorName() interface{} {
	if t.EditorId == 0 {
		return "Employee"
	}
	ins := &User{}
	err := DB.Find(ins, t.EditorId)
	if err != nil {
		return t.EditorId
	}
	return ins.Username
}

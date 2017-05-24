package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
)

type TicketStatus struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name      string    `json:"name" db:"name"`
}

// String is not required by pop and may be deleted
func (t TicketStatus) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TicketStatuses is not required by pop and may be deleted
type TicketStatuses []TicketStatus

// String is not required by pop and may be deleted
func (t TicketStatuses) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (t *TicketStatus) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: t.ID, Name: "ID"},
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (t *TicketStatus) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (t *TicketStatus) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

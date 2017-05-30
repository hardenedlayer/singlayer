package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type Progress struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	PrevID    uuid.UUID `json:"prev_id" db:"prev_id"`
	OrderID   uuid.UUID `json:"order_id" db:"order_id"`
	EditorID  uuid.UUID `json:"editor_id" db:"editor_id"`
	Action    string    `json:"action" db:"action"`
}

// String is not required by pop and may be deleted
func (p Progress) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Progresses is not required by pop and may be deleted
type Progresses []Progress

// String is not required by pop and may be deleted
func (p Progresses) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (p *Progress) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Action, Name: "Action"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (p *Progress) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (p *Progress) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

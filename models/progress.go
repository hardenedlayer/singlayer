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
	ID           uuid.UUID `json:"id" db:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	PrevID       uuid.UUID `json:"prev_id" db:"prev_id"`
	DirectLinkID uuid.UUID `json:"direct_link_id" db:"direct_link_id"`
	UpdateId     int       `json:"update_id" db:"update_id"`
	SingleID     uuid.UUID `json:"single_id" db:"single_id"`
	Action       string    `json:"action" db:"action"`
}

func (p Progress) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

type Progresses []Progress

func (p Progresses) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (p *Progress) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Action, Name: "Action"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (p *Progress) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (p *Progress) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

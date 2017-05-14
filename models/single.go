package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type Single struct {
	ID           uuid.UUID    `json:"id" db:"id"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
	Name         string       `json:"name" db:"name"`
	Email        string       `json:"email" db:"email"`
	Provider     string       `json:"provider" db:"provider"`
	UserID       string       `json:"user_id" db:"user_id"`
	AvatarUrl    string       `json:"avatar_url" db:"avatar_url"`
	Organization string       `json:"organization" db:"organization"`
	Note         nulls.String `json:"note" db:"note"`
	Permissions  string       `json:"permissions" db:"permissions"`
}

// String is not required by pop and may be deleted
func (s Single) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Singles is not required by pop and may be deleted
type Singles []Single

// String is not required by pop and may be deleted
func (s Singles) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (s *Single) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.Name, Name: "Name"},
		&validators.StringIsPresent{Field: s.Email, Name: "Email"},
		&validators.StringIsPresent{Field: s.Provider, Name: "Provider"},
		&validators.StringIsPresent{Field: s.UserID, Name: "UserID"},
		&validators.StringIsPresent{Field: s.Permissions, Name: "Permissions"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (s *Single) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (s *Single) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

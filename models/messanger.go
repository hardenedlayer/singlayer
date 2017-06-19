package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type Messanger struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	SingleID  uuid.UUID `json:"single_id" db:"single_id"`
	Name      string    `json:"name" db:"name"`
	Level     string    `json:"level" db:"level"`
	Method    string    `json:"method" db:"method"`
	Value     string    `json:"value" db:"value"`
}

func (m Messanger) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

type Messangers []Messanger

func (m Messangers) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (m *Messanger) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.Name, Name: "Name"},
		&validators.StringIsPresent{Field: m.Level, Name: "Level"},
		&validators.StringIsPresent{Field: m.Method, Name: "Method"},
		&validators.StringIsPresent{Field: m.Value, Name: "Value"},
	), nil
}

func (m *Messanger) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (m *Messanger) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

//// relational objects

func (m Messanger) Single() (single *Single) {
	single, err := FindSingle(m.SingleID)
	if err == nil {
		return single
	}
	return &Single{}
}

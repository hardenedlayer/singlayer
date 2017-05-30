package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type Vlan struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ID        int       `json:"id" db:"id"`
}

// String is not required by pop and may be deleted
func (v Vlan) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// Vlans is not required by pop and may be deleted
type Vlans []Vlan

// String is not required by pop and may be deleted
func (v Vlans) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (v *Vlan) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: v.ID, Name: "ID"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (v *Vlan) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (v *Vlan) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

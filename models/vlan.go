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
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	R1LinkID  uuid.UUID `json:"r1_link_id" db:"r1_link_id"`
	R2LinkID  uuid.UUID `json:"r2_link_id" db:"r2_link_id"`
	AccountId int       `json:"account_id" db:"account_id"`
	Booked    bool      `json:"booked" db:"booked"`
}

func (v Vlan) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

type Vlans []Vlan

func (v Vlans) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (v *Vlan) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: v.ID, Name: "ID"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (v *Vlan) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (v *Vlan) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

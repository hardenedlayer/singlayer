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

type User struct {
	ID                uuid.UUID    `json:"id" db:"id"`
	CreatedAt         time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at" db:"updated_at"`
	ID                int          `json:"id" db:"id"`
	AccountID         int          `json:"account_id" db:"account_id"`
	ParentID          int          `json:"parent_id" db:"parent_id"`
	Username          string       `json:"username" db:"username"`
	ApiKey            string       `json:"api_key" db:"api_key"`
	CompanyName       string       `json:"company_name" db:"company_name"`
	Email             string       `json:"email" db:"email"`
	FirstName         string       `json:"first_name" db:"first_name"`
	LastName          string       `json:"last_name" db:"last_name"`
	OpenTicketCount   uint         `json:"open_ticket_count" db:"open_ticket_count"`
	TicketCount       uint         `json:"ticket_count" db:"ticket_count"`
	HardwareCount     uint         `json:"hardware_count" db:"hardware_count"`
	VirtualGuestCount uint         `json:"virtual_guest_count" db:"virtual_guest_count"`
	Permissions       nulls.String `json:"permissions" db:"permissions"`
	SingleID          uuid.UUID    `json:"single_id" db:"single_id"`
	LastBatch         time.Time    `json:"last_batch" db:"last_batch"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: u.ID, Name: "ID"},
		&validators.IntIsPresent{Field: u.AccountID, Name: "AccountID"},
		&validators.IntIsPresent{Field: u.ParentID, Name: "ParentID"},
		&validators.StringIsPresent{Field: u.Username, Name: "Username"},
		&validators.StringIsPresent{Field: u.ApiKey, Name: "ApiKey"},
		&validators.StringIsPresent{Field: u.CompanyName, Name: "CompanyName"},
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: u.LastName, Name: "LastName"},
		&validators.TimeIsPresent{Field: u.LastBatch, Name: "LastBatch"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (u *User) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

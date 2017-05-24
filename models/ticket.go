package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
)

type Ticket struct {
	ID             int       `json:"id" db:"id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	AccountId      int       `json:"account_id" db:"account_id"`
	AssignedUserId nulls.Int `json:"assigned_user_id" db:"assigned_user_id"`
	SubjectId      nulls.Int `json:"subject_id" db:"subject_id"`
	GroupId        nulls.Int `json:"group_id" db:"group_id"`
	StatusId       int       `json:"status_id" db:"status_id"`
	Title          string    `json:"title" db:"title"`
	CreateDate     time.Time `json:"create_date" db:"create_date"`
	LastEditDate   time.Time `json:"last_edit_date" db:"last_edit_date"`
	LastEditType   string    `json:"last_edit_type" db:"last_edit_type"`
}

// String is not required by pop and may be deleted
func (t Ticket) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tickets is not required by pop and may be deleted
type Tickets []Ticket

// String is not required by pop and may be deleted
func (t Tickets) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (t *Ticket) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: t.ID, Name: "ID"},
		&validators.IntIsPresent{Field: t.AccountId, Name: "AccountId"},
		&validators.IntIsPresent{Field: t.StatusId, Name: "StatusId"},
		&validators.StringIsPresent{Field: t.Title, Name: "Title"},
		&validators.TimeIsPresent{Field: t.CreateDate, Name: "CreateDate"},
		&validators.TimeIsPresent{Field: t.LastEditDate, Name: "LastEditDate"},
		&validators.StringIsPresent{Field: t.LastEditType, Name: "LastEditType"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (t *Ticket) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (t *Ticket) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

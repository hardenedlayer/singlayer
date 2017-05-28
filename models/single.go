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

func (s Single) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

type Singles []Single

func (s Singles) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run everytime you call a "pop.Validate" method.
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
func (s *Single) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (s *Single) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

//// Association and Relationship based search for instances.
//// It need instance of Single so more expensive than raw query. FIXME later.

// Users() returns instance of Users struct.
func (s *Single) Users() (users *Users) {
	users = &Users{}
	err := DB.BelongsTo(s).All(users)
	if err != nil {
		return nil
	}
	return
}

// Accounts() returns instance of Accounts struct.
func (s *Single) Accounts() (accounts *Accounts) {
	accounts = &Accounts{}
	err := DB.BelongsToThrough(s, "users").All(accounts)
	if err != nil {
		return nil
	}
	return
}

// User(user_id) returns single instance of User
func (s *Single) User(user_id interface{}) (user *User) {
	user = &User{}
	err := DB.BelongsTo(s).Find(user, user_id)
	if err != nil {
		return nil
	}
	return
}

// UserByUsername(username) returns single instance of User
func (s *Single) UserByUsername(username interface{}) (user *User) {
	user = &User{}
	err := DB.BelongsTo(s).Where("username=?", username).First(user)
	//err := DB.BelongsTo(s).Find(user, user_id)
	if err != nil {
		return nil
	}
	return
}

// Account(account_id) returns single instance of Account.
func (s *Single) Account(account_id interface{}) (account *Account) {
	account = &Account{}
	err := DB.BelongsToThrough(s, "users").Find(account, account_id)
	if err != nil {
		return nil
	}
	return
}

// Tickets() returns all tickets associated to the single.
func (s *Single) MyTickets() (tickets *Tickets) {
	tickets = &Tickets{}
	err := pop.Q(DB).
		LeftJoin("users", "users.id = tickets.assigned_user_id").
		Where("users.single_id = ?", s.ID).
		Order("tickets.last_edit_date desc").
		All(tickets)
	if err != nil {
		Logger.Printf("Err: %v", err)
		return nil
	}
	return
}

// Tickets() returns all tickets associated to the single.
func (s *Single) Tickets() (tickets *Tickets) {
	tickets = &Tickets{}
	err := pop.Q(DB).
		LeftJoin("accounts", "accounts.id = tickets.account_id").
		LeftJoin("users", "accounts.id = users.account_id").
		Where("users.single_id = ?", s.ID).
		Order("tickets.last_edit_date desc").
		All(tickets)
	if err != nil {
		Logger.Printf("Err: %v", err)
		return nil
	}
	return
}

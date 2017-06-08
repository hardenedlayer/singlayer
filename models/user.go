package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/copier"
	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type User struct {
	ID                int          `json:"id" db:"id"`
	CreatedAt         time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at" db:"updated_at"`
	AccountId         int          `json:"account_id" db:"account_id"`
	ParentId          int          `json:"parent_id" db:"parent_id"`
	Username          string       `json:"username" db:"username"`
	APIKey            string       `json:"api_key" db:"api_key"`
	CompanyName       string       `json:"company_name" db:"company_name"`
	Email             string       `json:"email" db:"email"`
	FirstName         string       `json:"first_name" db:"first_name"`
	LastName          string       `json:"last_name" db:"last_name"`
	OpenTicketCount   int          `json:"open_ticket_count" db:"open_ticket_count"`
	TicketCount       int          `json:"ticket_count" db:"ticket_count"`
	HardwareCount     int          `json:"hardware_count" db:"hardware_count"`
	VirtualGuestCount int          `json:"virtual_guest_count" db:"virtual_guest_count"`
	Permissions       nulls.String `json:"permissions" db:"permissions"`
	SingleID          uuid.UUID    `json:"single_id" db:"single_id"`
	LastBatch         time.Time    `json:"last_batch" db:"last_batch"`
}

func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

type Users []User

func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: u.ID, Name: "ID"},
		&validators.IntIsPresent{Field: u.AccountId, Name: "AccountId"},
		&validators.IntIsPresent{Field: u.ParentId, Name: "ParentId"},
		&validators.StringIsPresent{Field: u.Username, Name: "Username"},
		&validators.StringIsPresent{Field: u.APIKey, Name: "APIKey"},
		&validators.StringIsPresent{Field: u.CompanyName, Name: "CompanyName"},
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: u.LastName, Name: "LastName"},
		&validators.TimeIsPresent{Field: u.LastBatch, Name: "LastBatch"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (u *User) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// Backend API Calls:

// Update() fills up user struct with response of api call.
func (u *User) Update() (err error) {
	sess := session.New(u.Username, u.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"
	service := services.GetAccountService(sess)
	sl_user, err := service.
		Mask("id;accountId;parentId;companyName;email;firstName;lastName;ticketCount;openTicketCount;hardwareCount;virtualGuestCount").
		GetCurrentUser()
	if err != nil {
		log.Errorf("softlayer api exception: %v --", err)
		return err
	}
	copier.Copy(u, sl_user)
	u.ID = *sl_user.Id
	inspect("updated user", u)
	return
}

//// Association and Relationship based search for instances.

// Account() returns related single Account instance for the User instance.
func (u *User) Account() (account *Account) {
	account = &Account{}
	err := DB.Find(account, u.AccountId)
	if err != nil {
		return nil
	}
	return
}

// Single() returns associated Single instance for this user.
func (u User) Single() interface{} {
	single := &Single{}
	err := DB.Find(single, u.SingleID)
	if err != nil {
		return "Unknown Singler"
	}
	return single
}

// MyTickets() returns all tickets assigned to me.
func (u *User) MyTickets(page, pp int) (*Tickets, *pop.Paginator) {
	tickets := &Tickets{}
	q := pop.Q(DB).Paginate(page, pp)
	err := q.Where("assigned_user_id = ?", u.ID).
		Order("tickets.last_edit_date desc").
		All(tickets)
	if err != nil {
		return nil, nil
	}
	return tickets, q.Paginator
}

// Tickets() returns all tickets from the user's account.
func (u *User) Tickets(page, pp int) (*Tickets, *pop.Paginator) {
	tickets := &Tickets{}
	q := pop.Q(DB).Paginate(page, pp)
	err := q.Where("account_id = ?", u.AccountId).
		Order("tickets.last_edit_date desc").
		All(tickets)
	if err != nil {
		return nil, nil
	}
	return tickets, q.Paginator
}

// DirectLinks() returns all directlinks from the user's account.
func (u *User) DirectLinks() (dlinks *DirectLinks) {
	dlinks = &DirectLinks{}
	err := DB.Where("account_id = ?", u.AccountId).
		Order("direct_links.created_at desc").
		All(dlinks)
	if err != nil {
		return nil
	}
	return
}

//// search functions

// Find a user with user_id.
func FindUser(user_id int) (user *User, err error) {
	user = &User{}
	err = DB.Find(user, user_id)
	return
}

// display functions:

func (u User) Contact() interface{} {
	return u.FirstName + " " + u.LastName + " <" + u.Email + ">"
}

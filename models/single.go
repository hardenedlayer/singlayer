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
	return s.Name
}

func (s Single) Marshal() string {
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

// Messanger() returns single messanger for this singler.
func (s *Single) Messanger(messanger_id interface{}) (messanger *Messanger) {
	messanger = &Messanger{}
	err := DB.BelongsTo(s).Find(messanger, messanger_id)
	if err != nil {
		return nil
	}
	return
}

// Messangers() returns all messanger for this singler.
func (s *Single) Messangers(levels ...string) (messangers *Messangers) {
	messangers = &Messangers{}
	q := DB.BelongsTo(s)
	for _, e := range levels { // tricky optional single argument.
		q = q.Where("level = ?", e)
	}
	q.Order("level, updated_at desc").All(messangers)
	return
}

func (s Single) NotifyTo() (messangers *Messangers) {
	messangers = s.Messangers("notification")
	m := &Messanger{}
	err := DB.BelongsTo(&s).
		Where("level = ? AND method = ?", "notification", "mail").First(m)
	if err != nil {
		m = &Messanger{
			Level: "Notification",
			Method: "Mail",
			Value: s.Email,
		}
		*messangers = append(*messangers, *m)
	}
	return
}

func (s Single) AlertTo() (messangers *Messangers) {
	return s.Messangers("alert")
}

func (s Single) Mail() (mail string) {
	mail = s.Email
	m := &Messanger{}
	err := DB.BelongsTo(&s).
		Where("level = ? AND method = ?", "notification", "mail").
		Order("updated_at desc").
		First(m)
	if err == nil {
		mail = m.Value
	}
	return
}

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

func (s Single) AccountCount() interface{} {
	accounts := &Accounts{}
	count, err := DB.BelongsToThrough(&s, "users").Count(accounts)
	if err == nil {
		return count
	}
	return -9
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

// UserByAccount(account_id) returns single instance of User
func (s *Single) UserByAccount(account_id interface{}) (user *User) {
	user = &User{}
	err := DB.BelongsTo(s).Where("account_id=?", account_id).First(user)
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
func (s *Single) MyTickets(page, pp int) (*Tickets, *pop.Paginator) {
	tickets := &Tickets{}
	q := pop.Q(DB).Paginate(page, pp)
	err := q.
		LeftJoin("users", "users.id = tickets.assigned_user_id").
		Where("users.single_id = ?", s.ID).
		Order("tickets.last_edit_date desc").
		All(tickets)
	if err != nil {
		log.Errorf("Err: %v", err)
		return nil, nil
	}
	return tickets, q.Paginator
}

// Tickets() returns all tickets associated to the single.
func (s *Single) Tickets(page, pp int) (*Tickets, *pop.Paginator) {
	tickets := &Tickets{}
	q := pop.Q(DB).Paginate(page, pp)
	err := q.
		LeftJoin("accounts", "accounts.id = tickets.account_id").
		LeftJoin("users", "accounts.id = users.account_id").
		Where("users.single_id = ?", s.ID).
		Order("tickets.last_edit_date desc").
		All(tickets)
	if err != nil {
		log.Errorf("Err: %v", err)
		return nil, nil
	}
	return tickets, q.Paginator
}

// Ticket() returns specified ticket associated to the single.
func (s *Single) Ticket(ticket_id interface{}) (ticket *Ticket) {
	ticket = &Ticket{}
	err := pop.Q(DB).
		LeftJoin("accounts", "accounts.id = tickets.account_id").
		LeftJoin("users", "accounts.id = users.account_id").
		Where("users.single_id = ?", s.ID).
		Find(ticket, ticket_id)
	if err != nil {
		log.Errorf("Err: %v", err)
		return nil
	}
	return
}

// DirectLinks() returns all directlinks associated to the single's accounts.
func (s *Single) DirectLinks() *DirectLinks {
	dlinks := &DirectLinks{}
	err := pop.Q(DB).
		LeftJoin("accounts", "accounts.id = direct_links.account_id").
		LeftJoin("users", "accounts.id = users.account_id").
		Where("users.single_id = ?", s.ID).
		Order("direct_links.created_at desc").
		All(dlinks)
	if err != nil {
		log.Errorf("Err: %v", err)
		return nil
	}
	return dlinks
}

// MyDirectLinks() returns all directlinks associated to the single directly.
func (s *Single) MyDirectLinks() (dlinks *DirectLinks) {
	dlinks = &DirectLinks{}
	err := DB.BelongsTo(s).Order("account_id, line_number").All(dlinks)
	if err != nil {
		return nil
	}
	return
}

// Computes() returns all compute instances associated with this single.
func (s *Single) Computes(page, pp int) (*Computes, *pop.Paginator) {
	comps := &Computes{}
	q := pop.Q(DB).Paginate(page, pp)
	err := q.
		LeftJoin("comp_user_maps", "comp_user_maps.compute_id = computes.id").
		LeftJoin("users", "comp_user_maps.user_id = users.id").
		Where("users.single_id = ?", s.ID).
		Order("account_id, hostname").
		All(comps)
	if err != nil {
		log.Errorf("Err: %v", err)
		return nil, nil
	}
	return comps, q.Paginator
}

// Compute() specified compute instance associated with this single.
func (s *Single) Compute(compute_id interface{}) (compute *Compute) {
	compute = &Compute{}
	err := pop.Q(DB).
		LeftJoin("comp_user_maps", "comp_user_maps.compute_id = computes.id").
		LeftJoin("users", "comp_user_maps.user_id = users.id").
		Where("users.single_id = ?", s.ID).
		Find(compute, compute_id)
	if err != nil {
		log.Errorf("Err: %v", err)
		return nil
	}
	return
}

//// search functions

// Find and a single with single_id
func FindSingle(single_id uuid.UUID) (single *Single, err error) {
	single = &Single{}
	err = DB.Find(single, single_id)
	return
}

func GetSinglesByPermission(perm string) (singles *Singles, err error) {
	singles = &Singles{}
	err = DB.Where("permissions LIKE ?", "%:"+perm+":%").All(singles)
	return
}

//// utilities:

// AdminMail wrapper
func (s *Single) AdminMail(obj Object, subj, to string, grps ...string) error {
	return AdminMail(s, obj, subj, to, grps...)
}

package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"text/template"
	"time"

	"github.com/jinzhu/copier"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

type Ticket struct {
	ID               int       `json:"id" db:"id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	AccountId        int       `json:"account_id" db:"account_id"`
	AssignedUserId   int       `json:"assigned_user_id" db:"assigned_user_id"`
	SubjectId        int       `json:"subject_id" db:"subject_id"`
	GroupId          int       `json:"group_id" db:"group_id"`
	StatusId         int       `json:"status_id" db:"status_id"`
	Title            string    `json:"title" db:"title"`
	TotalUpdateCount int       `json:"total_update_count" db:"total_update_count"`
	CreateDate       time.Time `json:"create_date" db:"create_date"`
	LastEditDate     time.Time `json:"last_edit_date" db:"last_edit_date"`
	LastEditType     string    `json:"last_edit_type" db:"last_edit_type"`
	LastSync         time.Time `json:"last_sync" db:"last_sync"`
}

func (t Ticket) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

type Tickets []Ticket

func (t Tickets) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

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

func (t *Ticket) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (t *Ticket) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

//// Backend API Calls:

// CreateDirectLinkTicket()
func CreateDirectLinkTicket(user *User, dlink *DirectLink) (int, error) {
	log.Infof("create order ticket... (use %v)", user.Username)
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"

	var sl_ticket datatypes.Ticket
	sl_ticket.AssignedUserId = sl.Int(user.ID)
	sl_ticket.SubjectId = sl.Int(1002)
	sl_ticket.GroupId = sl.Int(1014)
	sl_ticket.Title = sl.String("SK Direct Link " + dlink.Type + " - SEO01")
	inspect("new ticket", sl_ticket)

	tmpl, err := template.ParseFiles("templates/order.templ")
	if err != nil {
		return 0, err
	}
	buf := &bytes.Buffer{}
	if err = tmpl.Execute(buf, dlink); err != nil {
		return 0, err
	}
	content := buf.String()
	log.Debugf("new ticket contents: -------\n%v\n-------", content)

	if is_test {
		log.Infof("XXXXXXXXXXXXXXXXXXXXXXX EVERYTHING IS GOOD! BUT TEST.")
		return 0, nil
	}

	service := services.GetTicketService(sess)
	tick, err := service.CreateStandardTicket(&sl_ticket, &content,
		nil, nil, nil, nil, nil, nil)
	if err != nil {
		return 0, err
	}
	inspect("saved sl_ticket", tick)

	SyncTickets(user)
	return *tick.Id, nil
}

// AddUpdate() add an update to ticket and returns first added update.
func (t *Ticket) AddUpdate(user *User, entry string) (*TicketUpdate, error) {
	log.Infof("add update to ticket... (use %v)", user.Username)
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"

	var sl_update datatypes.Ticket_Update
	sl_update.Entry = sl.String(entry)
	inspect("new update", sl_update)

	service := services.GetTicketService(sess)
	upds, err := service.Id(t.ID).AddUpdate(&sl_update, nil)
	if err != nil {
		return &TicketUpdate{}, err
	}
	t.SyncTicketUpdates(user)
	return PickUpdate(*upds[0].Id)
}

// SyncTickets() creates and updates all Tickets of given user's account.
func SyncTickets(user *User) (count int, err error) {
	log.Infof("sync tickets... (use %v)", user.Username)
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"

	account := user.Account()
	if account == nil {
		log.Errorf("account link broken! %v of %v", user.ID, user.AccountId)
		return 0, errors.New("account link broken!")
	}
	log.Debugf("account: %v", account)

	date_since := account.LastBatch.AddDate(0, 0, -1).
		Format("01/02/2006 15:04:05")
	log.Infof("try to sync tickets from %v...", date_since)

	service := services.GetAccountService(sess)
	data, err := service.
		Mask("id;accountId;assignedUserId;subjectId;groupId;statusId;title;totalUpdateCount;createDate;lastEditDate;lastEditType").
		Filter(filter.Build(
			filter.Path("tickets.lastEditDate").DateAfter(date_since),
		)).
		GetTickets()
	if err != nil {
		log.Errorf("slapi error: %v", err)
		return 0, err
	}

	count = 0
	for _, el := range data {
		ticket := &Ticket{}
		copier.Copy(ticket, el)
		ticket.ID = *el.Id
		ticket.CreateDate, _ = time.Parse(time.RFC3339, el.CreateDate.String())
		ticket.LastEditDate, _ = time.Parse(time.RFC3339, el.LastEditDate.String())
		log.Debugf("ticket %v/%v --", ticket.AccountId, ticket.ID)
		for _, elu := range el.Updates {
			ticket_update := &TicketUpdate{}
			copier.Copy(ticket_update, elu)
			log.Debugf("--- %v ---", ticket_update)
		}

		err = ticket.Save()
		if err != nil {
			log.Errorf("cannot create ticket: %v, %v", err, ticket)
		} else {
			count++
		}
	}
	if len(data) == count {
		log.Infof("Bingo! all data were inserted to database! (%v)", count)
		account.LastBatch = time.Now()
		account.Save()
	} else {
		log.Errorf("Oops! some data not inserted! %v of %v saved.",
			count, len(data))
	}
	return count, nil
}

// SyncTicketUpdates() creates and updates all Updates of Ticket instance.
func (t *Ticket) SyncTicketUpdates(user *User) (count int, err error) {
	log.Infof("sync ticket updates... (use %v)", user.Username)
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"

	date_since := t.LastSync.AddDate(0, 0, -1).Format("01/02/2006 15:04:05")
	log.Debugf("try to sync updates from %v...", date_since)

	data, err := services.GetTicketService(sess).
		Id(t.ID).
		Mask("id;ticketId;editorId;editorType;entry;createDate;type").
		Filter(filter.Build(
			filter.Path("updates.createDate").DateAfter(date_since),
		)).
		GetUpdates()
	if err != nil {
		log.Errorf("slapi error: %v", err)
		return 0, err
	}

	count = 0
	errors := 0
	exists := 0
	for _, el := range data {
		update := &TicketUpdate{}
		copier.Copy(update, el)
		update.ID = *el.Id
		update.CreateDate, _ = time.Parse(time.RFC3339, el.CreateDate.String())
		log.Debugf("%v/%v --", update.TicketId, update.ID)
		if ok, _ := DB.Where("id=?", update.ID).Exists(update); ok {
			log.Debugf("update %v already exists!", update.ID)
			exists++
		} else {
			err = DB.Create(update)
			if err != nil {
				log.Errorf("cannot create update: %v, %v", err, update)
				errors++
			} else {
				log.Debugf("update %v created.", update.ID)
				count++
			}
		}
	}
	if len(data) == (count + exists) {
		log.Infof("Yo! all were inserted or up-to-date! (%v,%v)", count, exists)
		t.LastSync = time.Now()
		err = t.Save()
		if err != nil {
			log.Errorf("...but cannot save last sync time of ticket: %v", err)
		}
	} else {
		log.Errorf("there are %v errors out of total %v. skip time update.",
			errors, len(data))
	}
	return count, nil
}

//// search functions:

// Find a Ticket with ticket_id.
func FindTicket(ticket_id int) (ticket *Ticket, err error) {
	ticket = &Ticket{}
	err = DB.Find(ticket, ticket_id)
	return
}

func PickUpdate(update_id int) (update *TicketUpdate, err error) {
	update = &TicketUpdate{}
	err = DB.Find(update, update_id)
	return
}

//// Association and Relationship based search for instances.

// Updates() returns all updates for the ticket.
func (t *Ticket) Updates() (updates *TicketUpdates) {
	updates = &TicketUpdates{}
	err := DB.BelongsTo(t).
		Order("create_date desc").
		All(updates)
	if err != nil {
		return nil
	}
	inspect("'updates' from dbms", updates)
	return
}

func (t *Ticket) FirstUpdate() (update *TicketUpdate) {
	update = &TicketUpdate{}
	err := DB.BelongsTo(t).
		Order("create_date").
		First(update)
	if err != nil {
		return nil
	}
	inspect("'update' from dbms", update)
	return
}

//// DBMS Functions:

// Save() saves the Ticket instance. (create or update)
func (t *Ticket) Save() (err error) {
	old := &Ticket{}
	err = DB.Find(old, t.ID)
	origin, _ := time.Parse("2006-01-02", "1977-05-25")
	if err == nil {
		if t.LastSync.Before(origin) {
			log.Debugf("preserve old timestamp!")
			t.LastSync = old.LastSync
		}
		verrs, err := DB.ValidateAndUpdate(t)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return verrs
		}
	} else {
		lst, e := time.Parse(time.RFC3339, "1977-05-25T00:00:00+09:00")
		if e == nil {
			t.LastSync = lst
		} else {
			t.LastSync = time.Now()
		}

		verrs, err := DB.ValidateAndCreate(t)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return verrs
		}
	}
	return nil
}

// display helpers:

func (t Ticket) Account() interface{} {
	account := &Account{}
	DB.Find(account, t.AccountId)
	return account
}

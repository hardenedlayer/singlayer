package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/copier"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type Ticket struct {
	ID             int       `json:"id" db:"id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	AccountId      int       `json:"account_id" db:"account_id"`
	AssignedUserId int       `json:"assigned_user_id" db:"assigned_user_id"`
	SubjectId      int       `json:"subject_id" db:"subject_id"`
	GroupId        int       `json:"group_id" db:"group_id"`
	StatusId       int       `json:"status_id" db:"status_id"`
	Title          string    `json:"title" db:"title"`
	CreateDate     time.Time `json:"create_date" db:"create_date"`
	LastEditDate   time.Time `json:"last_edit_date" db:"last_edit_date"`
	LastEditType   string    `json:"last_edit_type" db:"last_edit_type"`
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

// SyncTickets() creates and updates all Tickets of given user's account.
func SyncTickets(user *User) (count int, err error) {
	Logger.Printf("sync tickets... (use %v)", user.Username)
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"

	account := user.Account()
	Logger.Printf("account: %v", account)

	date_since := account.LastBatch.Format("01/02/2006 15:04:05")
	Logger.Printf("try to sync tickets from %v...", date_since)

	service := services.GetAccountService(sess)
	data, err := service.
		Mask("id;accountId;assignedUserId;subjectId;groupId;statusId,title;totalUpdateCount;createDate;lastEditDate;lastEditType;updates.id;updates.ticketId;updates.editorType;updates.editorId").
		Filter(`{"tickets":{"lastEditDate":{"operation":"greaterThanDate","options":[{"name":"date","value":["` + date_since + `"]}]}}}`).
		GetTickets()
	if err != nil {
		Logger.Printf("slapi error: %v", err)
		return 0, err
	}

	count = 0
	errors := 0
	exists := 0
	for _, el := range data {
		ticket := &Ticket{}
		copier.Copy(ticket, el)
		ticket.ID = *el.Id
		ticket.CreateDate,_ = time.Parse(time.RFC3339, el.CreateDate.String())
		ticket.LastEditDate,_ = time.Parse(time.RFC3339, el.LastEditDate.String())
		Logger.Printf("ticket %v/%v --", ticket.AccountId, ticket.ID)
		for _, elu := range el.Updates {
			ticket_update := &TicketUpdate{}
			copier.Copy(ticket_update, elu)
			Logger.Printf("--- %v ---", ticket_update)
		}

		if ok, _ := DB.Where("id=?", ticket.ID).Exists(ticket); ok {
			Logger.Printf("ticket %v already exists!", ticket.ID)
			exists++
		} else {
			err = DB.Create(ticket)
			if err != nil {
				Logger.Printf("cannot create ticket: %v, %v", err, ticket)
				errors++
			} else {
				Logger.Printf("ticket %v created.", ticket.ID)
				count++
			}
		}
	}
	if len(data) == count {
		Logger.Printf("Bingo! all data were inserted to database! (%v)", count)
		account.LastBatch = time.Now()
		account.Save()
	} else {
		Logger.Printf("Oops! some data not inserted! x:%v, s:%v, e:%vi (%v)",
			exists, count, errors, len(data))
	}
	return count, nil
}

// SyncTicketUpdates() creates and updates all Updates of Ticket instance.
func (t *Ticket) SyncTicketUpdates(user *User) (count int, err error) {
	Logger.Printf("sync ticket updates... (use %v)", user.Username)
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"
	data, err := services.GetTicketService(sess).
		Id(t.ID).
		Mask("id;ticketId;editorId;editorType;entry;createDate;type").
		GetUpdates()
	if err != nil {
		Logger.Printf("slapi error: %v", err)
		return 0, err
	}

	count = 0
	errors := 0
	exists := 0
	for _, el := range data {
		update := &TicketUpdate{}
		copier.Copy(update, el)
		update.ID = *el.Id
		update.CreateDate,_ = time.Parse(time.RFC3339, el.CreateDate.String())
		Logger.Printf("%v/%v --", update.TicketId, update.ID)
		if ok, _ := DB.Where("id=?", update.ID).Exists(update); ok {
			Logger.Printf("update %v already exists!", update.ID)
			exists++
		} else {
			err = DB.Create(update)
			if err != nil {
				Logger.Printf("cannot create update: %v, %v", err, update)
				errors++
			} else {
				Logger.Printf("update %v created.", update.ID)
				count++
			}
		}
	}
	if len(data) == count {
		Logger.Printf("Bingo! all data were inserted to database! (%v)", count)
	} else {
		Logger.Printf("Oops! some data not inserted! x:%v, s:%v, e:%vi (%v)",
			exists, count, errors, len(data))
	}
	return count, nil
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
	inspect("updates from dbms", updates)
	return
}

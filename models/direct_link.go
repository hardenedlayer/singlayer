package models

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type DirectLink struct {
	ID            uuid.UUID `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	SingleID      uuid.UUID `json:"single_id" db:"single_id"`
	UserId        int       `json:"user_id" db:"user_id"`
	AccountId     int       `json:"account_id" db:"account_id"`
	SiblingID     uuid.UUID `json:"sibling_id" db:"sibling_id"`
	VlanId        int       `json:"vlan_id" db:"vlan_id"`
	TicketId      int       `json:"ticket_id" db:"ticket_id"`
	Type          string    `json:"type" db:"type"`
	Location      string    `json:"location" db:"location"`
	LineNumber    int       `json:"line_number" db:"line_number"`
	Port          string    `json:"port" db:"port"`
	Router        int       `json:"router" db:"router"`
	Speed         int       `json:"speed" db:"speed"`
	RoutingOption string    `json:"routing_option" db:"routing_option"`
	MultiPath     bool      `json:"multi_path" db:"multi_path"`
	ASN           int       `json:"asn" db:"asn"`
	Prefix        int       `json:"prefix" db:"prefix"`
	Migration     string    `json:"migration" db:"migration"`
	Comments      string    `json:"comments" db:"comments"`
	Signature     string    `json:"signature" db:"signature"`
	XCRIP         string    `json:"xcr_ip" db:"xcr_ip"`
	CERIP         string    `json:"cer_ip" db:"cer_ip"`
	Status        string    `json:"status" db:"status"`
	Notes         string    `json:"notes" db:"notes"`
}

func (d DirectLink) String() string {
	str := fmt.Sprintf("%v %vGbps %v Line#%v of %v",
		d.Type, d.Speed, d.RoutingOption, d.LineNumber, d.AccountNick())
	return str
}

func (d DirectLink) Marshal() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

func (d DirectLink) Hash() string {
	str := fmt.Sprintf("%v-%v-%v-%v",
		d.SingleID, d.AccountId, d.UserId, d.String())
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}

func (d DirectLink) HashShort() string {
	return d.Signature[0:10]
}

type DirectLinks []DirectLink

func (d DirectLinks) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (d *DirectLink) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: d.UserId, Name: "UserId"},
		&validators.IntIsPresent{Field: d.AccountId, Name: "AccountId"},
		&validators.IntIsPresent{Field: d.VlanId, Name: "VlanId"},
		&validators.StringIsPresent{Field: d.Type, Name: "Type"},
		&validators.StringIsPresent{Field: d.Location, Name: "Location"},
		&validators.IntIsPresent{Field: d.LineNumber, Name: "LineNumber"},
		&validators.StringIsPresent{Field: d.Port, Name: "Port"},
		&validators.IntIsPresent{Field: d.Router, Name: "Router"},
		&validators.IntIsPresent{Field: d.Speed, Name: "Speed"},
		&validators.StringIsPresent{Field: d.RoutingOption, Name: "RoutingOption"},
		&validators.IntIsPresent{Field: d.Prefix, Name: "Prefix"},
		&validators.StringIsPresent{Field: d.Migration, Name: "Migration"},
		&validators.StringIsPresent{Field: d.Signature, Name: "Signature"},
		&validators.StringIsPresent{Field: d.Status, Name: "Status"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (d *DirectLink) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (d *DirectLink) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	if d.Status == "draft" && d.TicketId > 0 {
		log.Infof("status is draft but has ticket. upgrading to ordered...")
		d.Status = "ordered"
	}
	if d.Status == "ordered" && d.XCRIP != "" && d.CERIP != "" {
		log.Infof("status is ordered but has ips. upgrading to configured...")
		d.Status = "configured"
	}
	return validate.NewErrors(), nil
}

//// selectors:

func PickDirectLink(directlink_id uuid.UUID) (dlink *DirectLink) {
	dlink = &DirectLink{}
	err := DB.Find(dlink, directlink_id)
	if err != nil {
		return nil
	}
	return
}

//// Association and Relationship based search for instances.
//// It need instance of Single so more expensive than raw query. FIXME later.

// Progresses() returns instance of Progresses struct.
func (d *DirectLink) Progresses() (progresses *Progresses) {
	progresses = &Progresses{}
	err := DB.BelongsTo(d).
		Order("created_at").
		All(progresses)
	if err != nil {
		return nil
	}
	return
}

func (d *DirectLink) Updates() (updates *TicketUpdates) {
	updates = &TicketUpdates{}
	err := DB.Where("ticket_id = ?", d.TicketId).
		Order("create_date desc").
		All(updates)
	if err != nil {
		return
	}
	return
}

func (d *DirectLink) Ticket() (ticket *Ticket) {
	ticket = &Ticket{}
	err := DB.Find(ticket, d.TicketId)
	if err != nil {
		return nil
	}
	return
}

func (d *DirectLink) Account() (account *Account) {
	account = &Account{}
	err := DB.Find(account, d.AccountId)
	if err != nil {
		return nil
	}
	return
}

//// display helpers:
func (d DirectLink) AccountName() interface{} {
	return d.Account().CompanyName
}

func (d DirectLink) AccountNick() interface{} {
	return d.Account().String()
}

func (d DirectLink) UserName() interface{} {
	user, err := FindUser(d.UserId)
	if err != nil {
		return d.UserId
	}
	return user.Username
}

func (d DirectLink) SingleName() interface{} {
	single, err := FindSingle(d.SingleID)
	if err != nil {
		return "Unknown!!!"
	}
	return single.Name
}

func (d DirectLink) PairTicketId() interface{} {
	plink := PickDirectLink(d.SiblingID)
	if plink != nil {
		return plink.TicketId
	}
	return nil
}

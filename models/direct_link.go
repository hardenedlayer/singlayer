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
	Router        string    `json:"router" db:"router"`
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
	str := fmt.Sprintf("%v %vGbps %v Line#%v",
		d.Type, d.Speed, d.RoutingOption, d.LineNumber)
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
		&validators.StringIsPresent{Field: d.Router, Name: "Router"},
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
	return validate.NewErrors(), nil
}

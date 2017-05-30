package models

import (
	"encoding/json"
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
	UserID        int       `json:"user_id" db:"user_id"`
	AccountID     int       `json:"account_id" db:"account_id"`
	SiblingID     uuid.UUID `json:"sibling_id" db:"sibling_id"`
	VlanID        int       `json:"vlan_id" db:"vlan_id"`
	TicketID      int       `json:"ticket_id" db:"ticket_id"`
	Type          string    `json:"type" db:"type"`
	Location      string    `json:"location" db:"location"`
	LinuNumber    int       `json:"linu_number" db:"linu_number"`
	Port          string    `json:"port" db:"port"`
	Router        string    `json:"router" db:"router"`
	Speed         int       `json:"speed" db:"speed"`
	RoutingOption string    `json:"routing_option" db:"routing_option"`
	MultiPath     bool      `json:"multi_path" db:"multi_path"`
	Asn           unit32    `json:"asn" db:"asn"`
	Prefix        int       `json:"prefix" db:"prefix"`
	Migration     string    `json:"migration" db:"migration"`
	Comments      string    `json:"comments" db:"comments"`
	Signature     string    `json:"signature" db:"signature"`
	XcrIp         string    `json:"xcr_ip" db:"xcr_ip"`
	CerIp         string    `json:"cer_ip" db:"cer_ip"`
	Status        string    `json:"status" db:"status"`
	Notes         string    `json:"notes" db:"notes"`
}

// String is not required by pop and may be deleted
func (d DirectLink) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// DirectLinks is not required by pop and may be deleted
type DirectLinks []DirectLink

// String is not required by pop and may be deleted
func (d DirectLinks) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (d *DirectLink) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: d.UserID, Name: "UserID"},
		&validators.IntIsPresent{Field: d.AccountID, Name: "AccountID"},
		&validators.IntIsPresent{Field: d.VlanID, Name: "VlanID"},
		&validators.IntIsPresent{Field: d.TicketID, Name: "TicketID"},
		&validators.StringIsPresent{Field: d.Type, Name: "Type"},
		&validators.StringIsPresent{Field: d.Location, Name: "Location"},
		&validators.IntIsPresent{Field: d.LinuNumber, Name: "LinuNumber"},
		&validators.StringIsPresent{Field: d.Port, Name: "Port"},
		&validators.StringIsPresent{Field: d.Router, Name: "Router"},
		&validators.IntIsPresent{Field: d.Speed, Name: "Speed"},
		&validators.StringIsPresent{Field: d.RoutingOption, Name: "RoutingOption"},
		&validators.IntIsPresent{Field: d.Prefix, Name: "Prefix"},
		&validators.StringIsPresent{Field: d.Migration, Name: "Migration"},
		&validators.StringIsPresent{Field: d.Comments, Name: "Comments"},
		&validators.StringIsPresent{Field: d.Signature, Name: "Signature"},
		&validators.StringIsPresent{Field: d.XcrIp, Name: "XcrIp"},
		&validators.StringIsPresent{Field: d.CerIp, Name: "CerIp"},
		&validators.StringIsPresent{Field: d.Status, Name: "Status"},
		&validators.StringIsPresent{Field: d.Notes, Name: "Notes"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (d *DirectLink) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (d *DirectLink) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

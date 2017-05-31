package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type Vlan struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	R1LinkID  uuid.UUID `json:"r1_link_id" db:"r1_link_id"`
	R2LinkID  uuid.UUID `json:"r2_link_id" db:"r2_link_id"`
	AccountId int       `json:"account_id" db:"account_id"`
	Booked    bool      `json:"booked" db:"booked"`
}

func (v Vlan) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

type Vlans []Vlan

func (v Vlans) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (v *Vlan) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: v.ID, Name: "ID"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (v *Vlan) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (v *Vlan) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

//// selectors:
func VLAN(vlan_id int) (vlan *Vlan) {
	vlan = &Vlan{}
	err := DB.Find(vlan, vlan_id)
	if err != nil {
		log.Errorf("vlan query error: %v", err)
	}
	return vlan
}

func NextVLAN(account_id int) (vlan *Vlan) {
	d, _ := time.ParseDuration("-5m")
	vlans := &Vlans{}
	t := time.Now().Add(d).UTC().Format(time.RFC3339)
	err := DB.
		Where("booked=? AND updated_at < ?", true, t).
		Where("r1_link_id=?", uuid.UUID{}).
		Where("r2_link_id=?", uuid.UUID{}).
		All(vlans)
	log.Debugf("remove locks for %v entries... %v", len(*vlans), vlans)
	for _, v := range *vlans {
		v.AccountId = 0
		v.Booked = false
		DB.Save(&v)
	}

	vlan = &Vlan{}
	err = DB.
		Where("booked=?", false).
		Where("account_id=?", 0).
		Where("r1_link_id=?", uuid.UUID{}).
		Where("r2_link_id=?", uuid.UUID{}).
		Order("id").First(vlan)
	vlan.AccountId = account_id
	vlan.Booked = true
	DB.Save(vlan)
	if err != nil {
		log.Errorf("next_vlan query error: %v", err)
	}
	return vlan
}

func NextRouter() int {
	slot_r1, _ := DB.Where("r1_link_id=?", uuid.UUID{}).Count(&Vlans{})
	slot_r2, _ := DB.Where("r2_link_id=?", uuid.UUID{}).Count(&Vlans{})
	log.Infof("slot count: r1 %v, r2 %v", slot_r1, slot_r2)
	if slot_r1 < slot_r2 {
		return 2
	} else {
		return 1
	}
}

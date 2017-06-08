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

type Account struct {
	ID          int       `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	BrandId     int       `json:"brand_id" db:"brand_id"`
	CompanyName string    `json:"company_name" db:"company_name"`
	Email       string    `json:"email" db:"email"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	LastBatch   time.Time `json:"last_batch" db:"last_batch"`
	NickName    string    `json:"nick_name" db:"nick_name"`
}

func (a Account) String() string {
	if a.NickName != "" {
		return a.NickName
	}
	return a.CompanyName
}

func (a Account) Marshal() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

type Accounts []Account

func (a Accounts) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (a *Account) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: a.ID, Name: "ID"},
		&validators.StringIsPresent{Field: a.CompanyName, Name: "CompanyName"},
		&validators.TimeIsPresent{Field: a.LastBatch, Name: "LastBatch"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (a *Account) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (a *Account) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

//// Backend API Call:

// UpdateAndSave() fills up Account data from api call result and saves it.
func (a *Account) UpdateAndSave(user *User) (err error) {
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"
	service := services.GetAccountService(sess)
	sl_acc, err := service.
		Mask("id;brandId;companyName;email;firstName;lastName").
		GetObject()
	if err != nil {
		log.Errorf("softlayer api exception: %v --", err)
		return err
	}
	copier.Copy(a, sl_acc)
	inspect("updated account", a)
	return a.Save()
}

//// DBMS Functions:

// Save() saves the Account instance. (create or update)
func (a *Account) Save() (err error) {
	old := &Account{}
	err = DB.Find(old, a.ID)
	origin, _ := time.Parse("2006-01-02", "1977-05-25")
	if err == nil {
		if a.LastBatch.Before(origin) {
			log.Debugf("preserve old timestamp!")
			a.LastBatch = old.LastBatch
		}
		verrs, err := DB.ValidateAndUpdate(a)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return verrs
		}
	} else {
		t, e := time.Parse(time.RFC3339, "1977-05-25T00:00:00+09:00")
		if e == nil {
			a.LastBatch = t
		} else {
			a.LastBatch = time.Now()
		}

		verrs, err := DB.ValidateAndCreate(a)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return verrs
		}
	}
	return nil
}

//// Association and Relationship based search for instances.
//// It need instance of Single so more expensive than raw query. FIXME later.

// DirectLinks() returns all directlinks associated to the account.
func (a *Account) DirectLinks() (dlinks *DirectLinks) {
	dlinks = &DirectLinks{}
	err := DB.BelongsTo(a).All(dlinks)
	if err != nil {
		return nil
	}
	return
}

// display functions:

func (a Account) Contact() interface{} {
	return a.FirstName + " " + a.LastName + " <" + a.Email + ">"
}

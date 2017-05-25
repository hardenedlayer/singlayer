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
}

// String is not required by pop and may be deleted
func (a Account) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Accounts is not required by pop and may be deleted
type Accounts []Account

// String is not required by pop and may be deleted
func (a Accounts) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (a *Account) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: a.ID, Name: "ID"},
		&validators.StringIsPresent{Field: a.CompanyName, Name: "CompanyName"},
		&validators.TimeIsPresent{Field: a.LastBatch, Name: "LastBatch"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (a *Account) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (a *Account) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// UpdateAndSave
func (a *Account) UpdateAndSave(user *User) (err error) {
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"
	service := services.GetAccountService(sess)
	sl_acc, err := service.
		Mask("id;brandId;companyName;email;firstName;lastName").
		GetObject()
	if err != nil {
		Logger.Printf("softlayer api exception: %v --", err)
		return err
	}
	copier.Copy(a, sl_acc)
	a.LastBatch = time.Now()

	inspect("updated account", a)
	return a.Save()
}

func (a *Account) Save() (err error) {
	old := &Account{}
	err = DB.Find(old, a.ID)
	if err == nil {
		verrs, err := DB.ValidateAndUpdate(a)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return verrs
		}
	} else {
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

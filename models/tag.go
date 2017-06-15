package models

import (
	"encoding/json"
	"time"
)

type Tag struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	AccountId int       `json:"account_id" db:"account_id"`
	Name      string    `json:"name" db:"name"`
}

func (t Tag) String() string {
	jc, _ := json.Marshal(t)
	return string(jc)
}

type Tags []Tag

func (t Tags) String() string {
	jc, _ := json.Marshal(t)
	return string(jc)
}

//// database fuctions:

func PickTag(id int) (t *Tag) {
	t = &Tag{}
	err := DB.Find(t, id)
	if err != nil {
		return nil
	}
	return
}

// Save()
func (t *Tag) Save() error {
	old := PickTag(t.ID)
	if old == nil {
		verrs, err := DB.ValidateAndCreate(t)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return verrs
		}
	} else {
		log.Debugf("saving existing %v...", t)
		verrs, err := DB.ValidateAndUpdate(t)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return verrs
		}
	}
	return nil
}

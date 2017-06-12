package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type Doc struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	SingleID  uuid.UUID `json:"single_id" db:"single_id"`
	Type      string    `json:"type" db:"type"`
	Category  string    `json:"category" db:"category"`
	Subject   string    `json:"subject" db:"subject"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Published bool      `json:"published" db:"published"`
}

func (d Doc) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

type Docs []Doc

// additional types for manage
type Category struct {
	Name string `json:"name" db:"category"`
}

type Categories []Category

type Subject struct {
	Category string `json:"category" db:"category"`
	Name     string `json:"name" db:"subject"`
}

type Subjects []Subject

func (d Docs) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (d *Doc) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: d.Type, Name: "Type"},
		&validators.StringIsPresent{Field: d.Category, Name: "Category"},
		&validators.StringIsPresent{Field: d.Subject, Name: "Subject"},
		&validators.StringIsPresent{Field: d.Content, Name: "Content"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (d *Doc) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (d *Doc) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func DocCategories() interface{} {
	cats := &Categories{}
	DB.RawQuery(`SELECT DISTINCT category
		FROM docs
		ORDER BY category`).
		All(cats)
	inspect("categories", cats)
	return cats
}

func DocSubjects() interface{} {
	subjects := &Subjects{}
	DB.RawQuery(`SELECT DISTINCT category, subject
		FROM docs
		ORDER BY category, subject`).
		All(subjects)
	inspect("subjects", subjects)
	return subjects
}

func (d Doc) Author() interface{} {
	single := &Single{}
	DB.Find(single, d.SingleID)
	return single.Name
}

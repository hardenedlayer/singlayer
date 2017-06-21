package models

import (
	"encoding/json"
	"time"

	"github.com/satori/go.uuid"
)

type Log struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	SingleID  uuid.UUID `json:"single_id" db:"single_id"`
	Category  string    `json:"category" db:"category"`
	Level     string    `json:"level" db:"level"`
	Message   string    `json:"message" db:"message"`
	IsRead    bool      `json:"is_read" db:"is_read"`
}

func (a Log) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

type Logs []Log

func (a Logs) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

//// display helpers:

func (a Log) Single() (single *Single) {
	single = &Single{}
	DB.Find(single, a.SingleID)
	return
}

func l(category, level, message string) error {
	log := &Log{
		Category: category,
		Level:    level,
		Message:  message,
	}
	err := DB.Create(log)
	if err != nil {
		return err
	}
	return nil
}

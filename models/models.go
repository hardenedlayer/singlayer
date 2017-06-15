package models

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gobuffalo/envy"
	"github.com/markbates/pop"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection
var log = logrus.New()
var is_test = false

var mail_sender string
var mail_admins []string

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"

	if env == "development" {
		log.Formatter = &logrus.TextFormatter{}
		log.Out = os.Stdout
		log.Level = logrus.DebugLevel
		is_test = true
	}

	mail_sender = os.Getenv("MAIL_SENDER")
	log.Infof("mail_sender: %v", mail_sender)
	if mail_sender == "" {
		log.Fatal("environment variable MAIL_SENDER not defined!")
	}

	// deprecated block. currently use admin mail. (AdminMail)
	admins := os.Getenv("MAIL_ADMINS")
	if len(admins) > 0 {
		for _, el := range strings.Split(admins, ";") {
			mail_admins = append(mail_admins, strings.TrimSpace(el))
		}
	}
	log.Infof("mail_admins: %v (deprecated)", mail_admins)

	pop.MapTableName("TicketStatus", "ticket_statuses")
	pop.MapTableName("TicketStatuses", "ticket_statuses")
}

//// Additional functions for my debugging.

// inspect: to check data type and value.
func inspect(desc string, data interface{}) {
	log.Debugf("{\"description\":\"%s\", \"datatype\":\"%T\", \"data\":%v}",
		desc, data, toJSON(data))
}

// toJSON returns JSON formatted string with given data.
func toJSON(data interface{}) string {
	str, _ := json.Marshal(data)
	return string(str)
}

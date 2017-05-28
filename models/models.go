package models

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gobuffalo/envy"
	"github.com/markbates/pop"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection
var Logger *log.Logger

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"

	Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	pop.MapTableName("TicketStatus", "ticket_statuses")
	pop.MapTableName("TicketStatuses", "ticket_statuses")
}

//// Additional functions for my debugging.

// inspect: to check data type and value.
func inspect(desc string, data interface{}) {
	Logger.Printf("{\"description\":\"%s\", \"datatype\":\"%T\", \"data\":%v}",
		desc, data, toJSON(data))
}

// toJSON returns JSON formatted string with given data.
func toJSON(data interface{}) string {
	str, _ := json.Marshal(data)
	return string(str)
}

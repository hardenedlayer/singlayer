package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/envy"
	"github.com/hardenedlayer/singlayer/actions"
)

func main() {
	port := envy.Get("PORT", "3000")
	log.Printf("Starting singlayer on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), actions.App()))
}

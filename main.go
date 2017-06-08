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
	use_tls := envy.Get("USE_TLS", "no")
	log.Printf("Starting singlayer on port %s\n", port)
	if use_tls == "yes" {
		tls_cert := envy.Get("TLS_CERT", "cert.pem")
		tls_key := envy.Get("TLS_KEY", "key.pem")
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%s", port),
			tls_cert, tls_key, actions.App()))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), actions.App()))
	}
}

// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"log"
	"os"

	"github.com/lacolaco/activitypub.lacolaco.net/web"
)

func main() {
	log.Print("starting server...")

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := web.Start(port); err != nil {
		log.Fatal(err)
	}
}

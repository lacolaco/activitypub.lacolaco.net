// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/web"
)

func main() {
	godotenv.Load()
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if err := web.Start(conf); err != nil {
		log.Fatal(err)
	}
}

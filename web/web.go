package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func Start(port string) error {
	http.HandleFunc("/", handler)

	// Start HTTP server.
	log.Printf("listening on http://localhost:%s", port)
	return http.ListenAndServe(":"+port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

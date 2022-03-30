package main

import (
	"log"
	"net/http"

	"github.com/I70l0teN4ik/go-rest/src/app"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", app.NewServer()))
}

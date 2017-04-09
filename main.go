package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	PORT := ":8080"
	r := httprouter.New()

	log.Print("Running server on " + PORT)
	log.Fatal(http.ListenAndServe(PORT, r))
}

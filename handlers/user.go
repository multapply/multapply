package handlers

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetUser - Returns user with given :id
func GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Print("Hello, there!")
}

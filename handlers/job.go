package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/multapply/multapply/models"
)

// CreateJob - Handler for creating a new Job - POST /jobs
func (env *Env) CreateJob(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	j := new(models.Job)
	err := json.NewDecoder(r.Body).Decode(j)
	if err != nil {
		http.Error(w, "Error parsing request", 400)
		return
	}

	// TODO: Define a trim func in models/job to make sure required fields are non-empty

	// TEMP: For now, just insert into DB
	err = models.InsertJob(env.DB, j)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to insert Job", 500)
		return
	}

	// TODO: Proper success response
}

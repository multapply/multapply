package handlers

import (
	"encoding/json"
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
}

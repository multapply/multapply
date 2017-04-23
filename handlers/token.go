package handlers

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/multapply/multapply/models"
	"github.com/multapply/multapply/util/constants"
)

// RequestNewToken - Handler for requesting a new access token
func (env *Env) RequestNewToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check context for tid (token_id)
	ctx := r.Context()
	tid := ctx.Value(constants.ContextKeyTokenID)
	if tid == nil {
		http.Error(w, "Missing claims", 401)
		return
	}

	tokenID := int(tid.(float64))

	u, err := models.GetUserByTokenID(env.DB, tokenID)
	if err != nil {
		http.Error(w, "Error requesting new access token", 500)
		return
	} else if u == nil {
		http.Error(w, "Error retrieving user info", 500)
		return
	}

	accessToken, err := models.GetAccessToken(u)
	if err != nil {
		http.Error(w, "Error creating new access token", 500)
		return
	}

	// store access token in a cookie for the client
	accessCookie := http.Cookie{
		Name:  "mAT",
		Value: accessToken,
		//Domain:   ".multapply.io", // TODO: When I have a domain, uncomment this
		Expires:  time.Now().Add(time.Minute * 30), // TODO: Use constant for this
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &accessCookie)

	// TODO: Return success
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/julienschmidt/httprouter"
	"github.com/multapply/multapply/models"
)

// GetUser - Returns user with given :id
// TODO: Actually implement this - currently just a placeholder to test auth flow
func (env *Env) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("Auth successful!"))
}

// CreateUser - Attempts to create a new User with the given
// parameters in the request body
// TODO: All errors that don't throw http.Error SHOULD, with code 500??
// TODO: Return actual JSON message after errors instead of just returning
// TODO: Split a lot of functionality into helper functions
func (env *Env) CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request body into user struct
	n := new(models.NewUser)
	err := json.NewDecoder(r.Body).Decode(n)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer r.Body.Close()

	// Trim all fields of any leading/trailing whitespace
	// and check for nonempty input
	n.Trim()
	if len(n.Username) == 0 {
		http.Error(w, "Username required", 400)
		return
	} else if len(n.Email) == 0 {
		http.Error(w, "Email address required", 400)
		return
	} else if len(n.Password) == 0 {
		http.Error(w, "Password required", 400)
		return
	} else if len(n.ConfirmPassword) == 0 {
		http.Error(w, "Password confirmation required", 400)
		return
	}

	// Check that username is unique
	// TODO: Make usernames case sensitive, i.e. store UPPER/lower of usernames/emails
	var exists bool
	exists, err = models.UsernameExists(env.DB, n.Username)
	if err != nil {
		log.Fatal(err)
	} else if exists {
		http.Error(w, "Username is taken", 409)
		return
	}

	// Check for valid email address
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, n.Email); !m {
		http.Error(w, "Invalid email address", 400)
		return
	}

	// Check that email address is unique
	exists, err = models.EmailExists(env.DB, n.Email)
	if err != nil {
		log.Fatal(err)
	} else if exists {
		http.Error(w, "Email address is taken", 409)
		return
	}

	// Make sure password and confirmation are equal
	// and if they are, hash the password
	// TODO: Add secret to password THEN hash
	if n.Password != n.ConfirmPassword {
		http.Error(w, "Password and password confirmation do not match", 400)
		return
	}
	var passwordHash []byte
	passwordHash, err = bcrypt.GenerateFromPassword([]byte(n.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	n.PasswordHash = string(passwordHash)

	// Create User from NewUser n then insert into DB
	u := models.CreateUser(n)
	uid, err := models.InsertUser(env.DB, u)
	if err != nil {
		http.Error(w, "Failed to create account", 500)
	}
	u.UserID = uid

	// Create access token
	var accessToken string
	accessToken, err = models.GetAccessToken(u)
	if err != nil {
		models.RemoveNewUser(env.DB, u.Username)
		http.Error(w, "Failed to create account", 500)
		log.Fatal("Could not create access token.")
	}

	// store access token in a cookie for the client
	accessCookie := http.Cookie{
		Name:  "mAT",
		Value: accessToken,
		//Domain:   ".multapply.io",
		Expires:  time.Now().Add(time.Minute * 30), // TODO: Use constant for this
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &accessCookie)

	// Create refresh token
	var refreshToken string
	refreshToken, err = models.GetRefreshToken(env.DB, u)
	if err != nil {
		models.RemoveNewUser(env.DB, u.Username)
		http.Error(w, "Failed to create account", 500)
		log.Fatal("Could not create refresh token:", err)
	}

	// store refresh token in cookie for client
	refreshCookie := http.Cookie{
		Name:  "mRT",
		Value: refreshToken,
		//Domain: ".multapply.io",
		Expires:  time.Now().Add(time.Minute * 525600), // TODO: Use constant for this
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &refreshCookie)
}

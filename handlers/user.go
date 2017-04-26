package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/julienschmidt/httprouter"
	"github.com/multapply/multapply/models"
)

// GetUser - Returns user with given :id
// TODO: Actually implement this - currently just a placeholder to test auth flow
func (env *Env) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	log.Println("IP:", r.RemoteAddr)
	w.Write([]byte("Auth successful!"))
	r.Body.Close()
}

// CreateUser - Attempts to create a new User with the given parameters in the request body
// TODO: All errors that don't throw http.Error SHOULD, with code 500??
// TODO: Return actual JSON message after errors instead of just returning
// TODO: Split a lot of functionality into helper functions
func (env *Env) CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	// Decode request body into user struct
	n := new(models.NewUser)
	err := json.NewDecoder(r.Body).Decode(n)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

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
		http.Error(w, "Error creating account.", 500)
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
		//Domain:   ".multapply.io", // TODO: When I have a domain, uncomment this
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
		//Domain: ".multapply.io", // TODO: When domain is hooked up, uncomment this
		Expires:  time.Now().Add(time.Minute * 525600), // TODO: Use constant for this
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &refreshCookie)
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginUser - Attempts to log in user with given credentials in request body
// If successful, returns access/refresh tokens to user
// TODO: Rate limit amount of login attempts to prevent attackers from brute forcing login
func (env *Env) LoginUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	// Decode request body into user struct
	l := new(loginRequest)
	err := json.NewDecoder(r.Body).Decode(l)
	if err != nil {
		http.Error(w, "Error logging in", 500)
		return
	}

	// Trim all fields of any leading/trailing whitespace
	// and check for nonempty input
	l.Username = strings.TrimSpace(l.Username)
	l.Password = strings.TrimSpace(l.Password)
	if len(l.Username) == 0 {
		http.Error(w, "Username missing", 400) // TODO: Replace all numbers with actual error constants
		return
	} else if len(l.Password) == 0 {
		http.Error(w, "Password missing", 400)
		return
	}

	// Check if user exists
	var exists bool
	exists, err = models.UsernameExists(env.DB, l.Username)
	if err != nil {
		http.Error(w, "Error logging in.", 500)
		return
	} else if !exists {
		http.Error(w, "An account doesn't exist with that username", 404)
		return
	}

	// Compare password hash in DB with given password
	// TODO: Get rid of instance where I declare var and then assign; I don't need to do that
	// TODO: Hide hashing cost somewhere else
	passwordHash, err := models.GetPasswordHash(env.DB, l.Username)
	if err != nil {
		http.Error(w, "Error retrieving password", 500)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(l.Password))
	if err != nil {
		http.Error(w, "Incorrect password", 401)
		return
	}

	// If we reach here, the password is correct and we can issue an access/refresh token pair
	// Create access token
	// TODO: Separate these two into seperate methods since we use this in CreateUser as well
	u, err := models.GetUserByUsername(env.DB, l.Username)
	if err != nil {
		http.Error(w, "Error retrieving user", 500)
		return
	}
	accessToken, err := models.GetAccessToken(u)
	if err != nil {
		http.Error(w, "Failed to create access token", 500)
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

	// Create refresh token
	refreshToken, err := models.GetRefreshToken(env.DB, u)
	if err != nil {
		http.Error(w, "Failed to create refresh token", 500)
		return
	}

	// store refresh token in cookie for client
	refreshCookie := http.Cookie{
		Name:  "mRT",
		Value: refreshToken,
		//Domain: ".multapply.io", // TODO: When domain is hooked up, uncomment this
		Expires:  time.Now().Add(time.Minute * 525600), // TODO: Use constant for this
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &refreshCookie)
}

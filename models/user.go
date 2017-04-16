package models

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

// User - A user account in the app
type User struct {
	UserID       int    `json:"user_id"`
	Username     string `json:"username" db:"username"`
	FirstName    string `json:"first_name" db:"first_name"`
	LastName     string `json:"last_name" db:"last_name"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	Roles        string `json:"roles" db:"roles"`
}

// NewUser - A struct representing new user info received from a POST
// request to create an account
type NewUser struct {
	Username        string `json:"username"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	PasswordHash    string `json:"password_hash"`
}

// Trim - trims all string attributes of this NewUser object
// TODO: Test this function
func (n *NewUser) Trim() {
	n.Username = strings.TrimSpace(n.Username)
	n.FirstName = strings.TrimSpace(n.FirstName)
	n.LastName = strings.TrimSpace(n.LastName)
	n.Email = strings.TrimSpace(n.Email)
	n.Password = strings.TrimSpace(n.Password)
	n.ConfirmPassword = strings.TrimSpace(n.ConfirmPassword)
}

// CreateUser - Create a new User from a given NewUser
func CreateUser(n *NewUser) *User {
	u := &User{
		Username:     n.Username,
		FirstName:    n.FirstName,
		LastName:     n.LastName,
		Email:        n.Email,
		PasswordHash: n.PasswordHash,
		Roles:        "USER",
	}

	return u
}

// InsertUser - Insert User into DB
func InsertUser(db *sqlx.DB, u *User) error {
	_, err := db.NamedExec(`INSERT INTO users 
		(username, first_name, last_name, email, password_hash, roles)
		VALUES (:username, :first_name, :last_name, :email, :password_hash, :roles)`, u)
	return err
}

// RemoveNewUser - Remove a user we just inserted into the DB
// This function is for use when something goes wrong after we create
// a user in handlers/user.go:CreateUser()
func RemoveNewUser(db *sqlx.DB, username string) error {
	_, err := db.Exec(`DELETE FROM users WHERE username=$1`, username)
	return err
}

//TEMP
func GetAllUsers(db *sqlx.DB) ([]User, error) {
	// rows, err := db.Query("select * from users")
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	users := []User{}
	err := db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	return users, nil
}

// UsernameExists - returns whether an account with the given username
// exists in the database already
func UsernameExists(db *sqlx.DB, username string) (bool, error) {
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username) // TODO: Test with SELECT *
	if err != nil {
		return false, err
	}

	return exists, nil
}

// EmailExists - returns whether an account with the given email
// exists in the database already
// TODO: Possibly merge this into more general function like AccountExistsWithParam
func EmailExists(db *sqlx.DB, email string) (bool, error) {
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email) // TODO: Test with SELECT *
	if err != nil {
		return false, err
	}

	return exists, nil
}

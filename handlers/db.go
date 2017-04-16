package handlers

import (
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/multapply/multapply/util/constants"
)

// Env - Used for dependency injection
type Env struct {
	DB *sqlx.DB
}

// InitDB - initializes database
func InitDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Ping DB
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Create tables if applicable
	tx := db.MustBegin()
	tx.MustExec(createUsersTable)
	tx.MustExec(createRefreshTokensTable)
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// DB initialization queries
var createUsersTable = `CREATE TABLE IF NOT EXISTS users (
	user_id SERIAL PRIMARY KEY NOT NULL UNIQUE,
	username varchar(` + strconv.Itoa(constants.UsernameLength) + `) UNIQUE NOT NULL,
	first_name varchar(` + strconv.Itoa(constants.FirstNameLength) + `),
	last_name varchar(` + strconv.Itoa(constants.LastNameLength) + `),
	email varchar(` + strconv.Itoa(constants.EmailLength) + `) UNIQUE NOT NULL,
	password_hash text NOT NULL,
	roles varchar(100));`

var createRefreshTokensTable = `CREATE TABLE IF NOT EXISTS refresh_tokens (
	token_id SERIAL PRIMARY KEY NOT NULL UNIQUE,
	username varchar(` + strconv.Itoa(constants.UsernameLength) + `) NOT NULL REFERENCES users(username));`

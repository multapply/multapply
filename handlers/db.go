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
	tx.MustExec(createJobsTable)
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// DB initialization queries
// TODO: Change password_hash to varchar???
var createUsersTable = `CREATE TABLE IF NOT EXISTS users (
	user_id SERIAL PRIMARY KEY NOT NULL UNIQUE,
	username VARCHAR(` + strconv.Itoa(constants.UsernameLength) + `) UNIQUE NOT NULL,
	first_name VARCHAR(` + strconv.Itoa(constants.FirstNameLength) + `),
	last_name VARCHAR(` + strconv.Itoa(constants.LastNameLength) + `),
	email VARCHAR(` + strconv.Itoa(constants.EmailLength) + `) UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	roles VARCHAR(100));`

var createRefreshTokensTable = `CREATE TABLE IF NOT EXISTS refresh_tokens (
	token_id SERIAL PRIMARY KEY NOT NULL UNIQUE,
	username VARCHAR(` + strconv.Itoa(constants.UsernameLength) + `) NOT NULL REFERENCES users(username));`

// TODO: Link up company_id as foreign key
// TODO: Maybe make separate table for job descriptions, also make them varchar vs text??
// TODO: Decide on limit for icon_url/url... 150 for now (see below)
// TODO: Figure out arrays in Postgres for foreign key array for Locations
var createJobsTable = `CREATE TABLE IF NOT EXISTS jobs (
	job_id SERIAL PRIMARY KEY NOT NULL UNIQUE,
	author_id INTEGER NOT NULL REFERENCES users(user_id),
	company_id INTEGER NOT NULL,
	title VARCHAR(` + strconv.Itoa(constants.JobTitleLength) + `) NOT NULL,
	description TEXT NOT NULL,
	views INTEGER DEFAULT 0,
	icon_url VARCHAR(150) NOT NULL,
	url VARCHAR(150) NOT NULL,
	is_active BOOLEAN DEFAULT TRUE);`

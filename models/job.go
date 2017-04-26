package models

import "github.com/jmoiron/sqlx"

// Job - A Job object - referenced by JobListings, Locations
// This is the core model for storing information about a Job
// TODO: Something to consider: Upload icon after successful job POST, so make it not mandatory?
type Job struct {
	JobID       int    `json:"job_id" db:"job_id"`
	AuthorID    int    `json:"author_id" db:"author_id"`
	CompanyID   int    `json:"company_id" db:"company_id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Views       int    `json:"views" db:"views"`
	IconURL     string `json:"icon_url" db:"icon_url"`
	URL         string `json:"url" db:"url"`
	IsActive    bool   `json:"is_active" db:"is_active"`
}

// InsertJob - Inserts the given Job into the DB
func InsertJob(db *sqlx.DB, j *Job) error {
	_, err := db.NamedExec(`INSERT INTO jobs
		(author_id, company_id, title, description, views, icon_url, url, is_active)
		VALUES (:author_id, :company_id, :title, :description, :views, :icon_url, :url, :is_active);`,
		j)
	return err
}

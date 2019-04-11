package runner

import (
	"database/sql"
)

// Submission in the database
type Submission struct {
	ID       int64
	Status   string
	FileName string
}

func getNextSubmission(db *sql.DB) *Submission {
	return &Submission{}
}

// ClaimSubmission gets the first new submission and sets its status to running
func ClaimSubmission(db *sql.DB) *Submission {
	return nil
}

// UpdateSubmissionResult updates the result of a submission to match status
func UpdateSubmissionResult(db *sql.DB, id int64, status string) {

}

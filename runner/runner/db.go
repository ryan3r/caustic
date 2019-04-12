package runner

import (
	"database/sql"
)

// Submission in the database
type Submission struct {
	ID       int64
	Status   string
	FileName string
	db       *sql.DB
}

func getNextSubmission(db *sql.DB) (*Submission, error) {
	submission := &Submission{
		db: db,
	}

	err := db.QueryRow("SELECT submissionId, status, fileName FROM submission WHERE status = 'new' LIMIT 1").
		Scan(&submission.ID, &submission.Status, &submission.FileName)

	if len(submission.Status) == 0 {
		submission = nil
	}

	return submission, err
}

// ClaimSubmission gets the first new submission and sets its status to running
func ClaimSubmission(db *sql.DB) (*Submission, error) {
	submission, err := getNextSubmission(db)

	if submission != nil && err == nil {
		err = submission.UpdateStatus("running")
	}

	return submission, err
}

// UpdateStatus updates the status of a submission
func (s *Submission) UpdateStatus(status string) error {
	_, err := s.db.Exec("UPDATE submission WHERE id = ? SET status = ?", s.ID, status)
	return err
}

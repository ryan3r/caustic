package main

import (
	"database/sql"
)

// SubmissionStatus for a submission
type SubmissionStatus int

const (
	New          SubmissionStatus = 0
	Running      SubmissionStatus = 1
	CompileError SubmissionStatus = 2
	Ok           SubmissionStatus = 3
	Wrong        SubmissionStatus = 4
	TimeLimit    SubmissionStatus = 5
	Exception    SubmissionStatus = 6
)

func (s SubmissionStatus) String() string {
	switch s {
	case New:
		return "New"
	case Running:
		return "Running"
	case CompileError:
		return "CompileError"
	case Ok:
		return "Ok"
	case Wrong:
		return "Wrong"
	case TimeLimit:
		return "TimeLimitExceded"
	}
	return "Exception"
}

// Submission in the database
type Submission struct {
	ID       int64
	Status   SubmissionStatus
	FileName string
	Problem  int64
	db       *sql.DB
}

func getNextSubmission(db *sql.DB) (*Submission, error) {
	submission := &Submission{
		db: db,
	}

	err := db.QueryRow("SELECT submission_id, status, file_name, problem FROM submission WHERE status = 'new' LIMIT 1").
		Scan(&submission.ID, &submission.Status, &submission.FileName, &submission.Problem)

	if err == sql.ErrNoRows {
		err = nil
		submission = nil
	}

	return submission, err
}

// ClaimSubmission gets the first new submission and sets its status to running
func ClaimSubmission(db *sql.DB) (*Submission, error) {
	submission, err := getNextSubmission(db)

	if submission != nil && err == nil {
		err = submission.UpdateStatus(Running)
	}

	return submission, err
}

// UpdateStatus updates the status of a submission
func (s *Submission) UpdateStatus(status SubmissionStatus) error {
	s.Status = status
	_, err := s.db.Exec("UPDATE submission SET status = ? WHERE submission_id = ?", s.Status, s.ID)
	return err
}

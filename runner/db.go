package main

import (
	"database/sql"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// SubmissionStatus for a submission
type SubmissionStatus int

const (
	// New submission
	New SubmissionStatus = 0
	// Running submission
	Running SubmissionStatus = 1
	// CompileError in submission
	CompileError SubmissionStatus = 2
	// Ok aka accepted submission
	Ok SubmissionStatus = 3
	// Wrong answer
	Wrong SubmissionStatus = 4
	// TimeLimit exceded
	TimeLimit SubmissionStatus = 5
	// Exception in submission
	Exception SubmissionStatus = 6
	// RunnerError means the runner failed to run it
	RunnerError SubmissionStatus = 7
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
	case RunnerError:
		return "RunnerError"
	}
	return "Exception"
}

// Submission in the database
type Submission struct {
	ID        int64
	Status    SubmissionStatus
	FileName  string
	Problem   string
	Submitter string
	db        *sql.DB
}

func getNextSubmission(db *sql.DB) (*Submission, error) {
	submission := &Submission{
		db: db,
	}

	var ID int64

	err := db.QueryRow("SELECT submission_id, status, file_name, problem, submitter FROM submission WHERE status = 'new' LIMIT 1").
		Scan(&ID, &submission.Status, &submission.FileName, &submission.Problem, &submission.Submitter)

	submission.ID = ID

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

	// Update the rejection count
	if status == Ok {
		var rejections int64
		// TODO: Don't lower score after a correct answer
		err := s.db.QueryRow("select count(*) from submission where submitter=? and problem=? and status > 3", s.Submitter, s.Problem).
			Scan(&rejections)
		if err != nil {
			return err
		}

		_, err = s.db.Exec("update submission set rejections=? where submitter=? and problem=?", rejections, s.Submitter, s.Problem)
		if err != nil {
			return err
		}
	}

	_, err := s.db.Exec("UPDATE submission SET status = ? WHERE submission_id = ?", s.Status, s.ID)
	return err
}

// FindContainer by label
func FindContainer(cli DockerClient, label string) (*types.ContainerJSON, error) {
	// Try to find the mysql container by label
	containers, err := cli.cli.ContainerList(cli.ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: label,
		}, filters.KeyValuePair{
			Key:   "status",
			Value: "running",
		}),
	})
	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return nil, errors.New("could not find a suitible container")
	}

	if len(containers) > 1 {
		return nil, errors.New("found multiple suitible containers")
	}

	// Get the container info
	ctr, err := cli.cli.ContainerInspect(cli.ctx, containers[0].ID)
	return &ctr, err
}

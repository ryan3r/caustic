package main

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
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

// CreateConnection to the mysql container in docker compose
func CreateConnection(cli DockerClient) (*sql.DB, error) {
	info, err := FindContainer(cli, "com.ryan3r.caustic.is-db")
	if err != nil {
		return nil, err
	}

	if len(info.NetworkSettings.Ports["3306/tcp"]) == 0 {
		return nil, errors.New("Port 3306 is not exposed on mysql")
	}

	binding := info.NetworkSettings.Ports["3306/tcp"][0]

	// Find the mysql password
	password := ""
	for _, pair := range info.Config.Env {
		parts := strings.Split(pair, "=")
		if parts[0] == "MYSQL_ROOT_PASSWORD" {
			password = parts[1]
		}
	}

	if password == "" {
		return nil, errors.New("could not find mysql root password")
	}

	conn := "root:" + password + "@tcp(" + binding.HostIP + ":" + binding.HostPort + ")/caustic"
	return sql.Open("mysql", conn)
}

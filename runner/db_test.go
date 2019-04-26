package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCanLoadCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error setting up mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"submissionId", "status", "fileName", "problem", "submitter"}).
		AddRow(1, New, "foo.java", 1, "foo")

	mock.ExpectQuery("SELECT (.+) FROM submission").
		WillReturnRows(rows)

	submission, _ := getNextSubmission(db)

	if submission.ID != 1 {
		t.Errorf("Picked the wrong submission (%v)", mock.ExpectationsWereMet())
	}
}

func TestLoadsNothing(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error setting up mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"submissionId", "status", "fileName", "problem", "submitter"})
	mock.ExpectQuery("^SELECT (.+) FROM submission").WillReturnRows(rows)

	if submission, _ := getNextSubmission(db); submission != nil {
		t.Errorf("Expected nil but got %v", submission)
	}
}

func TestClaimUpdatesStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error setting up mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"submissionId", "status", "fileName", "problem", "submitter"}).
		AddRow(2, New, "foo.java", 2, "foo")

	mock.ExpectQuery("^SELECT (.+) FROM submission").WillReturnRows(rows)
	mock.ExpectExec("UPDATE submission").WithArgs(Running, 2)

	ClaimSubmission(db)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations not met: %s", err)
	}
}

func TestUpdateStatusWithResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error setting up mock: %s", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE submission").WithArgs(Exception, 1)

	submission := &Submission{
		ID: 1,
		db: db,
	}

	submission.UpdateStatus(Exception)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations not met: %s", err)
	}
}

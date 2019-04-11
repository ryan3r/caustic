package runner

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

	rows := sqlmock.NewRows([]string{"submissionId", "status", "fileName"}).
		AddRow(1, "new", "foo.java").
		AddRow(2, "new", "bar.java")

	mock.ExpectQuery("^SELECT (.+) FROM submissions$").WillReturnRows(rows)

	submission := getNextSubmission(db)

	if submission.ID != 1 {
		t.Errorf("Picked the wrong submission (%v)", submission.ID)
	}
}

func TestOnlyLoadsNew(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error setting up mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"submissionId", "status", "fileName"}).
		AddRow(1, "running", "foo.java").
		AddRow(2, "new", "bar.java")

	mock.ExpectQuery("^SELECT (.+) FROM submissions$").WillReturnRows(rows)

	if submission := getNextSubmission(db); submission.ID != 2 {
		t.Errorf("Picked the wrong submission (%v)", submission.ID)
	}
}

func TestLoadsNothing(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error setting up mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"submissionId", "status", "fileName"}).
		AddRow(1, "running", "foo.java").
		AddRow(2, "running", "bar.java")

	mock.ExpectQuery("^SELECT (.+) FROM submissions$").WillReturnRows(rows)

	if submission := getNextSubmission(db); submission != nil {
		t.Errorf("Expected nil but got %v", submission)
	}
}

func TestClaimUpdatesStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error setting up mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"submissionId", "status", "fileName"}).
		AddRow(1, "new", "foo.java").
		AddRow(2, "running", "bar.java")

	mock.ExpectQuery("^SELECT (.+) FROM submissions$").WillReturnRows(rows)
	mock.ExpectExec("UPDATE submission WHERE id = 1")

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

	mock.ExpectExec("UPDATE submission WHERE id = 1")

	UpdateSubmissionResult(db, 1, "exception")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations not met: %s", err)
	}
}

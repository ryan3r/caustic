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
		AddRow(1, "new", "foo.java")

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

	rows := sqlmock.NewRows([]string{"submissionId", "status", "fileName"})
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

	rows := sqlmock.NewRows([]string{"submissionId", "status", "fileName"}).
		AddRow(2, "new", "foo.java")

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

	mock.ExpectExec("UPDATE submission").WithArgs(1, Exception)

	submission := &Submission{
		ID: 1,
		db: db,
	}

	submission.UpdateStatus(Exception)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations not met: %s", err)
	}
}

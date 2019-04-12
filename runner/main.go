package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"runner"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	flag.Parse()
	os.Chdir("/mnt/submissions")

	db, err := sql.Open("mysql", "root:"+os.Getenv("MYSQL_ROOT_PASSWORD")+"@tcp(db)/caustic")
	if err != nil {
		panic(err)
	}

	for {
		submission, err := runner.ClaimSubmission(db)
		if err != nil {
			fmt.Println("Failed to claim a submission:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if submission == nil {
			time.Sleep(5 * time.Second)
			continue
		}

		err = submission.UpdateStatus(runner.Test(submission.FileName, "1 2 3"))
		if err != nil {
			fmt.Println("Error status for %s: %s", submission.FileName, err)
		}
	}
}

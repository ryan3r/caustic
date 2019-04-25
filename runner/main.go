package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/client"
	_ "github.com/go-sql-driver/mysql"
)

const (
	// MaxPings to the db before crashing (is high because the docker db must initialize)
	MaxPings = 15
	// RecoveryTime to wait between pings or errors returned by the db
	RecoveryTime = 5 * time.Second
	// PollingInterval is the cool down interval for when we run out of submissions to claim
	PollingInterval = 1 * time.Second
	// Version is our current version
	Version = "1.0"
	// RunnerCount is the number runner workers in use
	RunnerCount = 5
)

// Load the language config and pull the language images
func loadLanguages(cli *DockerClient) {
	langFile, err := ioutil.ReadFile("languages.json")
	if err != nil {
		fmt.Printf("Failed to open laugage definitions: %s\n", err)
		os.Exit(127)
	}

	languageDefs = make(map[string]*LanguageDef)
	if err := json.Unmarshal(langFile, &languageDefs); err != nil {
		fmt.Println("Failed to parse languages")
		os.Exit(127)
	}

	fmt.Println("Pulling language images")
	if err := cli.PullAll(); err != nil {
		fmt.Println("An error occured while pulling language images")
		os.Exit(127)
	}
}

// Connect to the database
func connectDb() *sql.DB {
	dbURL, ok := os.LookupEnv("MYSQL_URL")
	if !ok {
		dbURL = "root:password@tcp(localhost:3307)/caustic"
	}

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		fmt.Printf("Failed to connect to local database: %s\n", err)
		os.Exit(127)
	}

	fmt.Printf("Attempting to connect to mysql (timeout in %v)\n", MaxPings*RecoveryTime)
	// Ping the db upto maxPings times before failing
	for pings := 0; db.Ping() != nil; pings++ {
		if pings > MaxPings {
			fmt.Printf("Failed to connect to db after %v attempts\n", MaxPings)
			os.Exit(127)
		}

		time.Sleep(RecoveryTime)
	}

	return db
}

// Update a submision's status or print an error
func updateStatus(submission *Submission, status SubmissionStatus) {
	err := submission.UpdateStatus(status)
	if err != nil {
		fmt.Printf("Error updating status for %v: %s\n", submission.ID, err)
	}
}

// A goroutine for claiming submissions
func claimSubmissions(db *sql.DB, submissionC chan *Submission) {
	for {
		submission, err := ClaimSubmission(db)
		if err != nil {
			fmt.Println("Failed to claim a submission:", err)
			time.Sleep(RecoveryTime)
			continue
		}

		// Sleep if there are no submissions
		if submission == nil {
			time.Sleep(PollingInterval)
			continue
		}

		submissionC <- submission
	}
}

// A goroutine for running submissions
func runSubmissions(cli *DockerClient, submissions, problems string, submissionC chan *Submission) {
	for {
		submission := <-submissionC
		fmt.Println("Running submission", submission.ID)

		strID := fmt.Sprintf("%v", submission.ID)
		status, err := Test(cli, filepath.Join(submissions, strID), submission.FileName, filepath.Join(problems, submission.Problem))
		if err != nil {
			fmt.Printf("Error testing submission: %v (%s)\n", submission.ID, err)

			// Put the submission back for another runner
			updateStatus(submission, New)
		}

		fmt.Printf("Submission status %v: %s\n", submission.ID, status)
		updateStatus(submission, status)
	}
}

func main() {
	fmt.Printf("Caustic runner v%s (--COMMIT-HASH-HERE--)\n", Version)

	apiClient, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		fmt.Printf("Failed to connect to docker: %s\n", err)
		fmt.Println("Make sure you have docker installed and it is running")
		os.Exit(127)
	}

	cli := &DockerClient{
		ctx: context.Background(),
		cli: apiClient,
	}

	loadLanguages(cli)
	db := connectDb()

	submissions, err := filepath.Abs("submissions")
	if err != nil {
		fmt.Printf("Error getting submissions path: %s\n", err)
		os.Exit(127)
	}
	problems, err := filepath.Abs("problems")
	if err != nil {
		fmt.Printf("Error getting problems path: %s\n", err)
		os.Exit(127)
	}

	// Start testing submissions
	fmt.Println("Connected\nWaiting for submissions")

	submissionC := make(chan *Submission)

	for i := 0; i < RunnerCount; i++ {
		go runSubmissions(cli, submissions, problems, submissionC)
	}

	claimSubmissions(db, submissionC)
}

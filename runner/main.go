package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
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
	PollingInterval = 3 * time.Second
)

func main() {
	apiClient, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		fmt.Printf("Failed to connect to docker: %s\n", err)
		fmt.Println("Make sure you have docker installed and it is running")
		os.Exit(127)
	}

	cli := DockerClient{
		ctx: context.Background(),
		cli: apiClient,
	}

	RegisterLanguage("java", &LanguageDef{
		Image:          "openjdk:13-jdk-alpine",
		CompileCommand: []string{"javac", "%f"},
		RunCommand:     []string{"java", "%n"},
		Artifacts:      []string{"%n.class"},
	})

	RegisterLanguage("cpp", &LanguageDef{
		Image:          "gcc:5",
		CompileCommand: []string{"g++", "%f", "-o", "%n"},
		RunCommand:     []string{"./%n"},
		Artifacts:      []string{"%n"},
	})

	RegisterLanguage("py", &LanguageDef{
		Image:      "python",
		RunCommand: []string{"python", "%f"},
	})

	fmt.Println("Pulling language images")
	if err := cli.PullAll(); err != nil {
		fmt.Println("An error occured while pulling language images")
		os.Exit(127)
	}

	db, err := CreateConnection(cli)
	if err != nil {
		fmt.Printf("Failed to find mysql in docker trying localhost (%s)\n", err)

		// This is to make spring boot development easier don't use this outside of development
		db, err = sql.Open("mysql", "root:Pa$$word1@tcp(localhost:3306)/caustic")
		if err != nil {
			fmt.Printf("Failed to connect to local database: %s\n", err)
			os.Exit(127)
		}
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

	// Start testing submissions
	fmt.Println("Connected\nWaiting for submissions")
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

		fmt.Println("Running submission", submission.ID)

		status, err := Test(cli, "test-files/0", submission.FileName, "test-files/problem")
		if err != nil {
			fmt.Printf("Error testing submission: %v\n", submission.ID)

			// Put the submission back for another runner
			err = submission.UpdateStatus(New)
			if err != nil {
				fmt.Printf("Error updating status for %v: %s\n", submission.ID, err)
			}
		}

		fmt.Printf("Submission status %v: %s\n", submission.ID, status)

		err = submission.UpdateStatus(status)
		if err != nil {
			fmt.Printf("Error updating status for %v: %s\n", submission.ID, err)
		}
	}
}

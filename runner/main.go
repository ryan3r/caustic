package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/docker/docker/client"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	apiClient, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	panicIf(err)

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

	panicIf(cli.PullAll())

	db, err := CreateConnection(cli)
	panicIf(err)

	for {
		submission, err := ClaimSubmission(db)
		if err != nil {
			fmt.Println("Failed to claim a submission:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if submission == nil {
			fmt.Println("Empty");
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println("Running");

		status, err := Test(cli, "test-files/0", submission.FileName, "test-files/problem")
		panicIf(err)

		fmt.Printf("Status %s: %s\n", submission.FileName, status)

		err = submission.UpdateStatus(status)
		if err != nil {
			fmt.Printf("Error updating status for %s: %s\n", submission.FileName, err)
		}
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

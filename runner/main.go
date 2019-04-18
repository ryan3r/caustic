package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/client"
)

var (
	ErrExitStatusError = errors.New("Program exited with non-zero exit status")
	ErrTimeLimit       = errors.New("Program took too long to run")
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: main <submission dir> <submission file>")
		os.Exit(1)
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	panicIf(err)

	panicIf(pullImage(ctx, cli, "openjdk:13-jdk-alpine"))

	msg, err := runAsContainer(ctx, cli, "openjdk:13-jdk-alpine", os.Args[1], []string{"javac", os.Args[2]})
	if err != nil {
		if err == ErrExitStatusError {
			fmt.Println("Result: compile error")
			return
		}
	}
	fmt.Println(msg)

	res, err := createContainer(ctx, cli, "openjdk:13-jdk-alpine", os.Args[1], false, []string{"sleep", "100"})
	panicIf(err)

	className, _ := detectType(os.Args[2])

	tests, err := ioutil.ReadDir(os.Args[3])
	panicIf(err)

	for _, file := range tests {
		_, fileType := detectType(file.Name())

		if fileType != "in" {
			continue
		}

		msg, err := execInContainer(ctx, cli, res.ID, 3*time.Second, filepath.Join(os.Args[3], file.Name()), []string{"java", className})
		if err != nil {
			if err == ErrExitStatusError {
				fmt.Println("Result: exception")
			} else if err == ErrTimeLimit {
				fmt.Println("Result: time limit exceded")
				res, err = createContainer(ctx, cli, "openjdk:13-jdk-alpine", os.Args[1], false, []string{"sleep", "100"})
				panicIf(err)
			} else {
				panic(err)
			}
		}

		fmt.Print(msg)
	}

	cli.ContainerStop(ctx, res.ID, nil)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

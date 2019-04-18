package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/client"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: main <submission dir> <submission file>")
		os.Exit(1)
	}

	apiClient, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	panicIf(err)

	cli := DockerClient{
		ctx: context.Background(),
		cli: apiClient,
	}

	panicIf(cli.Pull("openjdk:13-jdk-alpine"))

	compileCtr := Container{
		Docker:     cli,
		Image:      "openjdk:13-jdk-alpine",
		Cmd:        []string{"javac", os.Args[2]},
		WorkingDir: "/mnt",
		Out:        os.Stdout,
	}

	panicIf(compileCtr.BindDir(os.Args[1], "/mnt", false))
	panicIf(compileCtr.Run())

	if err := compileCtr.Wait(); err != nil {
		if err == ErrExitStatusError {
			fmt.Println("Result: compile error")
			return
		} else {
			panic(err)
		}
	}

	testCtr := Container{
		Docker:     cli,
		Image:      "openjdk:13-jdk-alpine",
		Cmd:        []string{"sleep", "100"},
		WorkingDir: "/mnt",
		Out:        os.Stdout,
		ReadOnly:   true,
	}

	panicIf(testCtr.BindDir(os.Args[1], "/mnt", true))
	panicIf(testCtr.Run())

	className, _ := detectType(os.Args[2])

	tests, err := ioutil.ReadDir(os.Args[3])
	panicIf(err)

	for _, file := range tests {
		_, fileType := detectType(file.Name())

		if fileType != "in" {
			continue
		}

		file, err := os.Open(filepath.Join(os.Args[3], file.Name()))
		panicIf(err)

		exec := ContainerExec{
			Container: &testCtr,
			Cmd:       []string{"java", className},
			In:        file,
			Out:       os.Stdout,
		}

		exec.Run()
		exec.StartKillTimer(3 * time.Second)

		if err := <-exec.ExitC; err != nil {
			if err == ErrExitStatusError {
				fmt.Println("Result: exception")
			} else if err == ErrTimeLimit {
				fmt.Println("Result: time limit exceded")
				testCtr.Run()
				panicIf(err)
			} else {
				panic(err)
			}
		}
	}

	panicIf(testCtr.Stop())
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

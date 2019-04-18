package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ahmetalpbalkan/dlog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

var (
	ErrExitStatusError = errors.New("Program exited with non-zero exit status")
	ErrTimeLimit       = errors.New("Program took too long to run")
)

type PullProgress struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

type PullResponse struct {
	Status         string       `json:"status"`
	ProgressDetail PullProgress `json:"progressDetail"`
	Progress       string       `json:"progress"`
	ID             string       `json:"id"`
}

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

// Detect the filetype and name of file
func detectType(filename string) (string, string) {
	idx := strings.LastIndex(filename, ".")
	return filename[:idx], filename[idx+1:]
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func createContainer(ctx context.Context, cli client.APIClient, image string, workDir string, readonly bool, cmd []string) (container.ContainerCreateCreatedBody, error) {
	workDir, err := filepath.Abs(workDir)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}

	res, err := cli.ContainerCreate(ctx, &container.Config{
		AttachStdin:     true,
		AttachStdout:    true,
		AttachStderr:    true,
		Image:           image,
		WorkingDir:      "/mnt",
		NetworkDisabled: true,
		Cmd:             cmd,
	}, &container.HostConfig{
		ReadonlyRootfs: true,
		AutoRemove:     true,
		Mounts: []mount.Mount{
			mount.Mount{
				Type:     mount.TypeBind,
				Source:   workDir,
				Target:   "/mnt",
				ReadOnly: readonly,
			},
		},
	}, nil, "")
	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}

	return res, cli.ContainerStart(ctx, res.ID, types.ContainerStartOptions{})
}

func execInContainer(ctx context.Context, cli client.APIClient, containerID string, maxTime time.Duration, inFileName string, cmd []string) (string, error) {
	execID, err := cli.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		Cmd:          cmd,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return "", err
	}

	isTLE := false
	done := make(chan bool, 1)
	go func() {
		timer := time.NewTimer(maxTime)
		select {
		case <-done:
		case <-timer.C:
			isTLE = true
			panicIf(cli.ContainerKill(ctx, containerID, "SIGKILL"))
		}
	}()

	con, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}

	inFile, err := os.Open(inFileName)
	if err != nil {
		return "", err
	}

	io.Copy(con.Conn, inFile)
	con.CloseWrite()

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(dlog.NewReader(con.Reader))

	info, err := cli.ContainerExecInspect(ctx, execID.ID)
	done <- true
	if err != nil {
		return "", err
	}

	if isTLE {
		return "", ErrTimeLimit
	}

	if info.ExitCode != 0 {
		return buffer.String(), ErrExitStatusError
	}

	return buffer.String(), nil
}

func pullImage(ctx context.Context, cli client.APIClient, image string) error {
	pullStats, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	needsNewline := false
	decoder := json.NewDecoder(pullStats)
	for decoder.More() {
		var progress PullResponse
		decoder.Decode(&progress)

		switch progress.Status {
		case "Downloading":
			needsNewline = true
			fmt.Printf("\rDownloading %s", progress.Progress)

		case "Extracting":
			needsNewline = true
			fmt.Printf("\rExtracting %s", progress.Progress)

		default:
			if needsNewline {
				fmt.Println()
				needsNewline = false
			}
			fmt.Println(progress.Status)
		}
	}

	return nil
}

func runAsContainer(ctx context.Context, cli client.APIClient, image string, dirName string, cmd []string) (string, error) {
	res, err := createContainer(ctx, cli, image, dirName, false, cmd)
	if err != nil {
		return "", err
	}

	reader, err := cli.ContainerLogs(ctx, res.ID, types.ContainerLogsOptions{
		ShowStdout: true,
	})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	buffer := new(bytes.Buffer)
	go func() {
		buffer.ReadFrom(dlog.NewReader(reader))
	}()

	doneC, errC := cli.ContainerWait(ctx, res.ID, "")
	select {
	case err = <-errC:
		return "", err
	case info := <-doneC:
		if info.StatusCode != 0 {
			return buffer.String(), ErrExitStatusError
		}
	}

	return buffer.String(), nil
}

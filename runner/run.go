package main

import (
	"bytes"
	"context"
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
	languageDefs       map[string]*LanguageDef
	ErrUnknownFileType = errors.New("Unknown filetype")
	ErrExitStatusError = errors.New("Program exited with non-zero exit status")
	ErrTimeLimit       = errors.New("Program took too long to run")
)

// LanguageDef defines how to handle a file type
type LanguageDef struct {
	Image          string   `json:"image"`
	CompileCommand []string `json:"compile"`
	RunCommand     []string `json:"run"`
	Artifacts      []string `json:"artifacts"`
}

// Detect the filetype and name of file
func detectType(fileName string) (string, string) {
	idx := strings.LastIndex(fileName, ".")
	return fileName[:idx], fileName[idx+1:]
}

func expandTemplate(template []string, fileName, name string) []string {
	templated := append([]string(nil), template...)

	for i, orig := range templated {
		tmp := strings.ReplaceAll(orig, "%n", name)
		templated[i] = strings.ReplaceAll(tmp, "%f", fileName)
	}

	return templated
}

// Compile a file
func Compile(cli DockerClient, problemDir, fileName string) error {
	name, ft := detectType(fileName)
	def := languageDefs[ft]

	if def == nil {
		return ErrUnknownFileType
	}

	// This language is interpreted
	if len(def.CompileCommand) == 0 {
		return nil
	}

	compileCommand := expandTemplate(def.CompileCommand, fileName, name)

	compileCtr := Container{
		Docker:     cli,
		Image:      def.Image,
		Cmd:        compileCommand,
		WorkingDir: "/mnt",
		Out:        os.Stdout,
	}

	if err := compileCtr.BindDir(problemDir, "/mnt", false); err != nil {
		return err
	}

	if err := compileCtr.Run(); err != nil {
		return err
	}

	return compileCtr.Wait()
}

// Runner is a runner for a single test
type Runner struct {
	problemDir string
	fileName   string
	container  *Container
	timeLimit  time.Duration
}

// NewRunner creates the container for running tests
func NewRunner(cli DockerClient, problemDir, fileName string, timeLimit time.Duration) (*Runner, error) {
	_, ft := detectType(fileName)
	def := languageDefs[ft]

	testCtr := &Container{
		Docker:     cli,
		Image:      def.Image,
		Cmd:        []string{"sleep", "100"},
		WorkingDir: "/mnt",
		Out:        os.Stdout,
		ReadOnly:   true,
	}

	if err := testCtr.BindDir(problemDir, "/mnt", true); err != nil {
		return nil, err
	}

	if err := testCtr.Run(); err != nil {
		return nil, err
	}

	return &Runner{
		problemDir: problemDir,
		fileName:   fileName,
		container:  testCtr,
		timeLimit:  timeLimit,
	}, nil
}

// Run the submission with the test case
func (r *Runner) Run(in io.Reader, out io.Writer) error {
	name, ft := detectType(r.fileName)
	def := languageDefs[ft]

	exec := ContainerExec{
		Container: r.container,
		Cmd:       expandTemplate(def.RunCommand, r.fileName, name),
		In:        in,
		Out:       out,
	}

	exec.Run()
	exec.StartKillTimer(r.timeLimit)

	err := <-exec.ExitC
	if err == ErrTimeLimit {
		r.container.Run()
	}
	return err
}

// Close stops the runner container
func (r *Runner) Close() error {
	return r.container.Stop()
}

// Test will compile, run and check a program
func Test(cli DockerClient, problemDir, fileName, solutionDir string) (SubmissionStatus, error) {
	defer cleanUpArtifacts(problemDir, fileName)

	// Compile the solution
	err := Compile(cli, problemDir, fileName)
	if err == ErrExitStatusError {
		return CompileError, nil
	} else if err != nil {
		return New, err
	}

	// Run the submissions
	tests, err := ioutil.ReadDir(solutionDir)
	if err != nil {
		return New, err
	}

	runner, err := NewRunner(cli, problemDir, fileName, 3*time.Second)
	if err != nil {
		return New, err
	}
	defer runner.Close()

	for _, file := range tests {
		name, fileType := detectType(file.Name())

		if fileType != "in" {
			continue
		}

		fileIn, err := os.Open(filepath.Join(solutionDir, file.Name()))
		if err != nil {
			return New, err
		}
		defer fileIn.Close()

		outBuffer := bytes.NewBufferString("")

		if err := runner.Run(fileIn, outBuffer); err != nil {
			if err == ErrExitStatusError {
				return Exception, nil
			} else if err == ErrTimeLimit {
				return TimeLimit, nil
			} else {
				return New, err
			}
		}

		// Verify the output of the submission
		outFile, err := ioutil.ReadFile(filepath.Join(solutionDir, name+".out"))
		if err != nil {
			return Ok, err
		}

		expectedOut := strings.Trim(string(outFile), "\r\n\t ")
		solutionOut := strings.Trim(string(outBuffer.String()), "\r\n\t ")

		if expectedOut != solutionOut {
			return Wrong, nil
		}
	}

	return Ok, nil
}

func cleanUpArtifacts(problemDir, fileName string) {
	name, ft := detectType(fileName)
	def := languageDefs[ft]

	if def.Artifacts == nil {
		return
	}

	for _, artifact := range expandTemplate(def.Artifacts, fileName, name) {
		if err := os.RemoveAll(filepath.Join(problemDir, artifact)); err != nil {
			fmt.Printf("Warn: Failed to clean up artifacts (%v, %v)\n", problemDir, artifact)
		}
	}
}

// DockerClient is a APIClient + Context
type DockerClient struct {
	ctx context.Context
	cli client.APIClient
}

// Container is an abstraction of a docker container
type Container struct {
	Docker         DockerClient
	ID             string
	Image          string
	Cmd            []string
	mounts         []mount.Mount
	ReadOnly       bool
	WorkingDir     string
	NetworkEnabled bool
	Out            io.Writer
}

// BindDir adds a bind mount to the container
func (c *Container) BindDir(src string, dest string, readonly bool) (err error) {
	src, err = filepath.Abs(src)

	c.mounts = append(c.mounts, mount.Mount{
		Type:     mount.TypeBind,
		Source:   src,
		Target:   dest,
		ReadOnly: readonly,
	})

	return
}

// Run creates and starts the container and connects the stdio
func (c *Container) Run() error {
	// Create the container
	res, err := c.Docker.cli.ContainerCreate(c.Docker.ctx, &container.Config{
		AttachStdout:    true,
		AttachStderr:    true,
		Image:           c.Image,
		WorkingDir:      c.WorkingDir,
		NetworkDisabled: !c.NetworkEnabled,
		Cmd:             c.Cmd,
	}, &container.HostConfig{
		ReadonlyRootfs: c.ReadOnly,
		AutoRemove:     true,
		Mounts:         c.mounts,
	}, nil, "")
	if err != nil {
		return err
	}

	c.ID = res.ID

	// Start the container
	err = c.Docker.cli.ContainerStart(c.Docker.ctx, c.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	// Wire up stdout/stderr
	reader, err := c.Docker.cli.ContainerLogs(c.Docker.ctx, res.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return err
	}

	go func() {
		io.Copy(c.Out, dlog.NewReader(reader))
		reader.Close()
	}()

	return nil
}

// Wait for the container to exit
func (c *Container) Wait() error {
	doneC, errC := c.Docker.cli.ContainerWait(c.Docker.ctx, c.ID, "")
	select {
	case err := <-errC:
		return err
	case info := <-doneC:
		if info.StatusCode != 0 {
			return ErrExitStatusError
		}
	}

	return nil
}

// Kill a container
func (c *Container) Kill() error {
	return c.Docker.cli.ContainerKill(c.Docker.ctx, c.ID, "SIGKILL")
}

// Stop a container
func (c *Container) Stop() error {
	return c.Docker.cli.ContainerStop(c.Docker.ctx, c.ID, nil)
}

// ContainerExec is a program executed inside a container other than the main program
type ContainerExec struct {
	Container   *Container
	ID          string
	Cmd         []string
	In          io.Reader
	Out         io.Writer
	ExitC       chan error
	isTimerKill bool
}

// StartKillTimer starts a timeout
func (c *ContainerExec) StartKillTimer(timeout time.Duration) {
	go func() {
		timer := time.NewTimer(timeout)
		select {
		case <-c.ExitC:
		case <-timer.C:
			c.isTimerKill = true
			c.Container.Kill()
		}
	}()
}

// Run creates the exec, starts it and wires up the stdio
func (c *ContainerExec) Run() error {
	// Create exec process
	docker := c.Container.Docker
	execID, err := docker.cli.ContainerExecCreate(docker.ctx, c.Container.ID, types.ExecConfig{
		Cmd:          c.Cmd,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return err
	}

	c.ID = execID.ID

	c.ExitC = make(chan error, 1)
	// Connect stdio
	conn, err := docker.cli.ContainerExecAttach(docker.ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	go func() {
		io.Copy(conn.Conn, c.In)
		conn.CloseWrite()
	}()

	go func() {
		io.Copy(c.Out, dlog.NewReader(conn.Reader))
		conn.Close()

		// Get the exit status (no wait option)
		info, err := docker.cli.ContainerExecInspect(docker.ctx, execID.ID)
		if err != nil {
			c.ExitC <- err
			return
		}

		if c.isTimerKill {
			c.ExitC <- ErrTimeLimit
			return
		}

		if info.ExitCode != 0 {
			c.ExitC <- ErrExitStatusError
			return
		}

		c.ExitC <- nil
	}()

	return nil
}

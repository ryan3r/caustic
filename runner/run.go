package main

import (
	"context"
	"errors"
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
	Image          string
	CompileCommand []string
	RunCommand     []string
}

// RegisterLanguage registers a language for the runner
func RegisterLanguage(ft string, def *LanguageDef) {
	if languageDefs == nil {
		languageDefs = make(map[string]*LanguageDef)
	}
	languageDefs[ft] = def
}

// Detect the filetype and name of file
func detectType(filename string) (string, string) {
	idx := strings.LastIndex(filename, ".")
	return filename[:idx], filename[idx+1:]
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
func Compile(cli DockerClient, problemDir, filename string) error {
	name, ft := detectType(filename)
	def := languageDefs[ft]

	if def == nil {
		return ErrUnknownFileType
	}

	// This language is interpreted
	if len(def.CompileCommand) == 0 {
		return nil
	}

	compileCommand := expandTemplate(def.CompileCommand, filename, name)

	compileCtr := Container{
		Docker:     cli,
		Image:      def.Image,
		Cmd:        compileCommand,
		WorkingDir: "/mnt",
		Out:        os.Stdout,
	}

	panicIf(compileCtr.BindDir(problemDir, "/mnt", false))
	panicIf(compileCtr.Run())

	return compileCtr.Wait()
}

// Runner is a runner for a single test
type Runner struct {
	problemDir string
	filename   string
	container  *Container
	timeLimit  time.Duration
}

// NewRunner creates the container for running tests
func NewRunner(cli DockerClient, problemDir, filename string, timeLimit time.Duration) (*Runner, error) {
	_, ft := detectType(filename)
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
		filename:   filename,
		container:  testCtr,
		timeLimit:  timeLimit,
	}, nil
}

// Run the submission with the test case
func (r *Runner) Run(in io.Reader, out io.Writer) error {
	name, ft := detectType(r.filename)
	def := languageDefs[ft]

	exec := ContainerExec{
		Container: r.container,
		Cmd:       expandTemplate(def.RunCommand, r.filename, name),
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
	// Compile the solution
	err := Compile(cli, os.Args[1], os.Args[2])
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
		_, fileType := detectType(file.Name())

		if fileType != "in" {
			continue
		}

		file, err := os.Open(filepath.Join(solutionDir, file.Name()))
		if err != nil {
			return New, err
		}

		if err := runner.Run(file, os.Stdout); err != nil {
			if err == ErrExitStatusError {
				return Exception, nil
			} else if err == ErrTimeLimit {
				return TimeLimit, nil
			} else {
				return New, err
			}
		}
	}

	return Ok, nil
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

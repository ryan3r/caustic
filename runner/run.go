package main

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/ahmetalpbalkan/dlog"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

// Detect the filetype and name of file
func detectType(filename string) (string, string) {
	idx := strings.LastIndex(filename, ".")
	return filename[:idx], filename[idx+1:]
}

// Compile a file
// func compile(filename string) error {
// 	compiler := "javac"
// 	_, ft := detectType(filename)

// 	if ft == "py" {
// 		return nil
// 	}

// 	if ft == "cpp" {
// 		compiler = "g++"
// 	}

// 	log.Printf("Compiling %v as %v\n", filename, ft)
// 	return exec.Command(compiler, filename).Run()
// }

// // Run a program
// func run(ctx context.Context, filename string, output chan string, errs chan error) {
// 	name, ft := detectType(filename)
// 	var cmd *exec.Cmd

// 	switch ft {
// 	case "java":
// 		cmd = exec.CommandContext(ctx, "java", name)
// 	case "cpp":
// 		cmd = exec.CommandContext(ctx, "./a.out")
// 	case "py":
// 		cmd = exec.CommandContext(ctx, "python", filename)
// 	default:
// 		log.Printf("Error unknown filetype %v\n", ft)
// 		errs <- errors.New("Unknown filetype")
// 		return
// 	}

// 	log.Printf("Running %v as %v", filename, ft)
// 	out, err := cmd.CombinedOutput()

// 	if err != nil {
// 		log.Printf("Error running %v: %v\n", filename, err.Error())
// 		errs <- err
// 	} else {
// 		log.Printf("Completed %v no errors\n", filename)
// 		output <- string(out)
// 	}
// }

// Test will compile, run and check a program
// func Test(filename string, expected string) SubmissionStatus {
// 	if err := compile(filename); err != nil {
// 		return CompileError
// 	}

// 	errors := make(chan error, 1)
// 	output := make(chan string, 1)

// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	run(ctx, filename, output, errors)
// 	cancel()

// 	select {
// 	case out := <-output: // process exited on time w/o errors
// 		if strings.Trim(out, "\r\n\t ") == expected {
// 			return Ok
// 		}
// 		return Wrong
// 	case err := <-errors: // process crashed or was killed
// 		if err.Error() == "signal: killed" {
// 			return TimeLimit
// 		}
// 		return Exception
// 	}
// }

var (
	ErrExitStatusError = errors.New("Program exited with non-zero exit status")
	ErrTimeLimit       = errors.New("Program took too long to run")
)

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
		AttachStdin:     true,
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

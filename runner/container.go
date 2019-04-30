package main

import (
	"context"
	"io"
	"path/filepath"
	"time"

	"github.com/ahmetalpbalkan/dlog"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

// DockerClient is a APIClient + Context
type DockerClient struct {
	ctx context.Context
	cli client.APIClient
}

// Container is an abstraction of a docker container
type Container struct {
	Docker         *DockerClient
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
func (c *ContainerExec) StartKillTimer(timeout time.Duration) chan bool {
	exitC := make(chan bool, 1)
	go func() {
		timer := time.NewTimer(timeout)
		select {
		case <-exitC:
			timer.Stop()
		case <-timer.C:
			c.ExitC <- ErrTimeLimit
			c.Container.Kill()
		}
	}()
	return exitC
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

		if info.ExitCode != 0 {
			c.ExitC <- ErrExitStatusError
			return
		}

		c.ExitC <- nil
	}()

	return nil
}

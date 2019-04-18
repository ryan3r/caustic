package main

import (
	"context"
	"io"
	"path/filepath"
	"strings"

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
	Out            *io.Writer
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
	if c.Out != nil {
		reader, err := c.Docker.cli.ContainerLogs(c.Docker.ctx, res.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
		if err != nil {
			return err
		}

		go func() {
			io.Copy(*c.Out, dlog.NewReader(reader))
			reader.Close()
		}()
	}

	return nil
}

type ContainerExec struct {
	Container *Container
	ID        string
	Cmd       []string
}

func (c *ContainerExec) Run() error {
	docker := c.Container.Docker
	execID, err := docker.cli.ContainerExecCreate(docker.ctx, c.ID, types.ExecConfig{
		Cmd:          c.Cmd,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return err
	}

	c.ID = execID.ID

	// isTLE := false
	// done := make(chan bool, 1)
	// go func() {
	// 	timer := time.NewTimer(maxTime)
	// 	select {
	// 	case <-done:
	// 	case <-timer.C:
	// 		isTLE = true
	// 		panicIf(cli.ContainerKill(ctx, containerID, "SIGKILL"))
	// 	}
	// }()

	// con, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	// if err != nil {
	// 	return "", err
	// }

	// inFile, err := os.Open(inFileName)
	// if err != nil {
	// 	return "", err
	// }

	// io.Copy(con.Conn, inFile)
	// con.CloseWrite()

	// buffer := new(bytes.Buffer)
	// buffer.ReadFrom(dlog.NewReader(con.Reader))

	// info, err := cli.ContainerExecInspect(ctx, execID.ID)
	// done <- true
	// if err != nil {
	// 	return "", err
	// }

	// if isTLE {
	// 	return "", ErrTimeLimit
	// }

	// if info.ExitCode != 0 {
	// 	return buffer.String(), ErrExitStatusError
	// }

	// return buffer.String(), nil
	return nil
}

// func runAsContainer(ctx context.Context, cli client.APIClient, image string, dirName string, cmd []string) (string, error) {
// 	res, err := createContainer(ctx, cli, image, dirName, false, cmd)
// 	if err != nil {
// 		return "", err
// 	}

// 	doneC, errC := cli.ContainerWait(ctx, res.ID, "")
// 	select {
// 	case err = <-errC:
// 		return "", err
// 	case info := <-doneC:
// 		if info.StatusCode != 0 {
// 			return buffer.String(), ErrExitStatusError
// 		}
// 	}

// 	return buffer.String(), nil
// }

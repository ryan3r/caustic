package main

import (
	"context"
	"io"
	"os"

	"github.com/ahmetalpbalkan/dlog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	panicIf(err)

	res, err := cli.ContainerCreate(context.Background(), &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		Image:        "openjdk:13-jdk-alpine",
		WorkingDir:   "/mnt",
		Cmd:          []string{"sleep", "100"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeVolume,
				Source: "cas",
				Target: "/mnt",
			},
		},
	}, nil, "caustic-running")
	panicIf(err)

	panicIf(cli.ContainerStart(context.Background(), res.ID, types.ContainerStartOptions{}))

	execID, err := cli.ContainerExecCreate(context.Background(), res.ID, types.ExecConfig{
		Cmd:          []string{"javac", "ok.java"},
		AttachStdout: true,
		AttachStderr: true,
	})
	panicIf(err)

	con, err := cli.ContainerExecAttach(context.Background(), execID.ID, types.ExecStartCheck{})
	panicIf(err)
	defer con.Close()

	io.Copy(os.Stdout, dlog.NewReader(con.Reader))

	for i := 0; i < 3; i++ {
		execID, err = cli.ContainerExecCreate(context.Background(), res.ID, types.ExecConfig{
			Cmd:          []string{"java", "ok"},
			AttachStdout: true,
		})
		panicIf(err)

		con, err = cli.ContainerExecAttach(context.Background(), execID.ID, types.ExecStartCheck{})
		panicIf(err)
		defer con.Close()

		io.Copy(os.Stdout, dlog.NewReader(con.Reader))
	}

	cli.ContainerStop(context.Background(), res.ID, nil)
	cli.ContainerRemove(context.Background(), res.ID, types.ContainerRemoveOptions{})
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

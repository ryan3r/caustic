package main

import (
	"context"
	"fmt"
	"os"

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

	RegisterLanguage("java", &LanguageDef{
		Image:          "openjdk:13-jdk-alpine",
		CompileCommand: []string{"javac", "%f"},
		RunCommand:     []string{"java", "%n"},
	})

	RegisterLanguage("cpp", &LanguageDef{
		Image:          "gcc:5",
		CompileCommand: []string{"g++", "%f", "-o", "%n"},
		RunCommand:     []string{"./%n"},
	})

	RegisterLanguage("py", &LanguageDef{
		Image:      "python",
		RunCommand: []string{"python", "%f"},
	})

	panicIf(cli.PullAll())

	status, err := Test(cli, os.Args[1], os.Args[2], os.Args[3])
	panicIf(err)

	fmt.Println("Status:", status)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

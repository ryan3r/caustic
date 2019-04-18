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
		Image: "openjdk:13-jdk-alpine",
		CompileCommand: func(name string) []string {
			return []string{"javac", name + ".java"}
		},
		RunCommand: func(name string) []string {
			return []string{"java", name}
		},
	})

	RegisterLanguage("cpp", &LanguageDef{
		Image: "gcc:5",
		CompileCommand: func(name string) []string {
			return []string{"g++", name + ".cpp", "-o", name}
		},
		RunCommand: func(name string) []string {
			return []string{"./" + name}
		},
	})

	RegisterLanguage("py", &LanguageDef{
		Image: "python",
		RunCommand: func(name string) []string {
			return []string{"python", name + ".py"}
		},
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

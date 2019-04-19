package main

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/docker/docker/client"
)

func TestCanDetectFileType(t *testing.T) {
	name, ft := detectType("foo.java")
	if name != "foo" {
		t.Errorf("Expected name to be foo but got %v", name)
	}
	if ft != "java" {
		t.Errorf("Expected file type to be java but got %v", ft)
	}
}

var cli *DockerClient

const TEST_FILES = "test-files/0"

func RunTestCase(t *testing.T, testFile string, expected SubmissionStatus) {
	if cli == nil {
		RegisterLanguage("java", &LanguageDef{
			Image:          "openjdk:13-jdk-alpine",
			CompileCommand: []string{"javac", "%f"},
			RunCommand:     []string{"java", "%n"},
			Artifacts:      []string{"%n.class"},
		})

		RegisterLanguage("cpp", &LanguageDef{
			Image:          "gcc:5",
			CompileCommand: []string{"g++", "%f", "-o", "%n"},
			RunCommand:     []string{"./%n"},
			Artifacts:      []string{"%n"},
		})

		RegisterLanguage("py", &LanguageDef{
			Image:      "python",
			RunCommand: []string{"python", "%f"},
		})

		apiClient, err := client.NewClientWithOpts(client.WithVersion("1.39"))
		if err != nil {
			t.Error(err)
		}
		cli = &DockerClient{
			ctx: context.Background(),
			cli: apiClient,
		}

		if err := cli.PullAll(); err != nil {
			t.Error(err)
		}
	}

	result, err := Test(*cli, TEST_FILES, testFile, filepath.Join(TEST_FILES, "/problem"))
	if err != nil {
		t.Error(err)
	}

	if result != expected {
		t.Errorf("Expected result ok but got %v", result)
	}
}

func TestCanTestValidPrograms(t *testing.T) {
	RunTestCase(t, "ok.java", Ok)
}

func TestCanTestIncorrectPrograms(t *testing.T) {
	RunTestCase(t, "fail.java", Wrong)
}

func TestCanHandleRuntimeErrors(t *testing.T) {
	RunTestCase(t, "crash.java", Exception)
}

func TestCanTestInValidPrograms(t *testing.T) {
	RunTestCase(t, "error.java", CompileError)
}

func TestCanHandleInfiniteLoops(t *testing.T) {
	RunTestCase(t, "tle.java", TimeLimit)
}

func TestCanHandleCpp(t *testing.T) {
	RunTestCase(t, "ok.cpp", Ok)
}

func TestCanHandlePython(t *testing.T) {
	RunTestCase(t, "ok.py", Ok)
}

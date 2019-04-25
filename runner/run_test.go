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

const TestFiles = "test-files"

func RunTestCase(t *testing.T, testFile string, expected SubmissionStatus) {
	if cli == nil {
		apiClient, err := client.NewClientWithOpts(client.WithVersion("1.39"))
		if err != nil {
			t.Error(err)
		}
		cli = &DockerClient{
			ctx: context.Background(),
			cli: apiClient,
		}

		loadLanguages(cli)
	}

	result, err := Test(cli, TestFiles+"/0", testFile, filepath.Join(TestFiles, "/problem"))
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

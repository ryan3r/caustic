package runner

import (
	"testing"
	"context"
	"os"
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

func TestCanCompileJava(t *testing.T) {
	os.Chdir("test-files")

	if err := compile("ok.java"); err != nil {
		t.Errorf("ok.java failed to compile (should have): %v", err)
	}
}

func TestCompileCanHandleErrors(t *testing.T) {
	if err := compile("error.java"); err == nil {
		t.Errorf("error.java compiled (should not have): %v", err)
	}
}

func TestRunCanGetOutput(t *testing.T) {
	errors := make(chan error, 1)
	output := make(chan string, 1)

	go run(context.Background(), "ok.java", output, errors)

	select {
	case <- output:
	case err := <- errors:
		t.Errorf("An error occured while running %v", err)
	}
}

func TestRunCanHandleErrors(t *testing.T) {
	errors := make(chan error, 1)
	output := make(chan string, 1)

	go run(context.Background(), "invalid.java", output, errors)

	select {
	case <- output:
		t.Errorf("No errors occured while running (expected one)")
	case <- errors:
	}
}

func TestCanTestValidPrograms(t *testing.T) {
	if result := Test("ok.java", "1 2 3"); result != "ok" {
		t.Errorf("Expected result ok but got %v", result)
	}
}

func TestCanTestIncorrectPrograms(t *testing.T) {
	if result := Test("fail.java", "1 2 3"); result != "wrong" {
		t.Errorf("Expected result wrong but got %v", result)
	}
}

func TestCanHandleRuntimeErrors(t *testing.T) {
	if result := Test("crash.java", "1 2 3"); result != "exception" {
		t.Errorf("Expected result exception but got %v", result)
	}
}

func TestCanTestInValidPrograms(t *testing.T) {
	if result := Test("error.java", "1 2 3"); result != "compile-error" {
		t.Errorf("Expected result compile-error but got %v", result)
	}
}

func TestCanHandleInfiniteLoops(t *testing.T) {
	if result := Test("tle.java", "1 2 3"); result != "time-limit" {
		t.Errorf("Expected result time-limit but got %v", result)
	}
}

func TestCanHandleCpp(t *testing.T) {
	if result := Test("ok.cpp", "1 2 3"); result != "ok" {
		t.Errorf("Expected result ok but got %v", result)
	}
}

func TestCanHandlePython(t *testing.T) {
	if result := Test("ok.py", "1 2 3"); result != "ok" {
		t.Errorf("Expected result ok but got %v", result)
	}
}
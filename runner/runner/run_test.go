package runner

import (
	"testing"
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
	errors := make(chan error)
	output := make(chan string)

	run("ok.java", output, errors)

	select {
	case <- output:
	case err := <- errors:
		t.Errorf("An error occured while running %v", err)
	}
}

func TestRunCanHandleErrors(t *testing.T) {
	errors := make(chan error)
	output := make(chan string)

	run("invalid.java", output, errors)

	select {
	case <- output:
		t.Errorf("No errors occured while running (expected one)")
	case <- errors:
	}
}

func TestCanTestValidPrograms(t *testing.T) {
	if result := Test("ok.java"); result != "ok" {
		t.Errorf("Expected result ok but got %v", result)
	}
}

func TestCanTestIncorrectPrograms(t *testing.T) {
	if result := Test("fail.java"); result != "wrong" {
		t.Errorf("Expected result wrong but got %v", result)
	}
}

func TestCanHandleRuntimeErrors(t *testing.T) {
	if result := Test("crash.java"); result != "exception" {
		t.Errorf("Expected result exception but got %v", result)
	}
}

func TestCanTestInValidPrograms(t *testing.T) {
	if result := Test("error.java"); result != "compile-error" {
		t.Errorf("Expected result compile-error but got %v", result)
	}
}

func TestCanHandleInfiniteLoops(t *testing.T) {
	if result := Test("tle.java"); result != "time-limit" {
		t.Errorf("Expected result time-limit but got %v", result)
	}
}
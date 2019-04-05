package runner

import (
	"testing"
)

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
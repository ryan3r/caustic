package runner

import (
	"os/exec"
)

func compile(filename string) error {
	return exec.Command("javac", filename).Run()
}
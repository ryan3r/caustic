package runner

import (
	"os/exec"
)

func compile(filename string) error {
	return exec.Command("javac", filename).Run()
}

func run(filename string, output chan []byte, errors chan error) {

}
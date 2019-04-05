package runner

import (
	"os/exec"
)

func detectType(filename string) (string, string) {
	return "", ""
}

func compile(filename string) error {
	return exec.Command("javac", filename).Run()
}

func run(filename string, output chan string, errors chan error) {

}

func Test(filename string) string {
	return ""
}
package runner

import (
	"context"
	"errors"
	"log"
	"os/exec"
	"strings"
	"time"
)

// Detect the filetype and name of file
func detectType(filename string) (string, string) {
	idx := strings.LastIndex(filename, ".")
	return filename[:idx], filename[idx+1:]
}

// Compile a file
func compile(filename string) error {
	compiler := "javac"
	_, ft := detectType(filename)

	if ft == "py" {
		return nil
	}

	if ft == "cpp" {
		compiler = "g++"
	}

	log.Printf("Compiling %v as %v\n", filename, ft)
	return exec.Command(compiler, filename).Run()
}

// Run a program
func run(ctx context.Context, filename string, output chan string, errs chan error) {
	name, ft := detectType(filename)
	var cmd *exec.Cmd

	switch ft {
	case "java":
		cmd = exec.CommandContext(ctx, "java", name)
	case "cpp":
		cmd = exec.CommandContext(ctx, "./a.out")
	case "py":
		cmd = exec.CommandContext(ctx, "python", filename)
	default:
		log.Printf("Error unknown filetype %v\n", ft)
		errs <- errors.New("Unknown filetype")
		return
	}

	log.Printf("Running %v as %v", filename, ft)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Error running %v: %v\n", filename, err.Error())
		errs <- err
	} else {
		log.Printf("Completed %v no errors\n", filename)
		output <- string(out)
	}
}

// Test will compile, run and check a program
func Test(filename string, expected string) SubmissionStatus {
	if err := compile(filename); err != nil {
		return CompileError
	}

	errors := make(chan error, 1)
	output := make(chan string, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	run(ctx, filename, output, errors)
	cancel()

	select {
	case out := <-output: // process exited on time w/o errors
		if strings.Trim(out, "\r\n\t ") == expected {
			return Ok
		}
		return Wrong
	case err := <-errors: // process crashed or was killed
		if err.Error() == "signal: killed" {
			return TimeLimit
		}
		return Exception
	}
}

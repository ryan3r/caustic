package runner

import (
	"os/exec"
	"strings"
	"context"
	"time"
	"errors"
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

	return exec.Command(compiler, filename).Run()
}

// Run a program
func run(ctx context.Context, filename string, output chan string, errs chan error) {
	name, ft := detectType(filename)
	var cmd *exec.Cmd;

	switch ft {
	case "java":
		cmd = exec.CommandContext(ctx, "java", name)
	case "cpp":
		cmd = exec.CommandContext(ctx, "./a.out")
	case "py":
		cmd = exec.CommandContext(ctx, "python", filename)
	default:
		errs <- errors.New("Unknown filetype")
		return
	}

	out, err := cmd.CombinedOutput()

	if err != nil {
		errs <- err
	} else {
		output <- string(out)
	}
}

// Compile, run and check a program
func Test(filename string, expected string) string {
	if err := compile(filename); err != nil {
		return "compile-error"
	}

	errors := make(chan error, 1)
	output := make(chan string, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	run(ctx, filename, output, errors)
	cancel()

	select {
	case out := <- output: // process exited on time w/o errors
		if strings.Trim(out, "\r\n\t ") == expected {
			return "ok"
		} else {
			return "wrong"
		}
	case err := <- errors: // process crashed or was killed
		if err.Error() == "signal: killed" {
			return "time-limit"
		} else {
			return "exception"
		}
	}
}
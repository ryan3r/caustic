package runner

import (
	"os/exec"
	"strings"
	"context"
	"time"
)

func detectType(filename string) (string, string) {
	idx := strings.LastIndex(filename, ".")
	return filename[:idx], filename[idx+1:]
}

func compile(filename string) error {
	return exec.Command("javac", filename).Run()
}

func run(ctx context.Context, filename string, output chan string, errors chan error) {
	name, _ := detectType(filename)
	out, err := exec.CommandContext(ctx, "java", name).CombinedOutput()

	if err != nil {
		errors <- err
	} else {
		output <- string(out)
	}
}

func Test(filename string, expected string) string {
	if err := compile(filename); err != nil {
		return "compile-error"
	}

	errors := make(chan error, 1)
	output := make(chan string, 1)

	ctx, _ := context.WithTimeout(context.Background(), 2 * time.Second)
	run(ctx, filename, output, errors)

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
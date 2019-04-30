package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	languageDefs map[string]*LanguageDef
	// ErrUnknownFileType means the filetype is unknown
	ErrUnknownFileType = errors.New("Unknown filetype")
	// ErrExitStatusError is for non-zero exit statuses
	ErrExitStatusError = errors.New("Program exited with non-zero exit status")
	// ErrTimeLimit is for time limit exceeded
	ErrTimeLimit = errors.New("Program took too long to run")
)

// LanguageDef defines how to handle a file type
type LanguageDef struct {
	Image          string   `json:"image"`
	CompileCommand []string `json:"compile"`
	RunCommand     []string `json:"run"`
	Artifacts      []string `json:"artifacts"`
}

// Detect the filetype and name of file
func detectType(fileName string) (string, string) {
	idx := strings.LastIndex(fileName, ".")
	if idx == -1 {
		return fileName, ""
	}
	return fileName[:idx], fileName[idx+1:]
}

func expandTemplate(template []string, fileName, name string) []string {
	templated := append([]string(nil), template...)

	for i, orig := range templated {
		tmp := strings.ReplaceAll(orig, "%n", name)
		templated[i] = strings.ReplaceAll(tmp, "%f", fileName)
	}

	return templated
}

// Compile a file
func Compile(cli *DockerClient, problemDir, fileName string, logFile *os.File) error {
	name, ft := detectType(fileName)
	def := languageDefs[ft]

	if def == nil {
		return ErrUnknownFileType
	}

	// This language is interpreted
	if len(def.CompileCommand) == 0 {
		return nil
	}

	compileCommand := expandTemplate(def.CompileCommand, fileName, name)

	compileCtr := Container{
		Docker:     cli,
		Image:      def.Image,
		Cmd:        compileCommand,
		WorkingDir: "/mnt",
		Out:        logFile,
	}

	if err := compileCtr.BindDir(problemDir, "/mnt", false); err != nil {
		return err
	}

	if err := compileCtr.Run(); err != nil {
		return err
	}

	return compileCtr.Wait()
}

// Runner is a runner for a single test
type Runner struct {
	problemDir string
	fileName   string
	container  *Container
	timeLimit  time.Duration
}

// NewRunner creates the container for running tests
func NewRunner(cli *DockerClient, problemDir, fileName string, timeLimit time.Duration) (*Runner, error) {
	_, ft := detectType(fileName)
	def := languageDefs[ft]

	testCtr := &Container{
		Docker:     cli,
		Image:      def.Image,
		Cmd:        []string{"sleep", "1000000000000"},
		WorkingDir: "/mnt",
		Out:        os.Stdout,
		ReadOnly:   true,
	}

	if err := testCtr.BindDir(problemDir, "/mnt", true); err != nil {
		return nil, err
	}

	if err := testCtr.Run(); err != nil {
		return nil, err
	}

	return &Runner{
		problemDir: problemDir,
		fileName:   fileName,
		container:  testCtr,
		timeLimit:  timeLimit,
	}, nil
}

// Run the submission with the test case
func (r *Runner) Run(in io.Reader, out io.Writer) error {
	name, ft := detectType(r.fileName)
	def := languageDefs[ft]

	exec := ContainerExec{
		Container: r.container,
		Cmd:       expandTemplate(def.RunCommand, r.fileName, name),
		In:        in,
		Out:       out,
	}

	exec.Run()
	cancelTimer := exec.StartKillTimer(r.timeLimit)

	err := <-exec.ExitC
	if err == ErrTimeLimit {
		r.container.Run()
	}
	cancelTimer <- true
	return err
}

// Close stops the runner container
func (r *Runner) Close() error {
	return r.container.Stop()
}

// Test will compile, run and check a program
func Test(cli *DockerClient, problemDir, fileName, solutionDir string) (SubmissionStatus, error) {
	defer cleanUpArtifacts(problemDir, fileName)

	logFile, err := os.Create(filepath.Join(problemDir, "log.txt"))
	if err != nil {
		return RunnerError, err
	}
	defer logFile.Close()

	logFile.Write([]byte("Compile:\n"))

	// Compile the solution
	err = Compile(cli, problemDir, fileName, logFile)
	if err == ErrExitStatusError {
		logFile.Write([]byte("Status: Compile Error"))
		return CompileError, nil
	} else if err != nil {
		logFile.Write([]byte("Status: Runner error\nError: Failed to run the compiler container: " + err.Error()))
		return RunnerError, err
	}

	// Run the submissions
	tests, err := ioutil.ReadDir(solutionDir)
	if err != nil {
		logFile.Write([]byte("Status: Runner error\nError: Failed to open solution directory: " + err.Error()))
		return RunnerError, err
	}

	runner, err := NewRunner(cli, problemDir, fileName, 3*time.Second)
	if err != nil {
		logFile.Write([]byte("Status: Runner error\nError: Failed to create runner container: " + err.Error()))
		return RunnerError, err
	}
	defer runner.Close()

	for _, file := range tests {
		name, fileType := detectType(file.Name())

		if fileType != "in" {
			continue
		}

		logFile.Write([]byte("Running " + file.Name() + ":\n"))

		fileIn, err := os.Open(filepath.Join(solutionDir, file.Name()))
		if err != nil {
			logFile.Write([]byte("Status: Runner error\nError: Failed to load input file " + file.Name() + ": " + err.Error()))
			return RunnerError, err
		}
		defer fileIn.Close()

		outBuffer := bytes.NewBufferString("")

		if err := runner.Run(fileIn, outBuffer); err != nil {
			if err == ErrExitStatusError {
				logFile.Write([]byte("Status: Exception"))
				return Exception, nil
			} else if err == ErrTimeLimit {
				logFile.Write([]byte("Status: Time Limit Exceeded"))
				return TimeLimit, nil
			} else {
				logFile.Write([]byte("Status: Runner error\nError: Failed to run submission:" + err.Error() + "\n"))
				return RunnerError, err
			}
		}

		// Verify the output of the submission
		outFile, err := ioutil.ReadFile(filepath.Join(solutionDir, name+".out"))
		if err != nil {
			logFile.Write([]byte("Status: Runner error\nError: Failed to load answer " + name + ".out: " + err.Error() + "\n"))
			return RunnerError, err
		}

		expectedOut := strings.Trim(string(outFile), "\r\n\t ")
		solutionOut := strings.Trim(string(outBuffer.String()), "\r\n\t ")

		logFile.Write([]byte(solutionOut + "\n"))

		if expectedOut != solutionOut {
			logFile.Write([]byte("Status: Wrong Answer"))
			return Wrong, nil
		}
	}

	logFile.Write([]byte("Status: Accepted"))
	return Ok, nil
}

func cleanUpArtifacts(problemDir, fileName string) {
	name, ft := detectType(fileName)
	def := languageDefs[ft]

	if def == nil || def.Artifacts == nil {
		return
	}

	for _, artifact := range expandTemplate(def.Artifacts, fileName, name) {
		if err := os.RemoveAll(filepath.Join(problemDir, artifact)); err != nil {
			fmt.Printf("Warn: Failed to clean up artifacts (%v, %v)\n", problemDir, artifact)
		}
	}
}

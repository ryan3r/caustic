package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
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
	// DefaultTimeLimit is the default time limit for problems
	DefaultTimeLimit = flag.Int("timelimit", 3, "The time limit to use for problems without a problem.json")
)

// LanguageDef defines how to handle a file type
type LanguageDef struct {
	Image          string   `json:"image"`
	CompileCommand []string `json:"compile"`
	RunCommand     []string `json:"run"`
	Artifacts      []string `json:"artifacts"`
}

// ProblemDef defines a problem
type ProblemDef struct {
	Time int    `json:"time"`
	Name string `json:"name"`
}

// loadProblem loads a problem from a json file
func loadProblem(problemDir string) (ProblemDef, error) {
	problem := ProblemDef{
		Time: *DefaultTimeLimit,
		Name: "Unnamed problem",
	}

	langFile, err := ioutil.ReadFile(filepath.Join(problemDir, "problem.json"))
	if err != nil {
		return problem, err
	}

	if err := json.Unmarshal(langFile, &problem); err != nil {
		return problem, err
	}

	return problem, nil
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
func Compile(cli *DockerClient, problemDir, fileName, ft string, logFile *os.File) error {
	name, _ := detectType(fileName)
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
	ft         string
}

// NewRunner creates the container for running tests
func NewRunner(cli *DockerClient, problemDir, fileName, ft string, timeLimit time.Duration) (*Runner, error) {
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
		ft:         ft,
	}, nil
}

// Run the submission with the test case
func (r *Runner) Run(in io.Reader, out io.Writer) error {
	name, _ := detectType(r.fileName)
	def := languageDefs[r.ft]

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
func Test(cli *DockerClient, problemDir, fileName, solutionDir, ft string) (SubmissionStatus, error) {
	defer cleanUpArtifacts(problemDir, fileName, ft)

	logFile, err := os.Create(filepath.Join(problemDir, "log.txt"))
	if err != nil {
		return RunnerError, err
	}
	defer logFile.Close()

	logFile.Write([]byte("Compile:\n"))

	// Compile the solution
	err = Compile(cli, problemDir, fileName, ft, logFile)
	if err == ErrExitStatusError {
		logFile.Write([]byte("Status: Compile Error\n"))
		return CompileError, nil
	} else if err != nil {
		logFile.Write([]byte("Status: Runner error\nError: Failed to run the compiler container: " + err.Error() + "\n"))
		return RunnerError, err
	}

	// Run the submissions
	tests, err := ioutil.ReadDir(solutionDir)
	if err != nil {
		logFile.Write([]byte("Status: Runner error\nError: Failed to open solution directory: " + err.Error() + "\n"))
		return RunnerError, err
	}

	problemDef, err := loadProblem(problemDir)
	if err != nil {
		logFile.Write([]byte("Failed to open solution definition (using defaults): " + err.Error() + "\n"))
	}

	runner, err := NewRunner(cli, problemDir, fileName, ft, time.Duration(problemDef.Time)*time.Second)
	if err != nil {
		logFile.Write([]byte("Status: Runner error\nError: Failed to create runner container: " + err.Error() + "\n"))
		return RunnerError, err
	}
	defer runner.Close()

	for _, file := range tests {
		name, fileType := detectType(file.Name())

		if fileType != "in" {
			continue
		}

		logFile.Write([]byte("Running " + file.Name() + ":\n"))

		// Get the inputs
		fileIn, err := os.Open(filepath.Join(solutionDir, file.Name()))
		if err != nil {
			logFile.Write([]byte("Status: Runner error\nError: Failed to load input file " + file.Name() + ": " + err.Error() + "\n"))
			return RunnerError, err
		}
		defer fileIn.Close()

		outBuffer := bytes.NewBufferString("")

		// Run the code
		if err := runner.Run(fileIn, outBuffer); err != nil {
			if err == ErrExitStatusError {
				logFile.Write([]byte("Status: Exception\n"))
				return Exception, nil
			} else if err == ErrTimeLimit {
				logFile.Write([]byte("Status: Time Limit Exceeded\n"))
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
			logFile.Write([]byte("Status: Wrong Answer\n"))
			return Wrong, nil
		}
	}

	logFile.Write([]byte("Status: Accepted"))
	return Ok, nil
}

func cleanUpArtifacts(problemDir, fileName, ft string) {
	name, _ := detectType(fileName)
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

package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"syscall"
)

var buildDate = "development"
var buildCommit = "development"
var buildVersion = "development"
var versionString = fmt.Sprintf(strings.TrimSpace(`
Version:     %s
Go version:  %s
Build date:  %s
Git commit:  %s
`),
	buildVersion,
	runtime.Version(),
	buildDate,
	buildCommit,
)

var usage = `envdir.

Run a command with environment variables
specified by the files in a directory.

Usage:
  envdir --version
  envdir --help
  envdir [-i] <directory> <command> [<arguments>...]

Arguments:
  <directory>  The directory of files representing environment variables.
  <command>    The command to run.
  <arguments>  The arguments of the command to run.

Options:
  -i, --ignore-environment  Start with an empty environment.
  --version  Show version.
  --help     Show help.

Interface:
  Each filename in <directory> is the name of an environment variable.
  The contents of the file is the value of the environment variable.
  The last newline of each file is ignored.
  If the file is empty (containing only 0 bytes or 1 newline),
    that environment variable is unset.

  envdir exits 111 if:
   * The directory's files can't be read
   * A filename contains "="
   * A file contains the null character
   * The command can't be run
`

const exitCodeUnsuccessful = 111

type environment map[string]string

func (e environment) Strings() []string {
	result := make([]string, len(e))
	index := 0
	for key, value := range e {
		result[index] = fmt.Sprintf("%s=%s", key, value)
		index++
	}
	return result
}

func makeEnvironmentMap(pairs []string) environment {
	environmentMap := make(environment)
	for _, pair := range pairs {
		keyValue := strings.SplitN(pair, "=", 2)
		key := strings.Join(keyValue[:1], "")
		value := strings.Join(keyValue[1:], "=")
		environmentMap[key] = value
	}
	return environmentMap
}

func main() {
	arguments, _ := docopt.Parse(usage, nil, true, versionString, true)
	freshEnvironment := arguments["--ignore-environment"].(bool)
	directory := arguments["<directory>"].(string)
	commandName := arguments["<command>"].(string)
	commandArguments := arguments["<arguments>"].([]string)
	environmentMap := getEnvironmentVariables(directory, freshEnvironment)
	environmentList := environmentMap.Strings()
	runCommand(commandName, commandArguments, environmentList)
}

func runCommand(
	commandName string,
	commandArguments []string,
	environmentList []string,
) {
	command := exec.Command(commandName, commandArguments...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = environmentList
	if err := command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			os.Exit(waitStatus.ExitStatus())
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(exitCodeUnsuccessful)
		}
	}
	os.Exit(0)
}

func getEnvironmentVariables(directory string, freshEnvironment bool) environment {
	var environmentMap environment
	if freshEnvironment {
		environmentMap = make(environment)
	} else {
		environmentMap = makeEnvironmentMap(os.Environ())
	}
	contents, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCodeUnsuccessful)
	}
	for _, fileInfo := range contents {
		if fileInfo.IsDir() {
			continue
		}
		fileName := fileInfo.Name()
		if strings.ContainsRune(fileName, '=') {
			fmt.Fprintf(os.Stderr, "Error: Filenames in the environment directory must not contain \"=\": %s\n", fileName)
			os.Exit(exitCodeUnsuccessful)
		}
		if fileInfo.Size() == 0 {
			delete(environmentMap, fileName)
			continue
		}
		fileLocation := path.Join(directory, fileName)
		fileData, err := ioutil.ReadFile(fileLocation)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(exitCodeUnsuccessful)
		}
		fileString := string(fileData)
		if stringHasNullCharacter(fileString) {
			fmt.Fprintf(os.Stderr, "Error: %s contains a null character\n", fileName)
			os.Exit(exitCodeUnsuccessful)
		}
		fileString = trimLastNewline(fileString)
		if len(fileString) == 0 {
			delete(environmentMap, fileName)
			continue
		}
		environmentMap[fileName] = fileString
	}
	return environmentMap
}

func stringHasNullCharacter(s string) bool {
	i := strings.IndexByte(s, '\x00')
	return i != -1
}

func trimLastNewline(s string) string {
	if len(s) >= 2 && s[len(s)-2:] == "\r\n" {
		return s[:len(s)-2]
	}
	if len(s) >= 1 {
		lastRune := s[len(s)-1:]
		if lastRune == "\n" {
			return s[:len(s)-1]
		}
		if lastRune == "\r" {
			return s[:len(s)-1]
		}
	}
	return s
}

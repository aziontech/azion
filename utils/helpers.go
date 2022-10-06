package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/aziontech/azion-cli/pkg/contracts"
)

const shell = "/bin/sh"

func CleanDirectory(dir string) error {

	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("%w - %s", ErrorCleaningDirectory, dir)
	}

	return nil
}

func IsDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		// Dir does not exist
		if errors.Is(err, os.ErrNotExist) {
			return true, nil
		}
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func LoadEnvVarsFromFile(varsFileName string) ([]string, error) {
	if _, err := os.Stat(varsFileName); errors.Is(err, os.ErrNotExist) {
		// Ignore error if not specified
		if varsFileName == "" {
			return nil, nil
		}
		return nil, err
	}

	f, err := os.Open(varsFileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileScan := bufio.NewScanner(f)
	fileVars := make([]string, 0)

	for fileScan.Scan() {
		fileVars = append(fileVars, fileScan.Text())
	}

	if err := fileScan.Err(); err != nil {
		return nil, err
	}

	return fileVars, nil
}

// RunCommandWithOutput returns the stringified command output, it's exit code and any errors
// Commands that exit with exit codes > 0 will return a non-nil error
func RunCommandWithOutput(envVars []string, comm string) (string, int, error) {
	command := exec.Command(shell, "-c", comm)
	if len(envVars) > 0 {
		command.Env = os.Environ()
		command.Env = append(command.Env, envVars...)
	}

	out, err := command.CombinedOutput()
	exitCode := command.ProcessState.ExitCode()

	return string(out), exitCode, err
}

func RunCommand(envVars []string, comm string) error {
	command := exec.Command(shell, "-c", comm)
	//Load environment variables (if they exist)
	if len(envVars) > 0 {
		command.Env = os.Environ()
		command.Env = append(command.Env, envVars...)
	}
	err := command.Run()
	if err != nil {
		return ErrorRunningCommand
	}

	return nil
}

func GetWorkingDir() (string, error) {
	pathWorkingDir, err := os.Getwd()
	if err != nil {
		return "", ErrorInternalServerError
	}
	return pathWorkingDir, nil
}

func ResponseToBool(response string) (bool, error) {

	response = strings.TrimSpace(response)

	if strings.ToLower(response) == "yes" {
		return true, nil
	}
	if strings.ToLower(response) == "no" || response == "" {
		return false, nil
	}

	return false, ErrorInvalidOption
}

func GetAzionJsonContent() (*contracts.AzionApplicationOptions, error) {
	path, err := GetWorkingDir()
	if err != nil {
		return nil, err
	}
	jsonConf := path + "/azion/azion.json"
	file, err := os.ReadFile(jsonConf)
	if err != nil {
		fmt.Println(&jsonConf)
		return nil, ErrorOpeningAzionJsonFile
	}

	conf := &contracts.AzionApplicationOptions{}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		return nil, ErrorUnmarshalAzionJsonFile
	}

	return conf, nil
}

func WriteAzionJsonContent(conf *contracts.AzionApplicationOptions) error {
	path, err := GetWorkingDir()
	if err != nil {
		return err
	}
	jsonConf := path + "/azion/azion.json"

	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return ErrorMarshalAzionJsonFile
	}

	err = os.WriteFile(jsonConf, data, 0644)
	if err != nil {
		return ErrorWritingAzionJsonFile
	}

	return nil
}

//Returns the correct error message for each HTTP Status code
func ErrorPerStatusCode(httpResp *http.Response, err error) error {

	// when the CLI times out, probably due to SSO communication, httpResp is null and/or http status is 500;
	// that's why we need this verification first
	if httpResp == nil || httpResp.StatusCode >= 500 {
		return checkStatusCode500Error(err)
	}

	statusCode := httpResp.StatusCode

	switch statusCode {

	case 400:
		return checkStatusCode400Error(httpResp)

	case 401:
		return ErrorToken401

	case 403:
		return ErrorForbidden403

	case 404:
		return ErrorNotFound404

	default:
		return err

	}
}

// checks varying errors that may occur when status code is 500
func checkStatusCode500Error(err error) error {

	if strings.Contains(err.Error(), "Client.Timeout") {
		return ErrorTimeoutAPICall
	}

	return ErrorInternalServerError
}

//read the body of the response and returns its content
func checkStatusCode400Error(httpResp *http.Response) error {
	responseBody, _ := ioutil.ReadAll(httpResp.Body)
	return fmt.Errorf("%s", responseBody)
}

//TODO ADD STATUS CODE

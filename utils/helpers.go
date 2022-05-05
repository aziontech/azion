package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aziontech/azion-cli/pkg/contracts"
)

const shell = "/bin/sh"

func ConvertIdsToInt(ids ...string) ([]int64, error) {
	converted_ids := make([]int64, len(ids))
	for index, id := range ids {
		converted_id, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		converted_ids[index] = int64(converted_id)
	}

	return converted_ids, nil

}

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

func GetAzionJsonContent() (*contracts.AzionJsonData, error) {
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

	conf := &contracts.AzionJsonData{}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		return nil, ErrorUnmarshalAzionJsonFile
	}

	return conf, nil
}

func WriteAzionJsonContent(conf *contracts.AzionJsonData) error {
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

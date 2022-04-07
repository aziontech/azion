package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
)

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

func LoadEnvVars(varsFileName string) ([]string, error) {
	// If no .env file, ignore and return
	if _, err := os.Stat(varsFileName); errors.Is(err, os.ErrNotExist) {
		return nil, nil
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

func RunCommand(envVars []string, comm string) error {
	const shell = "/bin/bash"
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

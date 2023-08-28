package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
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

func CommandRunInteractive(f *cmdutil.Factory, envVars []string, comm string) error {

	cmd := exec.Command(shell, "-c", comm)
	cmd.Stdin = f.IOStreams.In
	cmd.Stdout = f.IOStreams.Out
	cmd.Stderr = f.IOStreams.Err
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// RunCommandStreamOutput executes the provived command while streaming its logs (stdout+stderr) directly to terminal
func RunCommandStreamOutput(out io.Writer, envVars []string, comm string) error {
	command := exec.Command(shell, "-c", comm)
	if len(envVars) > 0 {
		command.Env = os.Environ()
		command.Env = append(command.Env, envVars...)
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	multi := io.MultiReader(stdout, stderr)

	// start the command after having set up the pipe
	if err := command.Start(); err != nil {
		return fmt.Errorf(ErrorRunningCommandStream.Error(), err)
	}

	// read command's stdout line by line
	in := bufio.NewScanner(multi)

	for in.Scan() {
		fmt.Fprintf(out, "%s\n", in.Text())
	}
	if err := in.Err(); err != nil {
		return fmt.Errorf(ErrorRunningCommandStream.Error(), err)
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

	_, err = os.Stat(path + "/azion/azion.json")
	if err != nil {
		logger.Debug("Error reading stats of azion.json file", zap.Error(err))
		return nil, ErrorOpeningAzionJsonFile
	}

	jsonConf := path + "/azion/azion.json"
	file, err := os.ReadFile(jsonConf)
	if err != nil {
		logger.Debug("Error reading azion.json file", zap.Error(err))
		return nil, ErrorOpeningAzionJsonFile
	}

	conf := &contracts.AzionApplicationOptions{}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		logger.Debug("Error reading unmarshalling azion.json file", zap.Error(err))
		return nil, ErrorUnmarshalAzionJsonFile
	}

	return conf, nil
}

func GetAzionJsonSimple() (*contracts.AzionApplicationSimple, error) {
	path, err := GetWorkingDir()
	if err != nil {
		return nil, err
	}
	jsonConf := path + "/azion/azion.json"
	file, err := os.ReadFile(jsonConf)
	if err != nil {
		return nil, ErrorOpeningAzionJsonFile
	}

	conf := &contracts.AzionApplicationSimple{}
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

// Returns the correct error message for each HTTP Status code
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

// read the body of the response and returns a personalized error or the body if the error is not identified
func checkStatusCode400Error(httpResp *http.Response) error {
	responseBody, _ := io.ReadAll(httpResp.Body)
	if err := checkNoProduct(string(responseBody)); err != nil {
		return err
	}
	if err := checkTlsVersion(string(responseBody)); err != nil {
		return err
	}
	if err := checkOriginlessCacheSettings(string(responseBody)); err != nil {
		return err
	}
	if err := checkDetail(string(responseBody)); err != nil {
		return err
	}
	if err := checkOrderField(string(responseBody)); err != nil {
		return err
	}

	return fmt.Errorf("%s", responseBody)
}

func checkNoProduct(body string) error {
	if strings.Contains(body, "user_has_no_product") {
		product := gjson.Get(body, "user_has_no_product")
		return fmt.Errorf("%w: %s", ErrorProductNotOwned, product.String())
	}
	return nil
}

func checkOriginlessCacheSettings(body string) error {
	if strings.Contains(body, "originless_cache_settings") {
		msgorigin := gjson.Get(body, "originless_cache_settings")
		return fmt.Errorf("%s", msgorigin.String())
	}
	return nil
}

func checkTlsVersion(body string) error {
	if strings.Contains(body, "minimum_tls_version") {
		return msg.ErrorMinTlsVersion
	}
	return nil
}

func checkDetail(body string) error {
	if strings.Contains(body, "detail") {
		msgDetail := gjson.Get(body, "detail")
		return fmt.Errorf("%s", msgDetail.String())
	}
	return nil
}

func checkOrderField(body string) error {
	if strings.Contains(body, "invalid_order_field") {
		msgDetail := gjson.Get(body, "invalid_order_field")
		return fmt.Errorf("%s", msgDetail.String())
	}
	return nil
}

func CreateVersionID() string {
	return time.Now().Format("20060102150405")
}

func AskForInput(in io.ReadCloser, out io.Writer, message string) (response string) {
	fmt.Fprintf(out, "%s: ", message)
	fmt.Fscanln(in, &response)
	return response
}

func TruncateString(str string) string {
	if len(str) > 30 {
		return str[:30] + "..."
	}
	return str
}

func IsEmpty(str string) bool {
	return len(str) < 1
}

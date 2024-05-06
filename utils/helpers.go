package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/manifoldco/promptui"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

const shell = "/bin/sh"

var (
	NameTaken = []string{"already taken", "name taken", "name already in use", "already in use", "already exists", "with the name", "409 Conflict"}
)

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

// CommandRunInteractive runs a command interactively.
func CommandRunInteractiveWithOutput(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
	cmd := exec.Command(shell, "-c", comm)
	if len(envVars) > 0 {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, envVars...)
	}
	var stdoutBuffer bytes.Buffer

	if !f.Silent {
		cmd.Stdin = f.IOStreams.In
		cmd.Stdout = &stdoutBuffer
	}

	cmd.Stderr = f.IOStreams.Err

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	output := stdoutBuffer.String()

	return output, nil
}

// CommandRunInteractive runs a command interactively.
func CommandRunInteractive(f *cmdutil.Factory, comm string) error {
	cmd := exec.Command(shell, "-c", comm)

	if !f.Silent {
		cmd.Stdin = f.IOStreams.In
		cmd.Stdout = f.IOStreams.Out
	}

	cmd.Stderr = f.IOStreams.Err

	return cmd.Run()
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

	case 409:
		return ErrorNameInUse

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
	if err := checkNameInUse(string(responseBody)); err != nil {
		return err
	}

	result := strings.ReplaceAll(string(responseBody), "{", "")
	result = strings.ReplaceAll(result, "}", "")
	result = strings.ReplaceAll(result, "[", "")
	result = strings.ReplaceAll(result, "]", "")

	return fmt.Errorf("%s", result)
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
		return ErrorMinTlsVersion
	}
	return nil
}

func checkNameInUse(body string) error {
	if strings.Contains(body, "name_already_in_use") || strings.Contains(body, "bucket name is already in use") || containsErrorMessageNameTaken(body) {
		return ErrorNameInUse
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

func TruncateString(str string) string {
	if len(str) > 30 {
		return str[:30] + "..."
	}
	return str
}

// IsEmpty returns true when the string is empty
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case []int:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case map[string]int:
		return len(v) == 0
	case map[string]string:
		return len(v) == 0
	case *string:
		return v == nil || *v == ""
	case *int:
		return v == nil
	case *bool:
		return v == nil
	case *float64:
		return v == nil
	case *[]int:
		return v == nil || len(*v) == 0
	case *[]string:
		return v == nil || len(*v) == 0
	case *map[string]int:
		return v == nil || len(*v) == 0
	case *map[string]string:
		return v == nil || len(*v) == 0
	}

	return false
}

func GetPackageManager() (string, error) {
	opts := []string{"npm", "yarn"}
	answer := ""
	prompt := &survey.Select{
		Message: "Choose a package manager:",
		Options: opts,
	}
	err := survey.AskOne(prompt, &answer)
	if err != nil {
		return "", err
	}
	return answer, nil
}

func AskInputEmpty(msg string) (string, error) {
	qs := []*survey.Question{
		{
			Name:     "id",
			Prompt:   &survey.Input{Message: msg},
			Validate: survey.MinLength(0),
		},
	}

	answer := ""

	err := survey.Ask(qs, &answer)
	if err == terminal.InterruptErr {
		logger.Error(ErrorCancelledContextInput.Error())
		os.Exit(0)
	} else if err != nil {
		logger.Debug("Error while parsing answer", zap.Error(err))
		return "", ErrorParseResponse
	}

	return answer, nil
}

func AskInput(msg string) (string, error) {
	qs := []*survey.Question{
		{
			Name:     "id",
			Prompt:   &survey.Input{Message: msg},
			Validate: survey.Required,
		},
	}

	answer := ""

	err := survey.Ask(qs, &answer)
	if err == terminal.InterruptErr {
		logger.Error(ErrorCancelledContextInput.Error())
		os.Exit(0)
	} else if err != nil {
		logger.Debug("Error while parsing answer", zap.Error(err))
		return "", ErrorParseResponse
	}

	return answer, nil
}

func AskPassword(msg string) (string, error) {
	qs := []*survey.Question{
		{
			Name:     "id",
			Prompt:   &survey.Password{Message: msg},
			Validate: survey.Required,
		},
	}

	answer := ""

	err := survey.Ask(qs, &answer)
	if err == terminal.InterruptErr {
		logger.Error(ErrorCancelledContextInput.Error())
		os.Exit(0)
	} else if err != nil {
		logger.Debug("Error while parsing answer", zap.Error(err))
		return "", ErrorParseResponse
	}

	return answer, nil
}

func LogAndRewindBody(httpResp *http.Response) error {
	logger.Debug("", zap.Any("Status Code", httpResp.StatusCode))
	logger.Debug("", zap.Any("Headers", httpResp.Header))
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		logger.Debug("Error while reading body of the http response", zap.Error(err))
		return ErrorPerStatusCode(httpResp, err)
	}

	// Convert the body bytes to string
	bodyString := string(bodyBytes)
	logger.Debug("", zap.Any("Body", bodyString))

	// Rewind the response body to the beginning
	httpResp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return nil
}

// FlagINUnmarshalFileJSON
// request interface{} always as a pointer
func FlagFileUnmarshalJSON(path string, request interface{}) error {
	var (
		file *os.File
		err  error
	)

	if path == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(path)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrorOpeningFile, path)
		}
		defer file.Close()
	}

	return cmdutil.UnmarshallJsonFromReader(file, &request)
}

func Select(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func Concat(strs ...string) string {
	var sb strings.Builder
	for i := 0; i < len(strs); i++ {
		sb.WriteString(strs[i])
	}
	return sb.String()
}

// Confirm is a function that provides a confirmation prompt to the user.
// It takes three parameters:
// - globalFlagAll: a boolean flag to skip the confirmation and return true directly.
// - msg: the message to display as part of the confirmation prompt.
// - defaultYes: a boolean flag indicating whether pressing enter should default to 'yes'.
func Confirm(globalFlagAll bool, msg string, defaultYes bool) bool {
	if globalFlagAll {
		return true
	}

	fmt.Printf("🤔 \x1b[32m%s \x1b[0m", msg)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	confirm := scanner.Text()

	if confirm == "" && !defaultYes {
		return false
	} else if confirm == "" && defaultYes {
		return true
	}

	switch confirm {
	case "y", "Y":
		return true
	case "n", "N":
		return false
	default:
		fmt.Printf("\x1b[33m%s\x1b[0m", "⚠️ Invalid input. Please enter 'y' or 'n'.\n")
		return Confirm(globalFlagAll, msg, defaultYes)
	}
}

func Format(input string) (int, error) {
	numberString := ""
	for _, char := range input {
		if unicode.IsDigit(char) {
			numberString += string(char)
		}
	}

	number, err := strconv.Atoi(numberString)
	if err != nil {
		return 0, err
	}

	return number, nil

}

func containsErrorMessageNameTaken(msg string) bool {
	for _, phrase := range NameTaken {
		if strings.Contains(msg, phrase) {
			return true
		}
	}
	return false
}

func Timestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

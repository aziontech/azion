package upbin

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
)

const (
	azioncli     string = "azioncli"
	repoURL      string = "https://github.com/aziontech/azion-cli.git"
	LastActivity string = "LAST_ACTIVITY"
	unknown      string = "unknown"
)

// Variable factory for unit testing
var (
	notifyFunc           = notify
	needToUpdateFunc     = needToUpdate
	wantToUpdateFunc     = wantToUpdate
	saveLastActivityFunc = saveLastActivity
	prepareURLFunc       = prepareURL
	whichFunc            = which
	openConfigFunc       = openConfig
)

func UpdateBin() error {
	notfy, err := notifyFunc()
	if err != nil {
		return err
	}

	if !notfy {
		return nil
	}

	if !needToUpdateFunc() {
		return nil
	}

	if !wantToUpdateFunc() {
		return saveLastActivityFunc()
	}

	install, err := managersPackagesFunc()
	if err != nil {
		return err
	}

	if install {
		return nil
	}

	filePath, err := whichFunc(azioncli)
	if err != nil {
		return err
	}

	err = replaceBinaryFunc(filePath)
	if err != nil {
		return err
	}

	return nil
}

func GetCurrentVersion() int {
	n, _ := Format(version.BinVersion)
	return n
}

func GetLatestVersion() int {
	r, _ := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: repoURL})

	tags, _ := r.Tags()
	t, _ := latestTag(tags)

	n, _ := Format(t)
	return n
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

type ReferenceIter interface {
	ForEach(func(*plumbing.Reference) error) error
}

// latestTag return value in format refs/tags/0.10.0
func latestTag(tags ReferenceIter) (tag string, err error) {
	var biggerVersionSoFar int = 0

	err = tags.ForEach(func(t *plumbing.Reference) error {
		tagCurrent := t.Name().String() // return this format "refs/tags/0.10.0"

		if !strings.Contains(tagCurrent, "dev") && !strings.Contains(tagCurrent, "beta") {
			current, _ := strconv.Atoi(getNumbersString(tagCurrent))
			if current > biggerVersionSoFar {
				biggerVersionSoFar = current
				tag = tagCurrent
			}
		}

		return err
	})

	return tag, err
}

// getNumbersString get numbers from a string
func getNumbersString(str string) string {
	var currentNumber string
	for _, char := range str {
		if unicode.IsDigit(char) {
			currentNumber += string(char)
		}
	}
	return currentNumber
}

func which(command string) (string, error) {
	paths := filepath.SplitList(os.Getenv("PATH"))

	for _, dir := range paths {
		executablePath := filepath.Join(dir, command)
		if _, err := os.Stat(executablePath); err == nil {
			return executablePath, nil
		}
	}

	return "", fmt.Errorf("command not found: %s", command)
}

func needToUpdate() bool {
	return GetLatestVersion() >= GetCurrentVersion()
}

func wantToUpdate() bool {
	prompt := promptui.Select{
		Label: "A new version of 'azioncli' was published. Do you wish to update to the latest version?",
		Items: []string{"Yes", "No"},
	}
	_, result, _ := prompt.Run()
	return result == "Yes"
}

func notify() (bool, error) {
	lastAct, err := getLastActivity()
	if err != nil {
		return false, err
	}

	// Difference is greater than or equal to half a day, notify user
	dateHour, _ := time.Parse(time.RFC3339, lastAct)
	diff := time.Since(dateHour)

	if diff >= 12*time.Hour {
		return true, nil
	} else {
		return false, nil
	}
}

func saveLastActivity() error {
	data, err := openConfig()
	if err != nil {
		return err
	}

	data[LastActivity] = time.Now().Format(time.RFC3339)

	err = writeConfig(data)
	if err != nil {
		return err
	}

	return nil
}

func openConfig() (map[string]interface{}, error) {
	path, err := pathYamlConfig()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func writeConfig(data map[string]interface{}) error {
	content, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}

	path, err := pathYamlConfig()
	if err != nil {
		return err
	}

	err = os.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func pathYamlConfig() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".config", "azioncli.yaml"), err
}

func getLastActivity() (string, error) {
	data, err := openConfigFunc()
	if err != nil {
		return "", err
	}

	lastAct := data[LastActivity]
	return fmt.Sprint(lastAct), nil
}

type SystemArch struct {
	System string
	Arch   string
}

func GetInfoSystem() SystemArch {
	ar := runtime.GOARCH
	os := runtime.GOOS

	switch os {
	case "darwin":
		switch ar {
		case "arm64":
			return SystemArch{System: "Darwin", Arch: "arm64"}
		case "amd64":
			return SystemArch{System: "Darwin", Arch: "x86_64"}
		}
	case "linux":
		switch ar {
		case "386":
			return SystemArch{System: "linux", Arch: "386"}
		case "amd64":
			return SystemArch{System: "linux", Arch: "amd64"}
		case "arm64":
			return SystemArch{System: "linux", Arch: "arm64"}
		case "arm":
			return SystemArch{System: "linux", Arch: "arm"}
		}
	}

	return SystemArch{System: unknown, Arch: unknown}
}

package upbin

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
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

type system string

const (
	azioncli string = "azioncli"

	windows system = "windows"
	darwin  system = "darwin"
	linux   system = "linux"
	unknown system = "unknown"

	repoURL string = "https://github.com/aziontech/azion-cli.git"
	url     string = "https://downloads.azion.com/%s/x86_64/azioncli"

	LastActivity    string = "LAST_ACTIVITY"
	packageAzioncli string = "aziontech/tap/azioncli"
)

// Variable factory for unit testing
var (
	notifyFunc           = notify
	needToUpdateFunc     = needToUpdate
	wantToUpdateFunc     = wantToUpdate
	saveLastActivityFunc = saveLastActivity
	prepareURLFunc       = prepareURL
	managersPackagesFunc = managersPackages
	whichFunc            = which
	downloadFileFunc     = downloadFile
	replaceFileFunc      = replaceFile

	openConfigFunc = openConfig

	packageManagerExistsFunc  = packageManagerExists
	installPackageManagerFunc = installPackageManager
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

	fileURL, err := prepareURLFunc()
	if err != nil {
		return err
	}

	install, err := managersPackagesFunc()
	if err != nil {
		return err
	}

	if install {
		return nil
	}

	filePath, _ := whichFunc(azioncli)

	err = downloadFileFunc(filePath, fileURL)
	if err != nil {
		return err
	}

	err = replaceFileFunc(filePath)
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

// latestTag return value in format refs/tags/v0.10.0
func latestTag(tags ReferenceIter) (tag string, err error) {
	var biggerVersionSoFar int = 0

	err = tags.ForEach(func(t *plumbing.Reference) error {
		tagCurrent := t.Name().String() // return this format "refs/tags/v0.10.0"

		if !strings.Contains(tagCurrent, "dev") {
			versionParts := strings.Split(tagCurrent, ".")

			major := strings.TrimPrefix(versionParts[0], "refs/tags/v")
			minor := versionParts[1]
			patch := versionParts[2]

			current, _ := strconv.Atoi(fmt.Sprintf("%s%s%s", major, minor, patch))

			if current > biggerVersionSoFar {
				biggerVersionSoFar = current
				tag = tagCurrent
			}
		}

		return err
	})

	return tag, err
}

func which(command string) (string, error) {
	path := os.Getenv("PATH")
	paths := filepath.SplitList(path)

	for _, dir := range paths {
		executablePath := filepath.Join(dir, command)
		_, err := os.Stat(executablePath)
		if err == nil {
			return executablePath, nil
		}
	}

	return "", fmt.Errorf("command '%s' not found", command)
}

func GetSystem() system {
	switch system(runtime.GOOS) {
	case linux:
		return linux
	case darwin:
		return darwin
	case windows:
		return windows
	default:
		return unknown
	}
}

func prepareURL() (string, error) {
	sys := GetSystem()
	if sys == unknown {
		return "", errors.New("unknown system")
	}
	return fmt.Sprintf(url, sys), nil
}

func needToUpdate() bool {
	return GetLatestVersion() >= GetCurrentVersion()
}

func wantToUpdate() bool {
	prompt := promptui.Select{
		Label: "Do you want to update 'azioncli'?",
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

func downloadFile(filePath, fileURL string) error {
	response, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func replaceFile(filePath string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	err = os.Rename(filePath, exe)
	if err != nil {
		return err
	}

	return nil
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

	err = ioutil.WriteFile(path, content, 0644)
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

func managersPackages() (bool, error) {
	packageManagers := []string{"brew"}

	var install bool = false
	var manager string
	for _, man := range packageManagers {
		exists := packageManagerExistsFunc(man)
		if exists {
			manager = man
			install = true
			break
		}
	}

	err := installPackageManagerFunc(manager)
	if err != nil {
		return false, err
	}

	return install, nil
}

func packageManagerExists(command string) bool {
	_, err := exec.LookPath(command)
	if err != nil {
		return false
	} else {
		return true
	}
}

func installPackageManager(manager string) error {
	var command string
	switch manager {
	case "brew":
		command = "install"
	}

	cmd := exec.Command(manager, command, packageAzioncli)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

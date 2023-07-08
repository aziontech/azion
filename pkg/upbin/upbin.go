package upbin

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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

type system string

const (
	azioncli string = "azioncli"

	windows system = "windows"
	darwin  system = "darwin"
	linux   system = "linux"
	unknown system = "unknown"

	repoURL string = "https://github.com/aziontech/azion-cli.git"
	url     string = "https://downloads.azion.com/%s/x86_64/azioncli"

	LastActivity string = "LAST_ACTIVITY"
)

type Version struct {
	Major int
	Minor int
	Patch int
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
	var previousVersion Version

	err = tags.ForEach(func(t *plumbing.Reference) error {
		tagCurrent := t.Name().String() // return this format "refs/tags/v0.10.0"
		if !strings.Contains(tagCurrent, "dev") {
			versionParts := strings.Split(tagCurrent, ".")

			majorCurrent, _ := strconv.Atoi(strings.TrimPrefix(versionParts[0], "v"))
			minorCurrent, _ := strconv.Atoi(versionParts[1])
			patchCurrent, _ := strconv.Atoi(versionParts[2])

			currentVersion := Version{
				Major: majorCurrent,
				Minor: minorCurrent,
				Patch: patchCurrent,
			}

			if currentVersion.Major > previousVersion.Major {
				previousVersion = currentVersion
				tag = tagCurrent
			} else if currentVersion.Major == previousVersion.Major && currentVersion.Minor > previousVersion.Minor {
				previousVersion = currentVersion
				tag = tagCurrent
			} else if currentVersion.Major == previousVersion.Major && currentVersion.Minor == previousVersion.Minor && currentVersion.Patch > previousVersion.Patch {
				previousVersion = currentVersion
				tag = tagCurrent
			}
		}

		return err
	})

	if err != nil {
		return tag, err
	}

	return tag, err
}

func Which(command string) (string, error) {
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

func notify() bool {
	lastAct, _ := getLastActivity()

	// Difference is greater than or equal to half a day, notify user
	dateHour, _ := time.Parse(time.RFC3339, lastAct)
	diff := time.Since(dateHour)

	if diff >= 12*time.Hour {
		return true
	} else {
		return false
	}
}

func UpdateBin() error {
	if !notify() {
		fmt.Println("foi notificado recentemente")
		return nil
	}

	if !needToUpdate() {
		return nil
	}

	if !wantToUpdate() {
		return saveLastActivity()
	}

	fileURL, err := prepareURL()
	if err != nil {
		fmt.Println("sistema desconhecido")
		return err
	}

	filePath, _ := Which(azioncli)

	err = downloadFile(filePath, fileURL)
	if err != nil {
		fmt.Println("Erro ao baixar o arquivo:", err)
		return err
	}

	fmt.Println("Download concluído!")

	err = replaceFile(filePath)
	if err != nil {
		fmt.Println("Erro ao substituir o arquivo:", err)
		return err
	}

	fmt.Println("Arquivo substituído com sucesso!")
	return nil
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
		fmt.Println("error ao escrever yaml config", err)
		return err
	}

	return nil
}

func openConfig() (map[string]interface{}, error) {
	path, err := pathYamlConfig()
	if err != nil {
		fmt.Println("Erro ao obter o diretório home do usuário:", err)
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Erro ao abrir ou criar o arquivo YAML:", err)
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo YAML:", err)
		return nil, err
	}

	data := make(map[string]interface{})
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		fmt.Println("Erro ao fazer o Unmarshal do arquivo YAML:", err)
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
	data, err := openConfig()
	if err != nil {
		return "", err
	}

	lastAct := data[LastActivity]
	return fmt.Sprint(lastAct), nil
}

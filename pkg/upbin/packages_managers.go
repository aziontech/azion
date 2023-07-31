package upbin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/aziontech/azion-cli/utils"
)

const (
	urlDownloadPackage string = "https://github.com/aziontech/azion-cli/releases/download/%s/%s"
	urlAssetsAzioncli  string = "https://api.github.com/repos/aziontech/azion-cli/releases/latest"
	tapAzioncli        string = "aziontech/tap/azioncli"
)

type Release struct {
	Assets []struct {
		Name string `json:"name"`
	} `json:"assets"`
}

var (
	managersPackagesFunc          = managersPackages
	packageManagerExistsFunc      = packageManagerExists
	installPackageManagerFunc     = installPackageManager
	formatUrlPackageAzioncliFunc  = formatUrlPackageAzioncli
	downloadAndInstallPackageFunc = downloadAndInstallPackage
)

func managersPackages() (bool, error) {
	url, err := formatUrlPackageAzioncliFunc()
	if err != nil {
		return false, err
	}

	if packageManagerExistsFunc("brew") {
		err = installPackageManagerFunc("brew")
		if err != nil {
			return false, err
		}
		return true, nil
	}

	err = downloadAndInstallPackageFunc(url)
	if err != nil {
		return false, err
	}
	return true, nil
}

func formatUrlPackageAzioncli() (string, error) {
	pack := checkingPackages()
	if len(pack) == 0 {
		return "", utils.ErrorCommandNotFound
	}

	listPacks, err := getAssetsNamesPackagesAzioncli()
	if err != nil {
		return "", utils.ErrorGetAssetsNamesAzioncli
	}

	sysAr := GetInfoSystem()
	itemPack := findItemByArchAndPackage(listPacks, sysAr.Arch, pack)
	versionPack := extractVersionFromUrl(itemPack)

	url := fmt.Sprintf(urlDownloadPackage, versionPack, itemPack)

	return url, nil
}

func findItemByArchAndPackage(items []string, arch, packages string) string {
	for _, item := range items {
		if strings.Contains(item, arch) && strings.Contains(item, packages) {
			return item
		}
	}
	return ""
}

func checkingPackages() string {
	var listPackages = []string{"brew", "dpkg", "rpm", "apk"}

	var packMan string = ""

	for _, v := range listPackages {
		if packageManagerExistsFunc(v) {
			if v == "dpkg" {
				return "deb"
			}
			packMan = v
		}
	}
	return packMan
}

func extractVersionFromUrl(str string) string {
	re := regexp.MustCompile(`azioncli_([0-9]+\.[0-9]+\.[0-9]+)`)
	match := re.FindStringSubmatch(str)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func downloadAndInstallPackage(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileName := getFileNameFromURL(url)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return installPackage(fileName)
}

func getFileNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func installPackage(fileName string) error {
	ext := getFileExtension(fileName)
	switch ext {
	case "deb":
		return exec.Command("dpkg", "-i", fileName).Run()
	case "rpm":
		return exec.Command("rpm", "-i", fileName).Run()
	case "apk":
		return exec.Command("apk", "add", fileName).Run()
	default:
		return fmt.Errorf("unsupported package %s", ext)
	}
}

func getFileExtension(fileName string) string {
	parts := strings.Split(fileName, ".")
	return parts[len(parts)-1]
}

// getAssetsNamesPackagesAzioncli retrieves the names of the assets/packages. It returns a slice of strings containing the package names and an error if any occurred.
func getAssetsNamesPackagesAzioncli() ([]string, error) {
	response, err := http.Get(urlAssetsAzioncli)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		return nil, err
	}

	packageNames := make([]string, len(release.Assets))
	for i, asset := range release.Assets {
		packageNames[i] = asset.Name
	}

	return packageNames, nil
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
	_, _, err := utils.RunCommandWithOutput([]string{}, "brew upgrade azioncli")
	if err != nil {
		return err
	}
	return nil
}

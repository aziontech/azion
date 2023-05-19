package insert

import (
	"github.com/manifoldco/promptui"
)

func Insert() (string, error) {
	prompt := promptui.Prompt{
		Label: "Project name",
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, err
}

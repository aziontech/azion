package choose

import (
	"github.com/manifoldco/promptui"
)

func Choose() (string, error) {
	prompt := promptui.Select{
		Label: "Select Type",
		Items: []string{"nextjs", "static", "cdn"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

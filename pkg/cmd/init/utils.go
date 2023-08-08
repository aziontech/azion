package init

import (
	"encoding/json"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"go.uber.org/zap"
)

func yesNoFlagToResponse(info *InitInfo) bool {
	if info.YesOption {
		return info.YesOption
	}

	return false
}

func (cmd *InitCmd) createTemplateAzion(info *InitInfo) error {

	err := cmd.Mkdir(info.PathWorkingDir+"/azion", 0755) // 0755 is the permission mode for the new directories
	if err != nil {
		return msg.ErrorFailedCreatingAzionDirectory
	}

	azionJson := &contracts.AzionApplicationOptions{
		Name:      info.Name,
		Env:       "production",
		Type:      info.Template,
		VersionID: "",
	}
	azionJson.Function.Name = "__DEFAULT__"
	azionJson.Function.File = "./out/worker.js"
	azionJson.Function.Args = "./azion/args.json"
	azionJson.Domain.Name = "__DEFAULT__"
	azionJson.Application.Name = "__DEFAULT__"
	azionJson.Origin.Name = "__DEFAULT__"
	azionJson.RtPurge.PurgeOnPublish = true

	return cmd.createJsonFile(azionJson, info)

}

func (cmd *InitCmd) fetchTemplates(info *InitInfo) error {
	//create temporary directory to clone template into
	dir, err := cmd.CreateTempDir(info.PathWorkingDir, ".template")
	if err != nil {
		logger.Debug("Error while creating temporary directory for clining template", zap.Error(err))
		return utils.ErrorInternalServerError
	}
	defer func() {
		_ = cmd.RemoveAll(dir)
	}()

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: REPO})
	if err != nil {
		logger.Debug("Error while fetching templates from github", zap.Error(err))
		return utils.ErrorFetchingTemplates
	}

	tags, err := r.Tags()
	if err != nil {
		logger.Debug("Error while getting github tags", zap.Error(err))
		return msg.ErrorGetAllTags
	}

	tag, err := sortTag(tags, TemplateMajor)
	if err != nil {
		logger.Debug("Error while sorting tags for correct template application", zap.Error(err))
		return msg.ErrorIterateAllTags
	}

	_, err = cmd.GitPlainClone(dir, false, &git.CloneOptions{
		URL:           REPO,
		ReferenceName: plumbing.ReferenceName(tag),
	})
	if err != nil {
		logger.Debug("Error while fetching templates from github", zap.Error(err))
		return utils.ErrorFetchingTemplates
	}

	azionDir := info.PathWorkingDir + "/azion"

	// changing to Vulcan in case we are using any other type... this will be removed once Vulcan becomes the
	// only adapter used by the cli
	typeLang := info.Template
	if typeLang != "nextjs" && typeLang != "static" && typeLang != "simple" {
		typeLang = "vulcan"
	}

	//move contents from temporary directory into final destination
	err = cmd.Rename(dir+"/webdev/"+typeLang, azionDir)
	if err != nil {
		logger.Debug("Error while copying files to current project directory", zap.Error(err))
		return utils.ErrorMovingFiles
	}

	return nil
}

func (cmd *InitCmd) createJsonFile(options *contracts.AzionApplicationOptions, info *InitInfo) error {
	data, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return msg.ErrorUnmarshalAzionFile
	}

	err = cmd.WriteFile(info.PathWorkingDir+"/azion/azion.json", data, 0644)
	if err != nil {
		return utils.ErrorInternalServerError
	}
	return nil
}

func shouldFetch(cmd *InitCmd, info *InitInfo) (bool, error) {
	var err error
	var shouldFetchTemplates bool
	if empty, _ := cmd.IsDirEmpty("./azion"); !empty {
		if info.NoOption || info.YesOption {
			shouldFetchTemplates = yesNoFlagToResponse(info)
		} else {
			prompt := &survey.Confirm{
				Message: "This project was already configured. Do you want to override the previous configuration?",
			}
			survey.AskOne(prompt, &shouldFetchTemplates)
		}

		if shouldFetchTemplates {
			err = cmd.CleanDir("./azion")
			if err != nil {
				logger.Debug("Error while trying to clean azion directory", zap.Error(err))
				return false, err
			}
		}
		return shouldFetchTemplates, nil
	}
	return true, nil
}

func askForInput(msg string, defaultIn string) (string, error) {
	var userInput string
	prompt := &survey.Input{
		Message: msg,
		Default: defaultIn,
	}

	// Prompt the user for input
	err := survey.AskOne(prompt, &userInput, survey.WithKeepFilter(true))
	if err != nil {
		return "", err
	}
	return userInput, nil
}

func (cmd *InitCmd) selectVulcanTemplates(info *InitInfo) error {
	logger.FInfo(cmd.Io.Out, msg.InitGettingTemplates)
	output, _, err := cmd.CommandRunner("npx --yes edge-functions@1.0.0 presets ls", []string{"CLEAN_OUTPUT_MODE=true"})
	if err != nil {
		return err
	}

	newLineSplit := strings.Split(output, "\n")
	newLineSplit[len(newLineSplit)-1] = "nextjs (faststore)"

	answer := ""
	template := ""
	mode := ""
	prompt := &survey.Select{
		Message: "Choose a template:",
		Options: newLineSplit,
	}
	survey.AskOne(prompt, &answer)

	modeSplit := strings.Split(answer, " ")
	template = modeSplit[0]
	mode = strings.Replace(strings.Replace(modeSplit[1], "(", "", -1), ")", "", -1)

	info.Template = template
	info.Mode = mode

	return nil

}

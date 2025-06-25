package link

import (
	"encoding/json"
	"path"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
)

func (cmd *LinkCmd) createTemplateAzion(info *LinkInfo) error {

	err := cmd.Mkdir(path.Join(info.PathWorkingDir, info.projectPath), 0755) // 0755 is the permission mode for the new directories
	if err != nil {
		return msg.ErrorFailedCreatingAzionDirectory
	}

	azionJson := &contracts.AzionApplicationOptions{
		Name:   info.Name,
		Env:    "production",
		Preset: strings.ToLower(info.Preset),
		Prefix: "",
	}
	azionJson.Function = []contracts.AzionJsonDataFunction{}
	azionJson.Domain.Name = "__DEFAULT__"
	azionJson.Application.Name = "__DEFAULT__"
	azionJson.RtPurge.PurgeOnPublish = true

	return cmd.createJsonFile(azionJson, info)

}

func (cmd *LinkCmd) createJsonFile(options *contracts.AzionApplicationOptions, info *LinkInfo) error {
	data, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return msg.ErrorUnmarshalAzionFile
	}

	err = cmd.WriteFile(path.Join(info.PathWorkingDir, info.projectPath, "azion.json"), data, 0644)
	if err != nil {
		return utils.ErrorInternalServerError
	}
	return nil
}

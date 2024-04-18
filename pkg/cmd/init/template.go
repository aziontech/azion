package init

import (
	"encoding/json"

	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
)

func (cmd *initCmd) createTemplateAzion() error {

	err := cmd.mkdir(cmd.pathWorkingDir+"/azion", 0755) // 0755 is the permission mode for the new directories
	if err != nil {
		return msg.ErrorFailedCreatingAzionDirectory
	}

	azionJson := &contracts.AzionApplicationOptions{
		Name:   cmd.name,
		Env:    "production",
		Preset: cmd.preset,
		Mode:   cmd.mode,
		Prefix: "",
	}

	azionJson.Function.Name = "__DEFAULT__"
	azionJson.Function.InstanceName = "__DEFAULT__"
	azionJson.Function.File = "./out/worker.js"
	azionJson.Function.Args = "./azion/args.json"
	azionJson.Domain.Name = "__DEFAULT__"
	azionJson.Application.Name = "__DEFAULT__"
	azionJson.Origin.Name = "__DEFAULT__"
	azionJson.RtPurge.PurgeOnPublish = true

	return cmd.createJsonFile(azionJson)

}

func (cmd *initCmd) createJsonFile(options *contracts.AzionApplicationOptions) error {
	data, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return msg.ErrorUnmarshalAzionFile
	}

	err = cmd.writeFile(cmd.pathWorkingDir+"/azion/azion.json", data, 0644)
	if err != nil {
		return utils.ErrorInternalServerError
	}
	return nil
}

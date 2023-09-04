package link

import (
	"encoding/json"
	"fmt"
	"os"

	msg "github.com/aziontech/azion-cli/messages/link"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func initSimple(cmd *LinkCmd, path string, info *LinkInfo, c *cobra.Command) error {
	var err error
	var shouldFetchTemplates bool
	options := &contracts.AzionApplicationSimple{}

	shouldFetchTemplates, err = shouldFetch(cmd, info)
	if err != nil {
		return err
	}

	if shouldFetchTemplates {
		pathWorker := path + "/azion"
		if err := cmd.Mkdir(pathWorker, os.ModePerm); err != nil {
			logger.Debug("Error while creating azion directory", zap.Error(err))
			return msg.ErrorFailedCreatingAzionDirectory
		}

		options.Name = info.Name
		options.Type = info.Preset
		options.Domain.Name = "__DEFAULT__"
		options.Application.Name = "__DEFAULT__"

		data, err := json.MarshalIndent(options, "", "  ")
		if err != nil {
			logger.Debug("Error while marshling json file", zap.Error(err))
			return msg.ErrorUnmarshalAzionFile
		}

		err = cmd.WriteFile(path+"/azion/azion.json", data, 0644)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return utils.ErrorInternalServerError
		}

		logger.FInfo(cmd.Io.Out, fmt.Sprintf(msg.EdgeApplicationsLinkSuccessful+"\n", info.Name))
	}

	return nil
}

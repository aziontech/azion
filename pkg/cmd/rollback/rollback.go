package rollback

import (
	"context"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/rollback"
	apiOrigin "github.com/aziontech/azion-cli/pkg/api/origin"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	originKey   string
	projectPath string
)

type RollbackCmd struct {
	AskInput              func(string) (string, error)
	GetAzionJsonContent   func(pathConf string) (*contracts.AzionApplicationOptions, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confPath string) error
}

func NewDeleteCmd(f *cmdutil.Factory) *RollbackCmd {
	return &RollbackCmd{
		GetAzionJsonContent:   utils.GetAzionJsonContent,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
		AskInput:              utils.AskInput,
	}
}

func NewCobraCmd(rollback *RollbackCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.USAGE,
		Short:         msg.SHORTDESCRIPTION,
		Long:          msg.LONGDESCRIPTION,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion rollback --origin-key aaaa-bbbb-cccc-dddd
		`),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("origin-key") {
				answer, err := rollback.AskInput(msg.ASKORIGIN)
				if err != nil {
					return err
				}
				originKey = answer
			}

			conf, err := rollback.GetAzionJsonContent(projectPath)
			if err != nil {
				logger.Debug("Error while building your project", zap.Error(err))
				return msg.ERRORAZION
			}

			if conf.Bucket == "" || conf.Prefix == "" {
				return msg.ERRORNEEDSDEPLOY
			}

			timestamp, err := checkForNewTimestamp(f, conf.Prefix, conf.Bucket)
			if err != nil {
				return msg.ERRORROLLBACK
			}

			clientOrigin := apiOrigin.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			request := apiOrigin.UpdateRequest{}
			request.SetPrefix(timestamp)

			_, err = clientOrigin.Update(context.Background(), conf.Application.ID, originKey, &request)
			if err != nil {
				return msg.ERRORROLLBACK
			}

			conf.Prefix = timestamp
			err = rollback.WriteAzionJsonContent(conf, projectPath)
			if err != nil {
				return msg.ERRORROLLBACK
			}

			return nil
		},
	}

	cobraCmd.Flags().StringVar(&originKey, "origin-key", "", msg.FLAGORIGINKEY)
	cobraCmd.Flags().StringVar(&projectPath, "config-dir", "azion", msg.CONFFLAG)
	cobraCmd.Flags().BoolP("help", "h", false, msg.FLAGHELP)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}

func checkForNewTimestamp(f *cmdutil.Factory, referenceTimestamp, bucketName string) (string, error) {
	logger.Debug("Checking if there are previous static files for the following bucket", zap.Any("Bucket name", bucketName))
	client := api.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))
	c := context.Background()
	options := &contracts.ListOptions{
		// Sort:     "desc",
		// OrderBy:  "last_modified",
		PageSize: 100000,
	}

	resp, err := client.ListObject(c, bucketName, options)
	if err != nil {
		return "", err
	}

	var prevTimestamp string
	for _, object := range resp.Results {
		parts := strings.Split(object.Key, "/")
		if len(parts) > 1 {
			timestamp := parts[0]
			if timestamp == referenceTimestamp {
				return prevTimestamp, nil
			} else {
				prevTimestamp = timestamp
				continue
			}
			// if timestamp != referenceTimestamp {
			// 	return timestamp, nil
			// }
		}
	}

	return referenceTimestamp, nil
}

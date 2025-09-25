package edgeapplication

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/application"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	app "github.com/aziontech/azion-cli/pkg/v3api/edge_applications"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
)

var ProjectConf string

type DeleteCmd struct {
	Io         *iostreams.IOStreams
	GetAzion   func(confPath string) (*contracts.AzionApplicationOptionsV3, error)
	f          *cmdutil.Factory
	UpdateJson func(cmd *DeleteCmd) error
	Cascade    func(ctx context.Context, del *DeleteCmd) error
	AskInput   func(string) (string, error)
	ReadFile   func(name string) ([]byte, error)
	WriteFile  func(name string, data []byte, perm fs.FileMode) error
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f))
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io:         f.IOStreams,
		GetAzion:   utils.GetAzionJsonContentV3,
		f:          f,
		UpdateJson: updateAzionJson,
		Cascade:    CascadeDelete,
		AskInput:   utils.AskInput,
		ReadFile:   os.ReadFile,
		WriteFile:  os.WriteFile,
	}
}

func NewCobraCmd(delete *DeleteCmd) *cobra.Command {
	var application_id int64
	cmd := &cobra.Command{
		Use:           "edge-application",
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion delete edge-application --application-id 1234
		$ azion delete edge-application --cascade
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete.run(cmd, application_id)
		},
	}

	cmd.Flags().Int64Var(&application_id, "application-id", 0, msg.FlagId)
	cmd.Flags().Bool("cascade", true, msg.CascadeFlag)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)
	cmd.Flags().StringVar(&ProjectConf, "config-dir", "azion", msg.CONFDIRFLAG)

	return cmd
}

func (del *DeleteCmd) run(cmd *cobra.Command, application_id int64) error {
	ctx := context.Background()

	if cmd.Flags().Changed("cascade") {
		err := del.Cascade(ctx, del)
		if err != nil {
			return err
		}
		return nil
	}

	if !cmd.Flags().Changed("application-id") {
		answer, err := del.AskInput(msg.AskInput)
		if err != nil {
			return err
		}

		num, err := strconv.ParseInt(answer, 10, 64)
		if err != nil {
			logger.Debug("Error while converting answer to int64", zap.Error(err))
			return msg.ErrorConvertId
		}

		application_id = num
	}

	client := app.NewClient(del.f.HttpClient, del.f.Config.GetString("api_url"), del.f.Config.GetString("token"))

	err := client.Delete(ctx, application_id)
	if err != nil {
		return fmt.Errorf(msg.ErrorFailToDeleteApplication.Error(), err)
	}

	deleteOut := output.GeneralOutput{
		Msg:   fmt.Sprintf(msg.OutputSuccess, application_id),
		Out:   del.f.IOStreams.Out,
		Flags: del.f.Flags,
	}
	return output.Print(&deleteOut)
}

func updateAzionJson(cmd *DeleteCmd) error {
	path, err := utils.GetWorkingDir()
	if err != nil {
		return utils.ErrorInternalServerError
	}
	azionJson := path + "/azion/azion.json"
	byteAzionJson, err := cmd.ReadFile(azionJson)
	if err != nil {
		logger.Debug("Error while parsing json", zap.Error(err))
		return utils.ErrorUnmarshalAzionJsonFile
	}
	jsonReplaceFunc, err := sjson.Set(string(byteAzionJson), "function.id", 0)
	if err != nil {
		return msg.ErrorFailedUpdateAzionJson
	}

	jsonReplaceApp, err := sjson.Set(jsonReplaceFunc, "application.id", 0)
	if err != nil {
		return msg.ErrorFailedUpdateAzionJson
	}

	jsonReplaceDomain, err := sjson.Set(jsonReplaceApp, "domain.id", 0)
	if err != nil {
		return msg.ErrorFailedUpdateAzionJson
	}

	jsonReplaceFuncInstanceID, err := sjson.Set(jsonReplaceDomain, "function.instance-id", 0)
	if err != nil {
		return msg.ErrorFailedUpdateAzionJson
	}

	jsonReplaceFirstRun, err := sjson.Set(jsonReplaceFuncInstanceID, "not-first-run", false)
	if err != nil {
		return msg.ErrorFailedUpdateAzionJson
	}

	err = cmd.WriteFile(azionJson, []byte(jsonReplaceFirstRun), 0644)
	if err != nil {
		return fmt.Errorf(utils.ErrorCreateFile.Error(), azionJson)
	}

	return nil
}

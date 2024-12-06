package edgeapplication

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/edge_application"
	app "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var ProjectConf string

type DeleteCmd struct {
	Io                    *iostreams.IOStreams
	GetAzion              func(confPath string) (*contracts.AzionApplicationOptions, error)
	f                     *cmdutil.Factory
	UpdateJson            func(cmd *DeleteCmd) error
	Cascade               func(ctx context.Context, del *DeleteCmd) error
	AskInput              func(string) (string, error)
	ReadFile              func(name string) ([]byte, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confPath string) error
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f))
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io:                    f.IOStreams,
		GetAzion:              utils.GetAzionJsonContent,
		f:                     f,
		UpdateJson:            updateAzionJson,
		Cascade:               CascadeDelete,
		AskInput:              utils.AskInput,
		ReadFile:              os.ReadFile,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
	}
}

func NewCobraCmd(delete *DeleteCmd) *cobra.Command {
	var application_id int64
	cmd := &cobra.Command{
		Use:           msg.Usage,
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
	paths, err := utils.GetWorkingDir()
	if err != nil {
		return utils.ErrorInternalServerError
	}
	azionJson := path.Join(paths, "azion/azion.json")

	azionJsonFile := &contracts.AzionApplicationOptions{
		Env:    "production",
		Prefix: "",
	}
	azionJsonFile.Function.Name = "__DEFAULT__"
	azionJsonFile.Function.InstanceName = "__DEFAULT__"
	azionJsonFile.Function.File = "./out/worker.js"
	azionJsonFile.Domain.Name = "__DEFAULT__"
	azionJsonFile.Application.Name = "__DEFAULT__"
	azionJsonFile.RtPurge.PurgeOnPublish = true

	err = cmd.WriteAzionJsonContent(azionJsonFile, ProjectConf)
	if err != nil {
		return fmt.Errorf(utils.ErrorCreateFile.Error(), azionJson)
	}

	return nil
}

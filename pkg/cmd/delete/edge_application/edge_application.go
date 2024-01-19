package edgeapplication

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/edge_application"
	app "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	fun "github.com/aziontech/azion-cli/pkg/api/edge_function"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
)

type DeleteCmd struct {
	Io         *iostreams.IOStreams
	GetAzion   func() (*contracts.AzionApplicationOptions, error)
	f          *cmdutil.Factory
	UpdateJson func(cmd *DeleteCmd) error
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f))
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io:         f.IOStreams,
		GetAzion:   utils.GetAzionJsonContent,
		f:          f,
		UpdateJson: updateAzionJson,
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

	return cmd
}

func (del *DeleteCmd) run(cmd *cobra.Command, application_id int64) error {
	ctx := context.Background()

	if cmd.Flags().Changed("cascade") {
		azionJson, err := del.GetAzion()
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return msg.ErrorMissingAzionJson
			} else {
				return err //message with error
			}
		}
		if azionJson.Application.ID == 0 {
			return msg.ErrorMissingApplicationIdJson
		}
		clientapp := app.NewClient(del.f.HttpClient, del.f.Config.GetString("api_url"), del.f.Config.GetString("token"))
		clientfunc := fun.NewClient(del.f.HttpClient, del.f.Config.GetString("api_url"), del.f.Config.GetString("token"))

		err = clientapp.Delete(ctx, azionJson.Application.ID)
		if err != nil {
			return fmt.Errorf(msg.ErrorFailToDeleteApplication.Error(), err)
		}

		if azionJson.Function.ID == 0 {
			fmt.Fprintf(del.f.IOStreams.Out, msg.MissingFunction)
		} else {
			err = clientfunc.Delete(ctx, azionJson.Function.ID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeleteApplication.Error(), err)
			}
		}

		fmt.Fprintf(del.f.IOStreams.Out, "%s\n", msg.CascadeSuccess)

		err = del.UpdateJson(del)
		if err != nil {
			return err
		}

		return nil

	}

	if !cmd.Flags().Changed("application-id") {
		qs := []*survey.Question{
			{
				Name:     "id",
				Prompt:   &survey.Input{Message: msg.AskInput},
				Validate: survey.Required,
			},
		}

		answer := ""

		err := survey.Ask(qs, &answer)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
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

	out := del.f.IOStreams.Out
	fmt.Fprintf(out, msg.OutputSuccess, application_id)

	return nil
}

func updateAzionJson(cmd *DeleteCmd) error {
	path, err := utils.GetWorkingDir()
	if err != nil {
		return utils.ErrorInternalServerError
	}
	azionJson := path + "/azion/azion.json"
	byteAzionJson, err := os.ReadFile(azionJson)
	if err != nil {
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

	err = os.WriteFile(azionJson, []byte(jsonReplaceDomain), 0644)
	if err != nil {
		return fmt.Errorf(utils.ErrorCreateFile.Error(), azionJson)
	}

	return nil
}

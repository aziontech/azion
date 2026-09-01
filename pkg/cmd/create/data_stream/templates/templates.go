package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/templates"
	api "github.com/aziontech/azion-cli/pkg/api/data_stream"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type Fields struct {
	Name    string
	DataSet string
	Active  string
	InPath  string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create data-stream templates --name "My Template" --data-set "./data_set.json"
        $ azion create data-stream templates --name "My Template" --data-set "./data_set.json" --active false
        $ azion create data-stream templates --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.NewCreateTemplateRequest()

			if cmd.Flags().Changed("file") {
				if err := utils.FlagFileUnmarshalJSON(fields.InPath, &request.TemplateRequest); err != nil {
					logger.Debug("Failed to unmarshal file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("name") {
					answer, err := utils.AskInput(msg.AskInputName)
					if err != nil {
						return err
					}
					fields.Name = answer
				}

				if !cmd.Flags().Changed("data-set") {
					answer, err := utils.AskInput(msg.AskInputDataSet)
					if err != nil {
						return err
					}
					fields.DataSet = answer
				}

				dataSet, err := os.ReadFile(fields.DataSet)
				if err != nil {
					logger.Debug("Failed to read the data set file", zap.Error(err))
					return msg.ErrorDataSetFlag
				}
				if !json.Valid(dataSet) {
					return msg.ErrorParseDataSet
				}

				request.TemplateRequest = *sdk.NewTemplateRequest(fields.Name, string(dataSet))

				if cmd.Flags().Changed("active") {
					isActive, err := strconv.ParseBool(fields.Active)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorActiveFlag, fields.Active)
					}
					request.TemplateRequest.SetActive(isActive)
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.CreateTemplate(ctx, request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateTemplate.Error(), err)
			}

			creatOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.CreateOutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&creatOut)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)

	return cmd
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.DataSet, "data-set", "", msg.FlagDataSet)
	flags.StringVar(&fields.Active, "active", "true", msg.FlagActive)
	flags.StringVar(&fields.InPath, "file", "", msg.FlagIn)
	flags.BoolP("help", "h", false, msg.CreateFlagHelp)
}

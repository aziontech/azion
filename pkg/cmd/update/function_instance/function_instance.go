package functioninstance

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/update/function_instance"
	api "github.com/aziontech/azion-cli/pkg/api/function_instance"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Name          string `json:"name"`
	Active        string `json:"active"`
	Args          string `json:"args"`
	Path          string
	ApplicationID string
	FunctionID    int64
	InstanceID    string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion update function-instance --name functionInstanceName
        $ azion update function-instance --name withargs --active true
        $ azion update function-instance --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := sdk.PatchedApplicationFunctionInstanceRequest{}
			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("application-id") {
					answer, err := utils.AskInput(msg.AskInputApplicationID)
					if err != nil {
						return err
					}

					fields.ApplicationID = answer
				}

				if !cmd.Flags().Changed("instance-id") {
					answer, err := utils.AskInput(msg.AskInputInstanceID)
					if err != nil {
						return err
					}

					fields.InstanceID = answer
				}

				if cmd.Flags().Changed("function-id") {
					request.SetFunction(fields.FunctionID)
				}

				if cmd.Flags().Changed("name") {
					request.SetName(fields.Name)
				}

				isActive, err := strconv.ParseBool(fields.Active)
				if err != nil {
					return fmt.Errorf("%w: %q", msg.ErrorIsActiveFlag, fields.Active)
				}
				request.SetActive(isActive)

				if cmd.Flags().Changed("args") {
					marshalledArgs, err := os.ReadFile(fields.Args)
					if err != nil {
						return fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
					}

					args := make(map[string]interface{})
					if err := json.Unmarshal(marshalledArgs, &args); err != nil {
						return fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
					}
					request.SetArgs(args)

				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			response, err := client.Update(context.Background(), fields.ApplicationID, fields.InstanceID, request)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdate.Error(), err)
			}

			createOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&createOut)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.Active, "active", "true", msg.FlagIsActive)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.StringVar(&fields.Args, "args", "", msg.FlagArgs)
	flags.StringVar(&fields.ApplicationID, "application-id", "", msg.FlagApplicationID)
	flags.Int64Var(&fields.FunctionID, "function-id", 0, msg.FlagFunctionID)
	flags.StringVar(&fields.InstanceID, "instance-id", "", msg.FlagInstanceID)
	flags.BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

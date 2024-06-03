package edgefunction

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_function"
	api "github.com/aziontech/azion-cli/pkg/api/edge_function"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type Fields struct {
	ID            int64
	Name          string
	Language      string
	Code          string
	Active        string
	InitiatorType string
	Args          string
	InPath        string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.UpdateShortDescription,
		Long:          msg.UpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion update edge-function --function-id 1234 --name 'Hello'
		$ azion update edge-function -f 4185 --code ./mycode/function.js --args ./mycode/myargs.json
		$ azion update edge-function -f 9123 --active true
		$ azion update edge-function -f 9123 --active false
		$ azion update edge-function --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			// either function-id or in path should be passed
			if !cmd.Flags().Changed("function-id") {
				answers, err := utils.AskInput(msg.UpdateAskEdgeFunctionID)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				id, err := strconv.Atoi(answers)
				if err != nil {
					logger.Debug("Error while parsing string to integer", zap.Error(err))
					return utils.ErrorConvertingStringToInt
				}

				fields.ID = int64(id)
			}

			request := api.UpdateRequest{}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {
				err := createRequestFromFlags(cmd, fields, &request)
				if err != nil {
					return err
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, &request, fields.ID)

			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateFunction.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.UpdateOutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&updateOut)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)

	return cmd
}

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.UpdateRequest) error {
	if cmd.Flags().Changed("active") {
		active, err := strconv.ParseBool(fields.Active)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorActiveFlag, fields.Active)
		}
		request.SetActive(active)
	}

	if cmd.Flags().Changed("code") {
		code, err := os.ReadFile(fields.Code)
		if err != nil {
			return fmt.Errorf("%s: %w", msg.ErrorCodeFlag, err)
		}
		request.SetCode(string(code))
	}

	if cmd.Flags().Changed("args") {
		marshalledArgs, err := os.ReadFile(fields.Args)
		if err != nil {
			return fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
		}
		args := make(map[string]interface{})
		if err := json.Unmarshal(marshalledArgs, &args); err != nil {
			return fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
		}
		request.SetJsonArgs(args)
	}

	if cmd.Flags().Changed("name") {
		request.SetName(fields.Name)
	}

	return nil
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.Int64Var(&fields.ID, "function-id", 0, msg.FlagID)
	flags.StringVar(&fields.Name, "name", "", msg.UpdateFlagName)
	flags.StringVar(&fields.Code, "code", "", msg.UpdateFlagCode)
	flags.StringVar(&fields.Args, "args", "", msg.UpdateFlagArgs)
	flags.StringVar(&fields.Active, "active", "", msg.UpdateFlagActive)
	flags.StringVar(&fields.InPath, "file", "", msg.UpdateFlagFile)
	flags.BoolP("help", "h", false, msg.UpdateHelpFlag)
}

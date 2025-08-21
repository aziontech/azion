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
	Name                 string
	Language             string
	Code                 string
	Active               string
	Args                 string
	ExecutionEnvironment string
	InPath               string
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
        $ azion create edge-function --name misfunction --code ./code/function.js --active false
        $ azion create edge-function --name with args --code ./code/function.js --args ./args.json --active true
        $ azion create edge-function --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.NewCreateRequest()

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &request)
				if err != nil {
					logger.Debug("Failed to unmarshal file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				err := createRequestFromFlags(cmd, fields, request)
				if err != nil {
					return err
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Create(ctx, request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateFunction.Error(), err)
			}

			creatOut := output.GeneralOutput{
				Msg: fmt.Sprintf(msg.CreateOutputSuccess, response.GetId()),
				Out: f.IOStreams.Out,
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
	flags.StringVar(&fields.Code, "code", "", msg.FlagCode)
	flags.StringVar(&fields.Active, "active", "", msg.FlagActive)
	flags.StringVar(&fields.Args, "args", "", msg.FlagArgs)
	flags.StringVar(&fields.ExecutionEnvironment, "execution-environment", "", msg.FlagExecutionEnvironment)
	flags.StringVar(&fields.InPath, "file", "", msg.FlagIn)
	flags.BoolP("help", "h", false, msg.CreateFlagHelp)
}

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.CreateRequest) error {

	if !cmd.Flags().Changed("name") {
		answers, err := utils.AskInput(msg.AskName)

		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}

		fields.Name = answers
	}

	if !cmd.Flags().Changed("code") {
		answers, err := utils.AskInput(msg.AskCode)

		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}

		fields.Code = answers
	}

	if !cmd.Flags().Changed("active") {
		answers, err := utils.Select(utils.NewSelectPrompter(msg.AskActive, []string{"true", "false"}))
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}

		fields.Active = answers
	}

	isActive, err := strconv.ParseBool(fields.Active)
	if err != nil {
		return fmt.Errorf("%w: %s", msg.ErrorActiveFlag, fields.Active)
	}
	request.SetActive(isActive)

	code, err := os.ReadFile(fields.Code)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.ErrorCodeFlag, err)
	}
	request.SetCode(string(code))

	if cmd.Flags().Changed("args") {
		marshalledArgs, err := os.ReadFile(fields.Args)
		if err != nil {
			return fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
		}

		args := make(map[string]interface{})
		if err := json.Unmarshal(marshalledArgs, &args); err != nil {
			return fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
		}
		request.SetDefaultArgs(args)
	}

	if cmd.Flags().Changed("execution-environment") {
		request.SetExecutionEnvironment(fields.ExecutionEnvironment)
	}

	request.SetName(fields.Name)

	return nil
}

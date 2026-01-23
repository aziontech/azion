package firewallinstance

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/create/firewall_instance"
	api "github.com/aziontech/azion-cli/pkg/api/firewall_instance"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Name       string `json:"name"`
	Active     string `json:"active"`
	Args       string `json:"args"`
	Path       string
	FirewallID int64
	FunctionID int64
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
        $ azion create firewall-instance --name firewallInstanceName --firewall-id 1234
        $ azion create firewall-instance --name withargs --active true --firewall-id 1234
        $ azion create firewall-instance --file "create.json" --firewall-id 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.NewCreateRequest()
			if !cmd.Flags().Changed("firewall-id") {
				answer, err := utils.AskInput(msg.AskInputFirewallID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertFirewallId
				}

				fields.FirewallID = num
			}
			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {

				if !cmd.Flags().Changed("function-id") {
					answer, err := utils.AskInput(msg.AskInputFunctionID)
					if err != nil {
						return err
					}

					num, err := strconv.ParseInt(answer, 10, 64)
					if err != nil {
						logger.Debug("Error while converting answer to int64", zap.Error(err))
						return msg.ErrorConvertFunctionID
					}

					fields.FunctionID = num

				}
				request.SetFunction(fields.FunctionID)

				if !cmd.Flags().Changed("name") {
					answer, err := utils.AskInput(msg.AskInputName)
					if err != nil {
						return err
					}

					fields.Name = answer
				}

				request.SetName(fields.Name)

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
			response, err := client.Create(context.Background(), fields.FirewallID, request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreate.Error(), err)
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
	flags.Int64Var(&fields.FirewallID, "firewall-id", 0, msg.FlagFirewallID)
	flags.Int64Var(&fields.FunctionID, "function-id", 0, msg.FlagFunctionID)
	flags.BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

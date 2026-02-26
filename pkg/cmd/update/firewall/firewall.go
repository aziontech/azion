package firewall

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/update/firewall"
	api "github.com/aziontech/azion-cli/pkg/api/firewall"
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
	ID                int64
	Name              string
	DebugRules        string
	Active            string
	FunctionsEnabled  string
	NetworkProtection string
	WafEnabled        string
	InPath            string
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
		$ azion update firewall --firewall-id 1234 --name 'My Firewall'
		$ azion update firewall --firewall-id 4185 --active true --waf-enabled true
		$ azion update firewall --firewall-id 9123 --debug-rules false
		$ azion update firewall --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("firewall-id") {
				answer, err := utils.AskInput(msg.UpdateAskFirewallID)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdFirewall
				}

				fields.ID = num
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

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, &request, fields.ID)

			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateFirewall.Error(), err)
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
	if cmd.Flags().Changed("name") {
		request.SetName(fields.Name)
	}

	if cmd.Flags().Changed("debug-rules") {
		isDebug, err := strconv.ParseBool(fields.DebugRules)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorDebugFlag, fields.DebugRules)
		}
		request.SetDebug(isDebug)
	}

	if cmd.Flags().Changed("active") {
		isActive, err := strconv.ParseBool(fields.Active)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorActiveFlag, fields.Active)
		}
		request.SetActive(isActive)
	}

	if cmd.Flags().Changed("functions-enabled") || cmd.Flags().Changed("network-protection") || cmd.Flags().Changed("waf-enabled") {
		modules := sdk.NewFirewallModulesRequest()

		if cmd.Flags().Changed("functions-enabled") {
			isFunctionsEnabled, err := strconv.ParseBool(fields.FunctionsEnabled)
			if err != nil {
				return fmt.Errorf("%w: %s", msg.ErrorFunctionsEnabledFlag, fields.FunctionsEnabled)
			}
			functionsModule := sdk.NewFirewallModuleRequest()
			functionsModule.SetEnabled(isFunctionsEnabled)
			modules.SetFunctions(*functionsModule)
		}

		if cmd.Flags().Changed("network-protection") {
			isNetworkProtection, err := strconv.ParseBool(fields.NetworkProtection)
			if err != nil {
				return fmt.Errorf("%w: %s", msg.ErrorNetworkProtectionFlag, fields.NetworkProtection)
			}
			networkModule := sdk.NewFirewallModuleRequest()
			networkModule.SetEnabled(isNetworkProtection)
			modules.SetNetworkProtection(*networkModule)
		}

		if cmd.Flags().Changed("waf-enabled") {
			isWafEnabled, err := strconv.ParseBool(fields.WafEnabled)
			if err != nil {
				return fmt.Errorf("%w: %s", msg.ErrorWafEnabledFlag, fields.WafEnabled)
			}
			wafModule := sdk.NewFirewallModuleRequest()
			wafModule.SetEnabled(isWafEnabled)
			modules.SetWaf(*wafModule)
		}

		request.SetModules(*modules)
	}

	return nil
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.Int64Var(&fields.ID, "firewall-id", 0, msg.FlagID)
	flags.StringVar(&fields.Name, "name", "", msg.UpdateFlagName)
	flags.StringVar(&fields.DebugRules, "debug-rules", "", msg.UpdateFlagDebug)
	flags.StringVar(&fields.Active, "active", "", msg.UpdateFlagActive)
	flags.StringVar(&fields.FunctionsEnabled, "functions-enabled", "", msg.UpdateFlagFunctionsEnabled)
	flags.StringVar(&fields.NetworkProtection, "network-protection", "", msg.UpdateFlagNetworkProtection)
	flags.StringVar(&fields.WafEnabled, "waf-enabled", "", msg.UpdateFlagWafEnabled)
	flags.StringVar(&fields.InPath, "file", "", msg.UpdateFlagFile)
	flags.BoolP("help", "h", false, msg.UpdateHelpFlag)
}

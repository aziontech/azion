package firewall

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/create/firewall"
	api "github.com/aziontech/azion-cli/pkg/api/firewall"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type Fields struct {
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
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create firewall --name "My Firewall" --active true --functions-enabled true
        $ azion create firewall --name "WAF Firewall" --waf-enabled true --network-protection true
        $ azion create firewall --file "create.json"
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
	flags.StringVar(&fields.DebugRules, "debug-rules", "", msg.FlagDebug)
	flags.StringVar(&fields.Active, "active", "", msg.FlagActive)
	flags.StringVar(&fields.FunctionsEnabled, "functions-enabled", "", msg.FlagFunctionsEnabled)
	flags.StringVar(&fields.NetworkProtection, "network-protection", "", msg.FlagNetworkProtection)
	flags.StringVar(&fields.WafEnabled, "waf-enabled", "", msg.FlagWafEnabled)
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

	request.SetName(fields.Name)

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

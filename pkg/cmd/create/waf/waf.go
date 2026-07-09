package waf

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/create/waf"
	api "github.com/aziontech/azion-cli/pkg/api/waf"
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
	Name           string
	Active         string
	ProductVersion string
	EngineVersion  string
	Type           string
	Rulesets       string
	Thresholds     string
	InPath         string
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
		$ azion create waf --name "My WAF" --active true
		$ azion create waf --name "My WAF" --active true --engine-version "2021-Q3" --type "score"
		$ azion create waf --name "My WAF" --rulesets "1,2,3" --thresholds "cross_site_scripting=highest,sql_injection=high"
		$ azion create waf --file "create.json"
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
				return fmt.Errorf(msg.ErrorCreateWAF.Error(), err)
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
	flags.StringVar(&fields.Active, "active", "", msg.FlagActive)
	flags.StringVar(&fields.ProductVersion, "product-version", "", msg.FlagProductVersion)
	flags.StringVar(&fields.EngineVersion, "engine-version", "", msg.FlagEngineVersion)
	flags.StringVar(&fields.Type, "type", "", msg.FlagType)
	flags.StringVar(&fields.Rulesets, "rulesets", "", msg.FlagRulesets)
	flags.StringVar(&fields.Thresholds, "thresholds", "", msg.FlagThresholds)
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

	if cmd.Flags().Changed("active") {
		isActive, err := strconv.ParseBool(fields.Active)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorActiveFlag, fields.Active)
		}
		request.SetActive(isActive)
	}

	if cmd.Flags().Changed("product-version") {
		request.SetProductVersion(fields.ProductVersion)
	}

	if cmd.Flags().Changed("engine-version") || cmd.Flags().Changed("type") || cmd.Flags().Changed("rulesets") || cmd.Flags().Changed("thresholds") {
		engineSettings := sdk.NewWAFEngineSettingsFieldRequest()

		if cmd.Flags().Changed("engine-version") {
			engineSettings.SetEngineVersion(fields.EngineVersion)
		}

		if cmd.Flags().Changed("type") {
			engineSettings.SetType(fields.Type)
		}

		if cmd.Flags().Changed("rulesets") || cmd.Flags().Changed("thresholds") {
			attributes := sdk.NewWAFEngineSettingsAttributesFieldRequest()

			if cmd.Flags().Changed("rulesets") {
				rulesets, err := parseRulesets(fields.Rulesets)
				if err != nil {
					return err
				}
				attributes.SetRulesets(rulesets)
			}

			if cmd.Flags().Changed("thresholds") {
				thresholds, err := parseThresholds(fields.Thresholds)
				if err != nil {
					return err
				}
				attributes.SetThresholds(thresholds)
			}

			engineSettings.SetAttributes(*attributes)
		}

		request.SetEngineSettings(*engineSettings)
	}

	return nil
}

func parseRulesets(rulesetsStr string) ([]int64, error) {
	parts := strings.Split(rulesetsStr, ",")
	var rulesets []int64
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", msg.ErrorRulesetsFlag, part)
		}
		rulesets = append(rulesets, id)
	}
	return rulesets, nil
}

func parseThresholds(thresholdsStr string) ([]sdk.ThresholdsConfigFieldRequest, error) {
	validThreats := map[string]bool{
		"cross_site_scripting":  true,
		"directory_traversal":   true,
		"evading_tricks":        true,
		"file_upload":           true,
		"identified_attack":     true,
		"remote_file_inclusion": true,
		"sql_injection":         true,
		"unwanted_access":       true,
	}

	validSensitivities := map[string]bool{
		"highest": true,
		"high":    true,
		"medium":  true,
		"low":     true,
		"lowest":  true,
	}

	parts := strings.Split(thresholdsStr, ",")
	var thresholds []sdk.ThresholdsConfigFieldRequest
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			return nil, msg.ErrorThresholdsFlag
		}
		threat := strings.TrimSpace(kv[0])
		sensitivity := strings.TrimSpace(kv[1])

		if !validThreats[threat] {
			return nil, msg.ErrorThresholdThreat
		}
		if !validSensitivities[sensitivity] {
			return nil, msg.ErrorThresholdSensitivity
		}

		threshold := sdk.NewThresholdsConfigFieldRequest(threat)
		threshold.SetSensitivity(sensitivity)
		thresholds = append(thresholds, *threshold)
	}
	return thresholds, nil
}

package waf

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/update/waf"
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
	ID            int64
	Name          string
	Active        string
	EngineVersion string
	Type          string
	Rulesets      string
	Thresholds    string
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
		$ azion update waf --waf-id 1234 --name "My WAF"
		$ azion update waf --waf-id 4185 --active true
		$ azion update waf --waf-id 9123 --engine-version "2021-Q3" --type "score"
		$ azion update waf --waf-id 5678 --rulesets "1,2,3" --thresholds "cross_site_scripting=highest,sql_injection=high"
		$ azion update waf --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("waf-id") {
				answer, err := utils.AskInput(msg.UpdateAskWafID)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdWaf
				}

				fields.ID = num
			}

			request := api.NewUpdateRequest()

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &request)
				if err != nil {
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
			response, err := client.Update(ctx, request, fields.ID)

			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateWAF.Error(), err)
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

	if cmd.Flags().Changed("active") {
		isActive, err := strconv.ParseBool(fields.Active)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorActiveFlag, fields.Active)
		}
		request.SetActive(isActive)
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

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.Int64Var(&fields.ID, "waf-id", 0, msg.FlagID)
	flags.StringVar(&fields.Name, "name", "", msg.UpdateFlagName)
	flags.StringVar(&fields.Active, "active", "", msg.UpdateFlagActive)
	flags.StringVar(&fields.EngineVersion, "engine-version", "", msg.UpdateFlagEngineVersion)
	flags.StringVar(&fields.Type, "type", "", msg.UpdateFlagType)
	flags.StringVar(&fields.Rulesets, "rulesets", "", msg.UpdateFlagRulesets)
	flags.StringVar(&fields.Thresholds, "thresholds", "", msg.UpdateFlagThresholds)
	flags.StringVar(&fields.InPath, "file", "", msg.UpdateFlagFile)
	flags.BoolP("help", "h", false, msg.UpdateHelpFlag)
}

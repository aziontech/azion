package application

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/update/application"
	api "github.com/aziontech/azion-cli/pkg/api/applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

// Fields struct of inputs
type Fields struct {
	ID                            int64  `json:"id,omitempty"`
	Name                          string `json:"name,omitempty"`
	EdgeCacheEnabled              string `json:"edge_cache_enabled,omitempty"`
	FunctionsEnabled              string `json:"functions_enabled,omitempty"`
	ApplicationAcceleratorEnabled string `json:"application_accelerator_enabled,omitempty"`
	ImageProcessorEnabled         string `json:"image_processor_enabled,omitempty"`
	TieredCacheEnabled            string `json:"tiered_cache_enabled,omitempty"`
	Active                        string `json:"active,omitempty"`
	DebugRules                    string `json:"debug,omitempty"`
	Path                          string
	OutPath                       string
	Format                        string
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
		$ azion update application --application-id 1234 --name 'Hello'
		$ azion update application --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") && !cmd.Flags().Changed("file") {

				answer, err := utils.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				fields.ID = num
			}

			if !returnAnyField(cmd) {
				return msg.ErrorNoFieldInformed
			}

			request := api.UpdateRequest{}
			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				request.Id = fields.ID

				if cmd.Flags().Changed("name") {
					request.SetName(fields.Name)
				}

				if !utils.IsEmpty(fields.DebugRules) {
					debugRules, err := strconv.ParseBool(fields.DebugRules)
					if err != nil {
						logger.Debug("Error while parsing <"+fields.DebugRules+"> ", zap.Error(err))
						return utils.ErrorConvertingStringToBool
					}

					request.SetDebug(debugRules)
				}

				modules := sdk.ApplicationModulesRequest{}

				if !utils.IsEmpty(fields.EdgeCacheEnabled) {
					edgeCache, err := strconv.ParseBool(fields.EdgeCacheEnabled)
					if err != nil {
						logger.Debug("Error while parsing <"+fields.EdgeCacheEnabled+"> ", zap.Error(err))
						return utils.ErrorConvertingStringToBool
					}

					eCache := sdk.CacheModuleRequest{
						Enabled: &edgeCache,
					}

					modules.SetEdgeCache(eCache)
				}

				if !utils.IsEmpty(fields.FunctionsEnabled) {
					edgeFunctions, err := strconv.ParseBool(fields.FunctionsEnabled)
					if err != nil {
						logger.Debug("Error while parsing <"+fields.FunctionsEnabled+"> ", zap.Error(err))
						return utils.ErrorConvertingStringToBool
					}

					eFunction := sdk.EdgeFunctionModuleRequest{
						Enabled: &edgeFunctions,
					}

					modules.SetFunctions(eFunction)
				}

				if !utils.IsEmpty(fields.ApplicationAcceleratorEnabled) {
					applicationAcc, err := strconv.ParseBool(fields.ApplicationAcceleratorEnabled)
					if err != nil {
						logger.Debug("Error while parsing <"+fields.ApplicationAcceleratorEnabled+"> ", zap.Error(err))
						return utils.ErrorConvertingStringToBool
					}

					aAcceleration := sdk.ApplicationAcceleratorModuleRequest{
						Enabled: &applicationAcc,
					}

					modules.SetApplicationAccelerator(aAcceleration)
				}

				if !utils.IsEmpty(fields.ImageProcessorEnabled) {
					imageProcessor, err := strconv.ParseBool(fields.ImageProcessorEnabled)
					if err != nil {
						logger.Debug("Error while parsing <"+fields.ImageProcessorEnabled+"> ", zap.Error(err))
						return utils.ErrorConvertingStringToBool
					}

					iProcessor := sdk.ImageProcessorModuleRequest{
						Enabled: &imageProcessor,
					}

					modules.SetImageProcessor(iProcessor)
				}

				if !utils.IsEmpty(fields.TieredCacheEnabled) {
					tieredCache, err := strconv.ParseBool(fields.TieredCacheEnabled)
					if err != nil {
						logger.Debug("Error while parsing <"+fields.TieredCacheEnabled+"> ", zap.Error(err))
						return utils.ErrorConvertingStringToBool
					}

					tCache := sdk.TieredCacheModuleRequest{
						Enabled: &tieredCache,
					}

					modules.SetTieredCache(tCache)
				}

				request.SetModules(modules)

				if !utils.IsEmpty(fields.Active) {
					active, err := strconv.ParseBool(fields.Active)
					if err != nil {
						logger.Debug("Error while parsing <"+fields.Active+"> ", zap.Error(err))
						return utils.ErrorConvertingStringToBool
					}

					request.SetActive(active)
				}

			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, &request)

			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateApplication.Error(), err.Error())
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&updateOut)
		},
	}

	flags := cmd.Flags()
	flags.Int64Var(&fields.ID, "application-id", 0, msg.FlagID)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.EdgeCacheEnabled, "edge-cache", "", msg.FlagCaching)
	flags.StringVar(&fields.FunctionsEnabled, "functions", "", msg.FlagFunctions)
	flags.StringVar(&fields.ApplicationAcceleratorEnabled, "application-accelerator", "", msg.FlagApplicationAcceleration)
	flags.StringVar(&fields.ImageProcessorEnabled, "image-processor", "", msg.FlagImageOptimization)
	flags.StringVar(&fields.TieredCacheEnabled, "tiered-cache", "", msg.FlagTieredCaching)
	flags.StringVar(&fields.Active, "active", "", msg.FlagName)
	flags.StringVar(&fields.DebugRules, "debug-enabled", "", msg.FlagDebugRules)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

func returnAnyField(cmd *cobra.Command) bool {
	anyFlagChanged := false
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		if flag.Changed && flag.Name != "application-id" {
			anyFlagChanged = true
		}
	})
	return anyFlagChanged
}

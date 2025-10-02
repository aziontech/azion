package applications

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/create/application"
	api "github.com/aziontech/azion-cli/pkg/api/applications"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const example = `
        $ azion create application --name "naruno"
        $ azion create application --file create.json
        `

type Fields struct {
	Name                          string `json:"name,omitempty"`
	EdgeCacheEnabled              string `json:"edge_cache_enabled,omitempty"`
	EdgeFunctionsEnabled          string `json:"edge_functions_enabled,omitempty"`
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
		Example:       heredoc.Doc(example),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.CreateRequest{}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				err := createRequestFromFlags(fields, &request)
				if err != nil {
					return err
				}
			}

			response, err := api.NewClient(
				f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"),
			).Create(context.Background(), &request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreate.Error(), err)
			}

			creatOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, response.GetId()),
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

func createRequestFromFlags(fields *Fields, request *api.CreateRequest) error {

	if utils.IsEmpty(fields.Name) {
		answers, err := utils.AskInput("Enter the new Application's name")
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}

		fields.Name = answers
	}

	if utils.IsEmpty(fields.Name) {
		return msg.ErrorMandatoryCreateFlags
	}

	request.SetName(fields.Name)

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

	if !utils.IsEmpty(fields.EdgeFunctionsEnabled) {
		edgeFunctions, err := strconv.ParseBool(fields.EdgeFunctionsEnabled)
		if err != nil {
			logger.Debug("Error while parsing <"+fields.EdgeFunctionsEnabled+"> ", zap.Error(err))
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

	return nil
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.EdgeCacheEnabled, "edge-cache", "", msg.FlagCaching)
	flags.StringVar(&fields.EdgeFunctionsEnabled, "function", "", msg.FlagEdgeFunctions)
	flags.StringVar(&fields.ApplicationAcceleratorEnabled, "application-accelerator", "", msg.FlagApplicationAcceleration)
	flags.StringVar(&fields.ImageProcessorEnabled, "image-processor", "", msg.FlagImageOptimization)
	flags.StringVar(&fields.TieredCacheEnabled, "tiered-cache", "", msg.FlagTieredCaching)
	flags.StringVar(&fields.Active, "active", "", msg.FlagActive)
	flags.StringVar(&fields.DebugRules, "debug-enabled", "", msg.FlagDebugRules)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.FlagHelp)
}

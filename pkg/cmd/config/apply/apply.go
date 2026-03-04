package apply

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/config/apply"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/command"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/manifest"
	"github.com/aziontech/azion-cli/pkg/output"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ApplyCmd struct {
	GetWorkDir            func() (string, error)
	FileReader            func(path string) ([]byte, error)
	WriteFile             func(filename string, data []byte, perm fs.FileMode) error
	GetAzionJsonContent   func(confPath string) (*contracts.AzionApplicationOptions, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confPath string) error
	Interpreter           func() *manifest.ManifestInterpreter
	F                     *cmdutil.Factory
	CommandRunInteractive func(f *cmdutil.Factory, comm string) error
}

type Fields struct {
	ConfigDir       string
	AzionConfigName string
}

func NewApplyCmd(f *cmdutil.Factory) *ApplyCmd {
	return &ApplyCmd{
		GetWorkDir:            utils.GetWorkingDir,
		FileReader:            os.ReadFile,
		WriteFile:             os.WriteFile,
		GetAzionJsonContent:   utils.GetAzionJsonContent,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
		Interpreter:           manifest.NewManifestInterpreter,
		F:                     f,
		CommandRunInteractive: command.CommandRunInteractive,
	}
}

func NewCobraCmd(apply *ApplyCmd) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion config apply
        $ azion config apply --config-dir ./my-project
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return apply.Run(fields)
		},
	}

	cmd.Flags().StringVar(&fields.ConfigDir, "config-dir", ".", msg.FlagConfigDir)
	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)

	return cmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewApplyCmd(f))
}

func (cmd *ApplyCmd) Run(fields *Fields) error {
	msgs := []string{}
	logger.Debug("Running config apply command")

	wd, err := cmd.GetWorkDir()
	if err != nil {
		return err
	}

	configDir := fields.ConfigDir

	validConfigExtensions := []string{".js", ".ts", ".mjs", ".cjs"}
	azionConfigPath := ""
	for _, ext := range validConfigExtensions {
		candidate := path.Join(wd, "azion.config"+ext)
		if _, err := os.Stat(candidate); err == nil {
			azionConfigPath = candidate
			break
		}
	}
	if azionConfigPath == "" {
		logger.FInfoFlags(cmd.F.IOStreams.Out, msg.AzionConfigNotFound, cmd.F.Format, cmd.F.Out)
		return msg.ErrorAzionConfigNotFound
	}

	logger.Debug("Generating manifest from azion.config")
	vul := vulcanPkg.NewVulcan()
	command := vul.Command("", "manifest generate", cmd.F)
	logger.Debug("Running the following command", zap.Any("Command", command))

	err = cmd.CommandRunInteractive(cmd.F, command)
	if err != nil {
		return err
	}

	conf, err := cmd.GetAzionJsonContent(configDir)
	if err != nil {
		if errors.Is(err, utils.ErrorOpeningAzionJsonFile) {
			logger.FInfoFlags(cmd.F.IOStreams.Out, msg.CreatingAzionJson, cmd.F.Format, cmd.F.Out)
			conf = &contracts.AzionApplicationOptions{}
			if err := cmd.WriteAzionJsonContent(conf, configDir); err != nil {
				logger.Debug("Error creating azion.json file", zap.Error(err))
				return msg.ErrorCreatingAzionJson
			}
		} else {
			return err
		}
	}

	interpreter := cmd.Interpreter()

	manifestPath, err := interpreter.ManifestPath()
	if err != nil {
		return err
	}

	manifestStructure, err := interpreter.ReadManifest(manifestPath, cmd.F, &msgs)
	if err != nil {
		logger.Debug("Error reading manifest", zap.Error(err))
		return msg.ErrorReadingManifest
	}

	rc := manifest.NewResourceContext(
		cmd.F,
		conf,
		manifestStructure,
		configDir,
		&msgs,
		cmd.WriteAzionJsonContent,
	)

	resourceCount := 0

	if len(manifestStructure.Storage) > 0 {
		if err := rc.ApplyStorage(manifestStructure.Storage); err != nil {
			return err
		}
		resourceCount++
	}

	if len(manifestStructure.Connectors) > 0 {
		if err := rc.ApplyConnectors(manifestStructure.Connectors); err != nil {
			return err
		}
		resourceCount++
	}

	if len(manifestStructure.Functions) > 0 {
		if err := rc.ApplyFunctions(manifestStructure.Functions); err != nil {
			return err
		}
		resourceCount++
	}

	if len(manifestStructure.Applications) > 0 {
		app := manifestStructure.Applications[0]

		if err := rc.ApplyEdgeApplication(app); err != nil {
			return err
		}
		resourceCount++

		if len(app.FunctionsInstances) > 0 {
			if err := rc.ApplyFunctionInstances(app.FunctionsInstances); err != nil {
				return err
			}
			resourceCount++
		}

		if len(app.CacheSettings) > 0 {
			if err := rc.ApplyCacheSettings(app.CacheSettings); err != nil {
				return err
			}
		}

		if len(app.Rules) > 0 {
			if err := rc.ApplyRulesEngine(app.Rules); err != nil {
				return err
			}
		}
	}

	if len(manifestStructure.Workloads) > 0 {
		if err := rc.ApplyWorkloads(manifestStructure.Workloads); err != nil {
			return err
		}
		resourceCount++
	}

	if len(manifestStructure.WorkloadDeployments) > 0 {
		if err := rc.ApplyWorkloadDeployments(manifestStructure.WorkloadDeployments); err != nil {
			return err
		}
		resourceCount++
	}

	if len(manifestStructure.Firewalls) > 0 {
		if err := rc.ApplyFirewalls(manifestStructure.Firewalls); err != nil {
			return err
		}
		resourceCount++
	}

	if len(manifestStructure.Purge) > 0 {
		if err := rc.ApplyPurge(manifestStructure.Purge); err != nil {
			return err
		}
		resourceCount++
	}

	if resourceCount == 0 {
		logger.FInfoFlags(cmd.F.IOStreams.Out, msg.NoResourcesToApply, cmd.F.Format, cmd.F.Out)
		return nil
	}

	if err := rc.DeleteOrphanedResources(); err != nil {
		return err
	}

	successMsg := fmt.Sprintf(msg.ApplySuccessful, resourceCount)
	logger.FInfoFlags(cmd.F.IOStreams.Out, successMsg, cmd.F.Format, cmd.F.Out)
	msgs = append(msgs, successMsg)

	outSlice := output.SliceOutput{
		Messages: msgs,
		GeneralOutput: output.GeneralOutput{
			Out:   cmd.F.IOStreams.Out,
			Flags: cmd.F.Flags,
		},
	}

	return output.Print(&outSlice)
}

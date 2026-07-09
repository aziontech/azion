package init

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/config/init"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type InitCmd struct {
	GetWorkDir func() (string, error)
	FileReader func(path string) ([]byte, error)
	WriteFile  func(filename string, data []byte, perm fs.FileMode) error
	Stat       func(path string) (fs.FileInfo, error)
	F          *cmdutil.Factory
}

type Fields struct {
	ConfigDir string
	Force     bool
}

func NewInitCmd(f *cmdutil.Factory) *InitCmd {
	return &InitCmd{
		GetWorkDir: func() (string, error) {
			return os.Getwd()
		},
		FileReader: os.ReadFile,
		WriteFile:  os.WriteFile,
		Stat:       os.Stat,
		F:          f,
	}
}

func NewCobraCmd(init *InitCmd) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion config init
        $ azion config init --config-dir ./my-project
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return init.Run(fields)
		},
	}

	cmd.Flags().StringVar(&fields.ConfigDir, "config-dir", ".", msg.FlagConfigDir)
	cmd.Flags().BoolVar(&fields.Force, "force", false, "Overwrite existing configuration file")
	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)

	return cmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewInitCmd(f))
}

func (cmd *InitCmd) Run(fields *Fields) error {
	logger.Debug("Running config init command")

	wd, err := cmd.GetWorkDir()
	if err != nil {
		return err
	}

	configDir := fields.ConfigDir
	if configDir == "." {
		configDir = wd
	} else if !path.IsAbs(configDir) {
		configDir = path.Join(wd, configDir)
	}

	configPath := path.Join(configDir, "azion.json")

	if _, err := cmd.Stat(configPath); err == nil && !fields.Force {
		existsMsg := fmt.Sprintf(msg.ConfigExists, configPath)
		logger.FInfoFlags(cmd.F.IOStreams.Out, existsMsg, cmd.F.Format, cmd.F.Out)
		return msg.ErrorConfigExists
	}

	logger.FInfoFlags(cmd.F.IOStreams.Out, msg.CreatingConfig, cmd.F.Format, cmd.F.Out)

	azionJson := &contracts.AzionApplicationOptions{}

	data, err := json.MarshalIndent(azionJson, "", "  ")
	if err != nil {
		logger.Debug("Error marshaling azion.json", zap.Error(err))
		return msg.ErrorCreatingConfig
	}

	if err := cmd.WriteFile(configPath, data, 0644); err != nil {
		logger.Debug("Error creating config file", zap.Error(err))
		return msg.ErrorCreatingConfig
	}

	successMsg := fmt.Sprintf(msg.InitSuccessful, configPath)
	logger.FInfoFlags(cmd.F.IOStreams.Out, successMsg, cmd.F.Format, cmd.F.Out)
	logger.FInfoFlags(cmd.F.IOStreams.Out, msg.DocsURL, cmd.F.Format, cmd.F.Out)

	outSlice := output.SliceOutput{
		Messages: []string{successMsg, msg.DocsURL},
		GeneralOutput: output.GeneralOutput{
			Out:   cmd.F.IOStreams.Out,
			Flags: cmd.F.Flags,
		},
	}

	return output.Print(&outSlice)
}

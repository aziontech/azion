package ls

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var ErrNotFound = errors.New("executable file not found in $PATH")

type LsCmd struct {
	Io            *iostreams.IOStreams
	CommandRunner func(cmd string, envvars []string) (string, int, error)
	LookPath      func(bin string) (string, error)
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewLsCmd(f))
}

func NewLsCmd(f *cmdutil.Factory) *LsCmd {
	return &LsCmd{
		Io: f.IOStreams,
		CommandRunner: func(cmd string, envvars []string) (string, int, error) {
			return utils.RunCommandWithOutput(envvars, cmd)
		},
		LookPath: exec.LookPath,
	}
}

func NewCobraCmd(ls *LsCmd) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.LsUsage,
		Short:         msg.LsShortDescription,
		Long:          msg.LsLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azioncli edge_applications ls
		$ azioncli edge_applications ls --help
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return ls.run()
		},
	}

	return cobraCmd
}

func (cmd *LsCmd) run() error {
	_, err := cmd.LookPath("vulcan")
	if err != nil {
		if err == ErrNotFound {
			_, _, err := cmd.CommandRunner("npm install edge-functions -g", []string{})
			return err
		}
		return err
	}

	_, _, err = cmd.CommandRunner("npm update -g edge-functions", []string{})
	if err != nil {
		return fmt.Errorf("%s: %w", msg.ErrorUpdatingVulcan, err)
	}

	output, _, err := cmd.CommandRunner("vulcan presets ls", []string{"CLEAN_OUTPUT_MODE=true"})
	if err != nil {
		return err
	}

	newLineSplit := strings.Split(output, "\n")
	types := make([]string, len(newLineSplit))
	modes := make([]string, len(newLineSplit))
	for i := range newLineSplit {
		modeSplit := strings.Split(newLineSplit[i], " ")
		if len(modeSplit) > 1 {
			types[i] = modeSplit[0]
			modes[i] = modeSplit[1]
		}
	}

	tbl := table.New("PRESET", "MODE")
	table.DefaultWriter = cmd.Io.Out
	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt)

	for i := range types {
		tbl.AddRow(types[i], strings.Replace(strings.Replace(modes[i], "(", "", -1), ")", "", -1))
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})
	tbl.PrintHeader(format)

	for _, row := range tbl.GetRows() {
		tbl.PrintRow(format, row)
	}

	return nil
}

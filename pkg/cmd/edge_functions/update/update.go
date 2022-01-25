package update

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Name          string
	Language      string
	Code          string
	Active        string
	InitiatorType string
	Args          string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           "update <edge_function_id> [flags]",
		Short:         "Update an Edge Function",
		Long:          "Update an Edge Function",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions update 4185 –code ./mycode/function.js –args ./mycode/myargs.json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing edge function id argument")
			}

			_, err := utils.ConvertIdsToInt(args[0])
			if err != nil {
				return fmt.Errorf("invalid edge function id: %q", args[0])
			}

			if cmd.Flags().Changed("active") {
				_, err := strconv.ParseBool(fields.Active)
				if err != nil {
					return fmt.Errorf("invalid --active flag: %q", fields.Active)
				}
			}

			if cmd.Flags().Changed("code") {
				_, err := ioutil.ReadFile(fields.Code)
				if err != nil {
					return fmt.Errorf("failed to read code file: %w", err)
				}
			}

			if cmd.Flags().Changed("args") {
				_, err := ioutil.ReadFile(fields.Args)
				if err != nil {
					return fmt.Errorf("failed to read args file: %w", err)
				}
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Name, "name", "", "Name of your Edge Function.")
	flags.StringVar(&fields.Language, "language", "", "Programming language of your Edge Function <javascript|lua>")
	flags.StringVar(&fields.Code, "code", "", "Path to the file containing your Edge Function code.")
	flags.StringVar(&fields.InitiatorType, "initiator-type", "", "Initiator of your Edge Function: <edge-application|edge-firewall>")
	flags.StringVar(&fields.Active, "active", "", "Whether or not your Edge Function should be active: <true|false>")

	return cmd
}

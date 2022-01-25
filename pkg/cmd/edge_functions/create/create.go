package create

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
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
		Use:           "create [flags]",
		Short:         "Create a new Edge Function",
		Long:          "Create a new Edge Function",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions create -–name myfunc -–language javascript –-code ./mycode/function.js  -–state active --initiator-type edge-application
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := strconv.ParseBool(fields.Active)
			if err != nil {
				return fmt.Errorf("invalid --active flag: %s", fields.Active)
			}

			_, err = ioutil.ReadFile(fields.Code)
			if err != nil {
				return fmt.Errorf("failed to read code file: %w", err)
			}

			if cmd.Flags().Changed("args") {
				_, err = ioutil.ReadFile(fields.Args)
				if err != nil {
					return fmt.Errorf("failed to read args file: %w", err)
				}
			}

			// TODO: Interact with SDK to create function

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&fields.Name, "name", "", "Name of your Edge Function.")
	flags.StringVar(&fields.Language, "language", "", "Programming language of your Edge Function <javascript|lua>")
	flags.StringVar(&fields.Code, "code", "", "Path to the file containing your Edge Function code.")
	flags.StringVar(&fields.InitiatorType, "initiator-type", "", "Initiator of your Edge Function: <edge-application|edge-firewall>")
	flags.StringVar(&fields.Active, "active", "", "Whether or not your Edge Function should be active: <true|false>")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("language")
	_ = cmd.MarkFlagRequired("code")
	_ = cmd.MarkFlagRequired("initiator-type")
	_ = cmd.MarkFlagRequired("active")

	return cmd
}

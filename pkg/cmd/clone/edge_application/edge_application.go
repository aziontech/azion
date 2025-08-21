package edgeapplication

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/clone"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const example = `
        $ azion create edge-application --name "naruno"
        $ azion create edge-application --file create.json
        `

type Fields struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Path    string
	OutPath string
	Format  string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.UsageEdgeApplication,
		Short:         msg.ShortDescriptionEdgeApplication,
		Long:          msg.LongDescriptionEdgeApplication,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc(example),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("application-id") {
				answer, err := utils.AskInput(msg.AskApplicationIdClone)
				if err != nil {
					return err
				}
				fields.Id = answer
			}

			if !cmd.Flags().Changed("name") {
				answer, err := utils.AskInput(msg.AskApplicationNameClone)
				if err != nil {
					return err
				}

				fields.Name = answer
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			err := client.Clone(context.Background(), fields.Name, fields.Id)
			if err != nil {
				return fmt.Errorf(msg.ErrorClone.Error(), err)
			}

			creatOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.CloneSuccessful, fields.Name),
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

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.Name, "name", "", msg.FlagNameEdgeApplication)
	flags.StringVar(&fields.Id, "application-id", "", msg.FlagIdEdgeApplication)
	flags.BoolP("help", "h", false, msg.FlagHelpEdgeApplication)
}

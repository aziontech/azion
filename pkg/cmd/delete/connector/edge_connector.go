package connector

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/connector"
	api "github.com/aziontech/azion-cli/pkg/api/connector"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var connectorID string

type DeleteCmd struct {
	Io              *iostreams.IOStreams
	ReadInput       func(string) (string, error)
	DeleteConnector func(context.Context, string) error
	AskInput        func(string) (string, error)
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		DeleteConnector: func(ctx context.Context, connectorID string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.Delete(ctx, connectorID)
		},
		AskInput: utils.AskInput,
	}
}

func NewCobraCmd(delete *DeleteCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.DeleteShortDescription,
		Long:          msg.DeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion delete connector --connector-id 1234
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if !cmd.Flags().Changed("connector-id") {
				answer, err := delete.AskInput(msg.AskEdgeConnectorID)
				if err != nil {
					return err
				}

				connectorID = answer
			}

			ctx := context.Background()

			err = delete.DeleteConnector(ctx, connectorID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeleteConnector.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DeleteOutputSuccess, connectorID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().StringVar(&connectorID, "connector-id", "", msg.FlagID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DeleteHelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}

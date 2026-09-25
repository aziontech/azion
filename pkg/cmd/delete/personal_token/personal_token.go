package personaltoken

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type DeleteCmd struct {
	Io         *iostreams.IOStreams
	AskInput   func(string) (string, error)
	DeleteFunc func(context.Context, string) error
	TokenID    string
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		DeleteFunc: func(ctx context.Context, id string) error {
			client := registry.PersonalToken(f)
			return client.Delete(ctx, id)
		},
	}
}

func NewCobraCmd(delete *DeleteCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		  $ azion delete personal-token --id 1234-123-321
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if !cmd.Flags().Changed("id") {
				answer, err := delete.AskInput(msg.AskDeleteInput)
				if err != nil {
					return err
				}
				delete.TokenID = answer
			}

			if utils.IsEmpty(delete.TokenID) {
				return utils.ErrorArgumentIsEmpty
			}

			ctx := context.Background()

			err = delete.DeleteFunc(ctx, delete.TokenID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDelete, err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, delete.TokenID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().StringVar(&delete.TokenID, "id", "", msg.FlagID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)
	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}

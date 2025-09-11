package functioninstance

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/function_instance"
	api "github.com/aziontech/azion-cli/pkg/api/function_instance"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	edgeApplicationID  string
	functionInstanceID string
)

type DeleteCmd struct {
	Io        *iostreams.IOStreams
	ReadInput func(string) (string, error)
	AskInput  func(string) (string, error)
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		AskInput: utils.AskInput,
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
		$ azion delete function-instance --application-id 1234 --instance-id 4321
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if !cmd.Flags().Changed("application-id") {
				answer, err := delete.AskInput(msg.AskDeleteInput)
				if err != nil {
					return err
				}

				edgeApplicationID = answer
			}

			if !cmd.Flags().Changed("instance-id") {
				answer, err := delete.AskInput(msg.AskDeleteInput)
				if err != nil {
					return err
				}

				functionInstanceID = answer
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err = client.Delete(ctx, edgeApplicationID, functionInstanceID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeletInstance.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, functionInstanceID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().StringVar(&edgeApplicationID, "application-id", "", msg.FlagId)
	cobraCmd.Flags().StringVar(&functionInstanceID, "instance-id", "", msg.FlagId)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}

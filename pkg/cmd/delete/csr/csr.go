package csr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/csr"
	api "github.com/aziontech/azion-cli/pkg/api/csr"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var csrID int64

type DeleteCmd struct {
	Io        *iostreams.IOStreams
	DeleteCSR func(context.Context, int64) error
	AskInput  func(string) (string, error)
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io: f.IOStreams,
		DeleteCSR: func(ctx context.Context, id int64) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.Delete(ctx, id)
		},
		AskInput: utils.AskInput,
	}
}

func NewCobraCmd(del *DeleteCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion delete csr --csr-id 1234
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("csr-id") {
				answer, err := del.AskInput(msg.AskInput)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertId
				}

				csrID = num
			}

			ctx := context.Background()

			err := del.DeleteCSR(ctx, csrID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeleteCSR, err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, csrID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().Int64Var(&csrID, "csr-id", 0, msg.FlagId)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}

package application

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/application"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// DeleteCmd has no --cascade, unlike its v4 namesake: cascade removes what
// `azion deploy` created, and the v6 deploy model is not defined yet.
type DeleteCmd struct {
	Io       *iostreams.IOStreams
	f        *cmdutil.Factory
	AskInput func(string) (string, error)
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f))
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io:       f.IOStreams,
		f:        f,
		AskInput: utils.AskInput,
	}
}

func NewCobraCmd(delete *DeleteCmd) *cobra.Command {
	var application_id int64
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion delete application --application-id 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete.run(cmd, application_id)
		},
	}

	cmd.Flags().Int64Var(&application_id, "application-id", 0, msg.FlagId)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

func (del *DeleteCmd) run(cmd *cobra.Command, application_id int64) error {
	ctx := context.Background()

	if !cmd.Flags().Changed("application-id") {
		answer, err := del.AskInput(msg.AskInput)
		if err != nil {
			return err
		}

		num, err := strconv.ParseInt(answer, 10, 64)
		if err != nil {
			logger.Debug("Error while converting answer to int64", zap.Error(err))
			return msg.ErrorConvertId
		}

		application_id = num
	}

	client := registry.ApplicationsV6(del.f)
	err := client.Delete(ctx, application_id)
	if err != nil {
		return fmt.Errorf(msg.ErrorFailToDeleteApplication.Error(), err)
	}

	deleteOut := output.GeneralOutput{
		Msg:   fmt.Sprintf(msg.OutputSuccess, application_id),
		Out:   del.f.IOStreams.Out,
		Flags: del.f.Flags,
	}
	return output.Print(&deleteOut)
}

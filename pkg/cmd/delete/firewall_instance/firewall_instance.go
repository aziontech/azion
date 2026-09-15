package firewallinstance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/firewall_instance"
	api "github.com/aziontech/azion-cli/pkg/api/firewall_instance"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DeleteCmd struct {
	Io                         *iostreams.IOStreams
	ReadInput                  func(string) (string, error)
	AskInput                   func(string) (string, error)
	FirewallFunctionInstanceID int64
	FirewallID                 int64
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
		$ azion delete firewall-instance --firewall-id 1234 --instance-id 4321
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if !cmd.Flags().Changed("firewall-id") {
				answer, err := delete.AskInput(msg.AskDeleteFirewallInput)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertFirewallId
				}

				delete.FirewallID = num
			}

			if !cmd.Flags().Changed("instance-id") {
				answer, err := delete.AskInput(msg.AskDeleteInput)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertFirewallFunctionInstanceId
				}

				delete.FirewallFunctionInstanceID = num
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err = client.Delete(ctx, delete.FirewallID, delete.FirewallFunctionInstanceID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeletInstance.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, delete.FirewallFunctionInstanceID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().Int64Var(&delete.FirewallID, "firewall-id", 0, msg.FlagId)
	cobraCmd.Flags().Int64Var(&delete.FirewallFunctionInstanceID, "instance-id", 0, msg.FlagId)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}

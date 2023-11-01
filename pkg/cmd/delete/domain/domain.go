package domain

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/domain"
	api "github.com/aziontech/azion-cli/pkg/api/domain"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var domain_id int64
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion delete domain --domain-id 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("domain-id") {

				answer, err := utils.AskInput(msg.AskDeleteInput)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertId
				}

				domain_id = num
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.Delete(ctx, domain_id)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeleteDomain.Error(), err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, msg.OutputSuccess, domain_id)

			return nil
		},
	}

	cmd.Flags().Int64Var(&domain_id, "domain-id", 0, msg.FlagId)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cmd
}

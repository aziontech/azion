package workloaddeployment

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/update/workloads"
	api "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Fields struct {
	Name              string
	WorkloadId        int64
	AlternateDomains  []string
	Active            string
	EdgeApplicationID int64
	EdgeFirewall      int64
	Path              string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion update workload --workload-id 1234 --name 'Hello'
		$ azion update workload --workload-id 1234 --alternate-domains www.testhere.com,www.pudim.com
		$ azion update workload --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.UpdateRequest{}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("workload-id") {
					answer, err := utils.AskInput(msg.AskInputWorkloadID)
					if err != nil {
						return err
					}
					num, err := strconv.ParseInt(answer, 10, 64)
					if err != nil {
						logger.Debug("Error while converting answer to int64", zap.Error(err))
						return msg.ErrorConvertWorkloadID

					}
					logger.Debug("Converted Workload ID", zap.Any("Workload ID", num))
					request.Id = num
				}

				if cmd.Flags().Changed("name") {
					request.SetName(fields.Name)
				}

				if len(fields.AlternateDomains) > 0 {
					request.SetAlternateDomains(fields.AlternateDomains)
				}

				if cmd.Flags().Changed("active") {
					isActive, err := strconv.ParseBool(fields.Active)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorActiveFlag, fields.Active)
					}
					request.SetActive(isActive)
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			ctx := context.Background()
			response, err := client.Update(ctx, &request)

			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateDomain.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&updateOut)
		},
	}

	flags := cmd.Flags()
	flags.Int64Var(&fields.WorkloadId, "domain-id", 0, msg.FlagWorkloadID)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringSliceVar(&fields.AlternateDomains, "alternate-domains", []string{}, msg.FlagAlternateDomains)
	flags.StringVar(&fields.Active, "active", "true", msg.FlagActive)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.HelpFlag)

	return cmd
}

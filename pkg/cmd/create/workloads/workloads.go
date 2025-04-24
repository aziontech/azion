package workloads

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/create/workloads"
	api "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Name              string   `json:"name"`
	AlternateDomains  []string `json:"alternate_domains"`
	Active            string   `json:"active"`
	EdgeApplicationID int64    `json:"edge_application"`
	EdgeFirewall      int64    `json:"edge_firewall"`
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
        $ azion create workload --edge-application 1231 --name workloadName
        $ azion create workload --name withargs --edge-application 1231 --active true
        $ azion create workload --alternate-domains "www.thisismydomain.com" --edge-application 1231
        $ azion create workload --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := new(api.CreateRequest)
			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("edge-application") {
					answer, err := utils.AskInput(msg.AskInputApplicationID)
					if err != nil {
						return err
					}
					num, err := strconv.ParseInt(answer, 10, 64)
					if err != nil {
						logger.Debug("Error while converting answer to int64", zap.Error(err))
						return msg.ErrorConvertApplicationID

					}
					fields.EdgeApplicationID = num
				}

				if !cmd.Flags().Changed("name") {
					answer, err := utils.AskInput(msg.AskInputName)
					if err != nil {
						return err
					}

					fields.Name = answer
				}

				request.SetName(fields.Name)
				if len(fields.AlternateDomains) > 0 {
					request.SetAlternateDomains(fields.AlternateDomains)
				}

				request.SetEdgeApplication(fields.EdgeApplicationID)

				isActive, err := strconv.ParseBool(fields.Active)
				if err != nil {
					return fmt.Errorf("%w: %q", msg.ErrorIsActiveFlag, fields.Active)
				}
				request.SetActive(isActive)

				if fields.EdgeFirewall > 0 {
					request.SetEdgeFirewall(fields.EdgeFirewall)
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			response, err := client.Create(context.Background(), request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreate.Error(), err)
			}

			createOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&createOut)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringSliceVar(&fields.AlternateDomains, "alternate-domains", []string{}, msg.FlagAlternateDomains)
	flags.Int64Var(&fields.EdgeFirewall, "edge-firewall", 0, msg.FlagEdgeFirewall)
	flags.Int64Var(&fields.EdgeApplicationID, "edge-application", 0, msg.FlagEdgeApplicationId)
	flags.StringVar(&fields.Active, "active", "true", msg.FlagIsActive)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

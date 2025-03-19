package workloads

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/create/domain"
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
	NetworkMap        string   `json:"network_map"`
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
        $ azion create domain --application-id 1231 --name domainName --cnames "asdf.com,asdfsdf.com,asdfd.com" --cname-access-only false
        $ azion create domain --name withargs --application-id 1231 --active true
		$ azion create domain --digital-certificate-id "lets_encrypt" --cnames "www.thisismycname.com" --application-id 1231
        $ azion create domain --file "create.json"
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
				if !cmd.Flags().Changed("application-id") {
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
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
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
	flags.Int64Var(&fields.EdgeFirewall, "edge_firewall", 0, msg.FlagDigitalCertificateID)
	flags.Int64Var(&fields.EdgeApplicationID, "application-id", 0, msg.FlagEdgeApplicationId)
	flags.StringVar(&fields.Active, "active", "true", msg.FlagIsActive)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

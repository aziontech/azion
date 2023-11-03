package domain

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/create/domain"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	api "github.com/aziontech/azion-cli/pkg/api/domain"
	"github.com/aziontech/azion-cli/pkg/logger"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Name                 string   `json:"name"`
	Cnames               []string `json:"cnames"`
	CnameAccessOnly      string   `json:"cname_access_only"`
	EdgeApplicationID    int      `json:"edge_application_id"`
	DigitalCertificateID int      `json:"digital_certificate_id"`
	IsActive             string   `json:"is_active"`
	Path                 string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           domain.Usage,
		Short:         domain.ShortDescription,
		Long:          domain.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create domain --application-id 1231 --name domainName --cnames "asdf.com,asdfsdf.com,asdfd.com" --cname-access-only false
        $ azion create domain --name withargs --application-id 1231 --active true
        $ azion create domain --in "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := new(api.CreateRequest)
			if cmd.Flags().Changed("in") {
				err := utils.FlagINUnmarshalFileJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("application-id") {
					answer, err := utils.AskInput(domain.AskInputApplicationID)
					if err != nil {
						return err
					}

					num, err := strconv.ParseInt(answer, 10, 64)
					if err != nil {
						logger.Debug("Error while converting answer to int64", zap.Error(err))
						return domain.ErrorConvertApplicationID

					}

					fields.EdgeApplicationID = int(num)
				}

				if !cmd.Flags().Changed("name") {
					answer, err := utils.AskInput(domain.AskInputName)
					if err != nil {
						return err
					}

					fields.Name = answer
				}

				cnameAccessOnly, err := strconv.ParseBool(fields.CnameAccessOnly)
				if err != nil {
					return fmt.Errorf("%w: %q", domain.ErrorCnameAccessOnlyFlag, fields.CnameAccessOnly)
				}
				request.SetCnameAccessOnly(cnameAccessOnly)

				if cnameAccessOnly {
					if len(fields.Cnames) < 1 {
						return domain.ErrorMissingCnames
					}
				}

				request.SetName(fields.Name)
				request.SetCnames(fields.Cnames)
				request.SetEdgeApplicationId(int64(fields.EdgeApplicationID))
				if fields.DigitalCertificateID > 0 {
					request.SetDigitalCertificateId(int64(fields.DigitalCertificateID))
				}

				isActive, err := strconv.ParseBool(fields.IsActive)
				if err != nil {
					return fmt.Errorf("%w: %q", domain.ErrorIsActiveFlag, fields.IsActive)
				}
				request.SetIsActive(isActive)

			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.Create(context.Background(), request)
			if err != nil {
				return fmt.Errorf(domain.ErrorCreate.Error(), err)
			}

			logger.LogSuccess(f.IOStreams.Out, fmt.Sprintf(domain.OutputSuccess, response.GetId()))
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Name, "name", "", domain.FlagName)
	flags.StringSliceVar(&fields.Cnames, "cnames", []string{}, domain.FlagCnames)
	flags.StringVar(&fields.CnameAccessOnly, "cname-access-only", "false", domain.FlagCnameAccessOnly)
	flags.IntVar(&fields.DigitalCertificateID, "digital-certificate-id", 0, domain.FlagDigitalCertificateID)
	flags.IntVar(&fields.EdgeApplicationID, "application-id", 0, domain.FlagEdgeApplicationId)
	flags.StringVar(&fields.IsActive, "active", "true", domain.FlagIsActive)
	flags.StringVar(&fields.Path, "in", "", domain.FlagIn)
	flags.BoolP("help", "h", false, domain.HelpFlag)
	return cmd
}

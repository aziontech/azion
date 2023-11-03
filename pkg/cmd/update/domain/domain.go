package domain

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/update/domain"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/domain"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Fields struct {
	DomainID           int64
	ApplicationID      int64
	Name               string
	CnameAccessOnly    string
	Active             string
	InPath             string
	Cnames             []string
	DigitalCertificate string
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
		$ azion update domain --domain-id 1234 --name 'Hello'
		$ azion update domain --domain-id 1234 --application-id 4321
		$ azion update domain --domain-id 1234 --cnames www.testhere.com,www.pudim.com
		$ azion update domain --domain-id 9123 --cname-access-only true
		$ azion update domain --domain-id 9123 --cname-access-only false
		$ azion update domain --domain-id 9123 --application-id 192837
		$ azion update domain --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.UpdateRequest{}

			if cmd.Flags().Changed("in") {
				err := utils.FlagINUnmarshalFileJSON(fields.InPath, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.InPath+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("domain-id") {
					answer, err := utils.AskInput(domain.AskInputDomainID)
					if err != nil {
						return err
					}

					num, err := strconv.ParseInt(answer, 10, 64)
					if err != nil {
						logger.Debug("Error while converting answer to int64", zap.Error(err))
						return domain.ErrorConvertDomainID
					}

					fields.DomainID = num
				}

				request.Id = fields.DomainID

				if cmd.Flags().Changed("name") {
					request.SetName(fields.Name)
				}

				if cmd.Flags().Changed("application-id") {
					request.SetEdgeApplicationId(fields.ApplicationID)
				}

				if cmd.Flags().Changed("cnames") {
					request.SetCnames(fields.Cnames)
				}

				if cmd.Flags().Changed("digital-certificate-id") {
					if fields.DigitalCertificate == "null" {
						request.SetDigitalCertificateIdNil()
					} else {
						n, err := strconv.ParseInt(fields.DigitalCertificate, 10, 64)
						if err != nil {
							return domain.ErrorDigitalCertificateFlag
						}
						request.SetDigitalCertificateId(n)
					}
				}

				if cmd.Flags().Changed("cname-access-only") {
					active, err := strconv.ParseBool(fields.CnameAccessOnly)
					if err != nil {
						return fmt.Errorf("%w: %q", domain.ErrorActiveFlag, fields.CnameAccessOnly)
					}
					request.SetCnameAccessOnly(active)
				}

				if cmd.Flags().Changed("active") {
					active, err := strconv.ParseBool(fields.Active)
					if err != nil {
						return fmt.Errorf("%w: %q", domain.ErrorActiveFlag, fields.Active)
					}
					request.SetIsActive(active)
				}

			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, &request)

			if err != nil {
				return fmt.Errorf(domain.ErrorUpdateDomain.Error(), err)
			}

			logger.LogSuccess(f.IOStreams.Out, fmt.Sprintf(domain.OutputSuccess, response.GetId()))
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64Var(&fields.DomainID, "domain-id", 0, domain.FlagDomainID)
	flags.Int64Var(&fields.ApplicationID, "application-id", 0, domain.FlagApplicationID)
	flags.StringVar(&fields.DigitalCertificate, "digital-certificate-id", "", domain.FlagDigitalCertificateID)
	flags.StringVar(&fields.Name, "name", "", domain.FlagName)
	flags.StringSliceVar(&fields.Cnames, "cnames", []string{}, domain.FlagCnames)
	flags.StringVar(&fields.Active, "active", "", domain.FlagActive)
	flags.StringVar(&fields.CnameAccessOnly, "cname-access-only", "", domain.FlagCnameAccessOnly)
	flags.StringVar(&fields.InPath, "in", "", domain.FlagIn)
	flags.BoolP("help", "h", false, domain.HelpFlag)

	return cmd
}

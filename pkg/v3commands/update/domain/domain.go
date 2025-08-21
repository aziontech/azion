package domain

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/update/domain"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	api "github.com/aziontech/azion-cli/pkg/v3api/domain"
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
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion update domain --domain-id 1234 --name 'Hello'
		$ azion update domain --domain-id 1234 --application-id 4321
		$ azion update domain --domain-id 1234 --cnames www.testhere.com,www.pudim.com
		$ azion update domain --domain-id 9123 --cname-access-only true
		$ azion update domain --domain-id 9123 --cname-access-only false
		$ azion update domain --domain-id 9123 --application-id 192837
		$ azion update domain --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.UpdateRequest{}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.InPath+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("domain-id") {
					answer, err := utils.AskInput(msg.AskInputDomainID)
					if err != nil {
						return err
					}

					num, err := strconv.ParseInt(answer, 10, 64)
					if err != nil {
						logger.Debug("Error while converting answer to int64", zap.Error(err))
						return msg.ErrorConvertDomainID
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
					request.SetDigitalCertificateId(fields.DigitalCertificate)
				}

				if cmd.Flags().Changed("cname-access-only") {
					active, err := strconv.ParseBool(fields.CnameAccessOnly)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorActiveFlag, fields.CnameAccessOnly)
					}
					request.SetCnameAccessOnly(active)
				}

				if cmd.Flags().Changed("active") {
					active, err := strconv.ParseBool(fields.Active)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorActiveFlag, fields.Active)
					}
					request.SetIsActive(active)
				}

			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

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
	flags.Int64Var(&fields.DomainID, "domain-id", 0, msg.FlagDomainID)
	flags.Int64Var(&fields.ApplicationID, "application-id", 0, msg.FlagApplicationID)
	flags.StringVar(&fields.DigitalCertificate, "digital-certificate-id", "", msg.FlagDigitalCertificateID)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringSliceVar(&fields.Cnames, "cnames", []string{}, msg.FlagCnames)
	flags.StringVar(&fields.Active, "active", "", msg.FlagActive)
	flags.StringVar(&fields.CnameAccessOnly, "cname-access-only", "", msg.FlagCnameAccessOnly)
	flags.StringVar(&fields.InPath, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.HelpFlag)

	return cmd
}

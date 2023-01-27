package create

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/domains"
	api "github.com/aziontech/azion-cli/pkg/api/domains"

	"github.com/aziontech/azion-cli/utils"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type Fields struct {
	Name                 string   `json:"name"`
	Cnames               []string `json:"cnames,omitempty"`
	EdgeApplicationId    int      `json:"edge_application_id"`
	DigitalCertificateId int      `json:"digital_certificate_id,omitempty"`
	CnameAccessOnly      bool     `json:"cname_access-only,omitempty"`
	IsActive             bool     `json:"is_active,omitempty"`
	Path                 string   `json:"path,omitempty"`
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.DomainsCreateUsage,
		Short:         msg.DomainsCreateShortDescription,
		Long:          msg.DomainsCreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli domains create --name asdfçlkj --cnames "asdf.com,asdfsdf.com,asdfd.com" --cname-access-only false
        $ azioncli domains create --name withargs --edge-application-id 1231 --active true
        $ azioncli domains create --in "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("in") {
				f, err := os.ReadFile(fields.Path)
				if err != nil {
					return fmt.Errorf("%s %s", utils.ErrorOpeningFile, fields.Path)
				}

				err = json.Unmarshal(f, &fields)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			}

			if len(fields.Name) < 1 || fields.EdgeApplicationId < 1 {
				return msg.ErrorMandatoryCreateFlags
			}

			request := new(api.CreateRequest)
			request.SetName(fields.Name)
			request.SetCnames(fields.Cnames)
			request.SetCnameAccessOnly(fields.CnameAccessOnly)
			request.SetEdgeApplicationId(int64(fields.EdgeApplicationId))
			if fields.DigitalCertificateId > 0 {
				request.SetDigitalCertificateId(int64(fields.DigitalCertificateId))
			}
			request.SetIsActive(fields.IsActive)
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.Create(context.Background(), request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateDomain.Error(), err)
			}

			fmt.Fprintf(f.IOStreams.Out, msg.DomainsCreateOutputSuccess, response.GetId())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Name, "name", "", msg.DomainsCreateFlagName)
	flags.StringSliceVar(&fields.Cnames, "cnames", []string{}, msg.DomainsCreateFlagCnames)
	flags.BoolVar(&fields.CnameAccessOnly, "cname-access-only", false, msg.DomainsCreateFlagCnameAccessOnly)
	flags.IntVarP(&fields.DigitalCertificateId, "digital-certificate-id", "d", 0, msg.DomainsCreateFlagDigitalCertificateId)
	flags.IntVarP(&fields.EdgeApplicationId, "edge-application-id", "e", 0, msg.DomainsCreateFlagEdgeApplicationId)
	flags.BoolVar(&fields.IsActive, "active", false, msg.DomainsCreateFlagIsActive)
	flags.StringVar(&fields.Path, "in", "", msg.DomainsCreateFlagIn)
	flags.BoolP("help", "h", false, msg.DomainsCreateHelpFlag)
	return cmd
}

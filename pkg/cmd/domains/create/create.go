package create

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	//msg "github.com/aziontech/azion-cli/messages/edge_functions"
	msg "github.com/aziontech/azion-cli/messages/domains"
	api "github.com/aziontech/azion-cli/pkg/api/domains"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type Fields struct {
	Cnames                                  []string
	Name, Path                              string
	EdgeApplicationId, DigitalCertificateId int
	CnameAccessOnly, IsActive               bool
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
		$ azioncli domains create --name "max" --cnames  "asdfg.com",max.com,123.com -c true -d 1234 -e 42312434 -a true --in example.json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, _ := json.MarshalIndent(fields, " ", "	")
			fmt.Println(string(b))

			request := new(api.CreateRequest)

			if !cmd.Flags().Changed("name") || !cmd.Flags().Changed("edge-application-id") {
				return msg.ErrorMandatoryCreateFlags
			}

			request.SetName(fields.Name)
			request.SetCnames(fields.Cnames)
			request.SetCnameAccessOnly(fields.CnameAccessOnly)
			request.SetEdgeApplicationId(int64(fields.EdgeApplicationId))
			request.SetDigitalCertificateId(int64(fields.DigitalCertificateId))
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
	flags.BoolVarP(&fields.CnameAccessOnly, "cname-access-only", "c", false, msg.DomainsCreateFlagCnameAccessOnly)
	flags.IntVarP(&fields.DigitalCertificateId, "digital-certificate-id", "d", 0, msg.DomainsCreateFlagDigitalCertificateId)
	flags.IntVarP(&fields.EdgeApplicationId, "edge-application-id", "e", 0, msg.DomainsCreateFlagEdgeApplicationId)
	flags.BoolVarP(&fields.IsActive, "active", "a", false, msg.DomainsCreateFlagIsActive)
	flags.StringVar(&fields.Path, "in", "", msg.DomainsCreateFlagIn)
	flags.BoolP("help", "h", false, msg.DomainsCreateHelpFlag)
	return cmd
}

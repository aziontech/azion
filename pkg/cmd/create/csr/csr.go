package csr

import (
	"context"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/create/csr"
	api "github.com/aziontech/azion-cli/pkg/api/csr"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Fields struct {
	Name              string
	CommonName        string
	Country           string
	State             string
	Locality          string
	Organization      string
	OrganizationUnity string
	Email             string
	AlternativeNames  string
	CertificateType   string
	KeyAlgorithm      string
	InPath            string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create csr --name "My CSR" --common-name "example.com" --country "US" --state "California" --locality "San Francisco" --organization "Example Corp" --organization-unity "IT" --email "admin@example.com"
        $ azion create csr --name "My CSR" --common-name "example.com" --country "US" --state "California" --locality "San Francisco" --organization "Example Corp" --organization-unity "IT" --email "admin@example.com" --alternative-names "www.example.com,api.example.com" --key-algorithm rsa_2048
        $ azion create csr --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			ctx := context.Background()

			req := api.NewCreateRequest()

			if cmd.Flags().Changed("file") {
				if err := utils.FlagFileUnmarshalJSON(fields.InPath, &req.CertificateSigningRequest); err != nil {
					return msg.ErrorInvalidJSON
				}
				return doCreate(ctx, f, client, req)
			}

			if err := requestFromFlags(cmd, fields, req); err != nil {
				return err
			}

			return doCreate(ctx, f, client, req)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)

	return cmd
}

func requestFromFlags(cmd *cobra.Command, fields *Fields, req *api.CreateRequest) error {
	if err := askIfEmpty(cmd, "name", &fields.Name, msg.AskName); err != nil {
		return err
	}
	req.SetName(fields.Name)

	if err := askIfEmpty(cmd, "common-name", &fields.CommonName, msg.AskCommonName); err != nil {
		return err
	}
	req.SetCommonName(fields.CommonName)

	if err := askIfEmpty(cmd, "country", &fields.Country, msg.AskCountry); err != nil {
		return err
	}
	req.SetCountry(fields.Country)

	if err := askIfEmpty(cmd, "state", &fields.State, msg.AskState); err != nil {
		return err
	}
	req.SetState(fields.State)

	if err := askIfEmpty(cmd, "locality", &fields.Locality, msg.AskLocality); err != nil {
		return err
	}
	req.SetLocality(fields.Locality)

	if err := askIfEmpty(cmd, "organization", &fields.Organization, msg.AskOrganization); err != nil {
		return err
	}
	req.SetOrganization(fields.Organization)

	if err := askIfEmpty(cmd, "organization-unity", &fields.OrganizationUnity, msg.AskOrganizationUnity); err != nil {
		return err
	}
	req.SetOrganizationUnity(fields.OrganizationUnity)

	if err := askIfEmpty(cmd, "email", &fields.Email, msg.AskEmail); err != nil {
		return err
	}
	req.SetEmail(fields.Email)

	if cmd.Flags().Changed("alternative-names") {
		names := strings.Split(fields.AlternativeNames, ",")
		for i := range names {
			names[i] = strings.TrimSpace(names[i])
		}
		req.SetAlternativeNames(names)
	}

	if cmd.Flags().Changed("certificate-type") {
		req.SetType(fields.CertificateType)
	}

	if cmd.Flags().Changed("key-algorithm") {
		req.SetKeyAlgorithm(fields.KeyAlgorithm)
	}

	return nil
}

func askIfEmpty(cmd *cobra.Command, flag string, target *string, prompt string) error {
	if cmd.Flags().Changed(flag) {
		return nil
	}
	answer, err := utils.AskInput(prompt)
	if err != nil {
		return err
	}
	*target = answer
	return nil
}

func doCreate(ctx context.Context, f *cmdutil.Factory, client *api.Client, req *api.CreateRequest) error {
	response, err := client.Create(ctx, req)
	if err != nil {
		return fmt.Errorf(msg.ErrorCreateCSR.Error(), err)
	}

	out := output.GeneralOutput{
		Msg:   fmt.Sprintf(msg.CreateOutputSuccess, response.GetId()),
		Out:   f.IOStreams.Out,
		Flags: f.Flags,
	}
	if err := output.Print(&out); err != nil {
		return err
	}

	if csr := response.GetCsr(); csr != "" {
		logger.FInfo(f.IOStreams.Out, fmt.Sprintf("\n%s\n", csr))
	}

	return nil
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.CommonName, "common-name", "", msg.FlagCommonName)
	flags.StringVar(&fields.Country, "country", "", msg.FlagCountry)
	flags.StringVar(&fields.State, "state", "", msg.FlagState)
	flags.StringVar(&fields.Locality, "locality", "", msg.FlagLocality)
	flags.StringVar(&fields.Organization, "organization", "", msg.FlagOrganization)
	flags.StringVar(&fields.OrganizationUnity, "organization-unity", "", msg.FlagOrganizationUnity)
	flags.StringVar(&fields.Email, "email", "", msg.FlagEmail)
	flags.StringVar(&fields.AlternativeNames, "alternative-names", "", msg.FlagAlternativeNames)
	flags.StringVar(&fields.CertificateType, "certificate-type", "", msg.FlagCertificateType)
	flags.StringVar(&fields.KeyAlgorithm, "key-algorithm", "", msg.FlagKeyAlgorithm)
	flags.StringVar(&fields.InPath, "file", "", msg.FlagIn)
	flags.BoolP("help", "h", false, msg.CreateFlagHelp)
}

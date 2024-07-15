package domains

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/domain"
	api "github.com/aziontech/azion-cli/pkg/api/domain"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	domainID string
)

type DescribeCmd struct {
	Io       *iostreams.IOStreams
	AskInput func(string) (string, error)
	Get      func(context.Context, string) (api.DomainResponse, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		Get: func(ctx context.Context, domainID string) (api.DomainResponse, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.Get(ctx, domainID)
		},
	}
}

func NewCobraCmd(describe *DescribeCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cobraCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion describe domain --domain-id 4312
		$ azion describe domain --domain-id 1337 --out "./tmp/test.json" --format json
		$ azion describe domain --domain-id 1337 --format json
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !cmd.Flags().Changed("domain-id") {
				answer, err := describe.AskInput(msg.AskInputDomainID)
				if err != nil {
					return err
				}

				domainID = answer
			}

			ctx := context.Background()
			domain, err := describe.Get(ctx, domainID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetDomain.Error(), err.Error())
			}

			fields := make(map[string]string)
			fields["Id"] = "ID"
			fields["Name"] = "Name"
			fields["DomainName"] = "Domain"
			fields["CnameAccessOnly"] = "Cname Access Only"
			fields["Cnames"] = "Cnames"
			fields["EdgeApplicationId"] = "Application ID"
			fields["DigitalCertificateId"] = "Digital Certificate ID"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
					Out:   f.IOStreams.Out,
				},
				Fields: fields,
				Values: domain,
			}
			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().StringVar(&domainID, "domain-id", "", msg.FlagDomainID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}

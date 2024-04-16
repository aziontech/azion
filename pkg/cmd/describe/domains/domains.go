package domains

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/domain"

	api "github.com/aziontech/azion-cli/pkg/api/domain"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var domainID string
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
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
				answer, err := utils.AskInput(msg.AskInputDomainID)
				if err != nil {
					return err
				}

				domainID = answer
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			domain, err := client.Get(ctx, domainID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetDomain.Error(), err.Error())
			}

			domainID := strconv.FormatInt(domain.GetId(), 10)
			cnameAccessOnly := strconv.FormatBool(domain.GetCnameAccessOnly())
			edgeApplicationID := strconv.FormatInt(domain.GetDigitalCertificateId(), 10)
			digitalCertifacateID := strconv.FormatInt(domain.GetDigitalCertificateId(), 10)

			fields := [][]string{
				{"ID", domainID},
				{"Name", domain.GetName()},
				{"Domain", domain.GetDomainName()},
				{"Cname Access Only", cnameAccessOnly},
				{"Cnames", fmt.Sprintf("%v", domain.GetCnames())},
				{"Application ID", edgeApplicationID},
				{"Digital Certificate ID", digitalCertifacateID},
			}

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Msg:         fmt.Sprintf(msg.FileWritten, filepath.Clean(opts.OutPath)),
					FlagOutPath: opts.OutPath,
					FlagFormat:  opts.Format,
				},
				Fields: fields,
			}
			describeOut.Out = f.IOStreams.Out
			return output.Print(&describeOut)
		},
	}

	cmd.Flags().StringVar(&domainID, "domain-id", "", msg.FlagDomainID)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.FlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.FlagFormat)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cmd
}

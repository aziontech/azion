package describe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/domains"

	api "github.com/aziontech/azion-cli/pkg/api/domains"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var domain_id string
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           msg.DomainsDescribeUsage,
		Short:         msg.DomainsDescribeShortDescription,
		Long:          msg.DomainsDescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli domains describe --domain-id 4312
        $ azioncli domains describe --domain-id 1337 --out "./tmp/test.json" --format json
        $ azioncli domains describe --domain-id 1337 --format json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("domain-id") {
				return msg.ErrorMissingApplicationIdArgument
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			domain, err := client.Get(ctx, domain_id)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetDomains.Error(), err)
			}

			out := f.IOStreams.Out
			formattedFuction, err := format(cmd, domain)
			if err != nil {
				return utils.ErrorFormatOut
			}

			if cmd.Flags().Changed("out") {
				err := cmdutil.WriteDetailsToFile(formattedFuction, opts.OutPath, out)
				if err != nil {
					return fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
				}
				fmt.Fprintf(out, msg.DomainsFileWritten, filepath.Clean(opts.OutPath))
			} else {
				_, err := out.Write(formattedFuction[:])
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&domain_id, "domain-id", "d", "", msg.DomainsFlagId)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.DomainsDescribeFlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.DomainsDescribeFlagFormat)
	cmd.Flags().BoolP("help", "h", false, msg.DomainsDescribeHelpFlag)

	return cmd
}

func format(cmd *cobra.Command, domain api.DomainResponse) ([]byte, error) {
	var b bytes.Buffer

	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		file, err := json.MarshalIndent(domain, "", " ")
		if err != nil {
			return nil, err
		}
		return file, nil
	} else {
		b.Write([]byte(fmt.Sprintf("ID: %d\n", uint64(domain.GetId()))))
		b.Write([]byte(fmt.Sprintf("Name: %s\n", domain.GetDomainName())))
		return b.Bytes(), nil
	}
}

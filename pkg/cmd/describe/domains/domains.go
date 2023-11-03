package domains

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/describe/domain"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/MaxwelMazur/tablecli"
	"github.com/fatih/color"

	api "github.com/aziontech/azion-cli/pkg/api/domain"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var domainID string
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           domain.Usage,
		Short:         domain.ShortDescription,
		Long:          domain.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion describe domain --domain-id 4312
        $ azion describe domain --domain-id 1337 --out "./tmp/test.json" --format json
        $ azion describe domain --domain-id 1337 --format json
        `),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !cmd.Flags().Changed("domain-id") {
				answer, err := utils.AskInput(domain.AskInputDomainID)
				if err != nil {
					return err
				}

				domainID = answer
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			domain, err := client.Get(ctx, domainID)
			if err != nil {
				return fmt.Errorf(domain.ErrorGetDomain.Error(), err.Error())
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
				logger.LogSuccess(out, fmt.Sprintf(domain.FileWritten, filepath.Clean(opts.OutPath)))
				return nil
			}

			logger.FInfo(out, string(formattedFuction[:]))
			return nil
		},
	}

	cmd.Flags().StringVar(&domainID, "domain-id", "", domain.FlagDomainID)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", domain.FlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", domain.FlagFormat)
	cmd.Flags().BoolP("help", "h", false, domain.HelpFlag)

	return cmd
}

func format(cmd *cobra.Command, domain api.DomainResponse) ([]byte, error) {
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		return json.MarshalIndent(domain, "", " ")
	}

	tbl := tablecli.New("", "")
	tbl.WithFirstColumnFormatter(color.New(color.FgGreen).SprintfFunc())

	tbl.AddRow("ID: ", domain.GetId())
	tbl.AddRow("Name: ", domain.GetName())
	tbl.AddRow("Domain: ", domain.GetDomainName())
	tbl.AddRow("Cname Access Only: ", domain.GetCnameAccessOnly())
	if domain.GetCnameAccessOnly() {
		Cnames := domain.GetCnames()
		tbl.AddRow("Cnames: ")
		for _, cname := range Cnames {
			tbl.AddRow("	", cname)
		}
	}
	tbl.AddRow("Application ID: ", domain.GetEdgeApplicationId())
	tbl.AddRow("Digital Certificate ID: ", domain.GetDigitalCertificateId())

	return tbl.GetByteFormat(), nil
}

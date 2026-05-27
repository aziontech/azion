package crl

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/crl"
	api "github.com/aziontech/azion-cli/pkg/api/crl"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var crlID int64

type DescribeCmd struct {
	Io       *iostreams.IOStreams
	AskInput func(string) (string, error)
	Get      func(context.Context, int64) (sdk.CertificateRevocationList, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		Get: func(ctx context.Context, id int64) (sdk.CertificateRevocationList, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.Get(ctx, id)
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
		$ azion describe crl --crl-id 4312
		$ azion describe crl --crl-id 1337 --out "./tmp/test.json"
		$ azion describe crl --crl-id 1337 --format json
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !cmd.Flags().Changed("crl-id") {
				answer, err := describe.AskInput(msg.AskInputCRLID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdCRL
				}
				crlID = num
			}

			ctx := context.Background()
			resp, err := describe.Get(ctx, crlID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetCRL, err.Error())
			}

			fields := map[string]string{
				"Id":             "ID",
				"Name":           "Name",
				"Active":         "Active",
				"Issuer":         "Issuer",
				"LastEditor":     "Last Editor",
				"ProductVersion": "Product Version",
				"CreatedAt":      "Created At",
				"LastModified":   "Last Modified",
				"LastUpdate":     "Last Update",
				"NextUpdate":     "Next Update",
			}

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Out:   f.IOStreams.Out,
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
				},
				Fields: fields,
				Values: &resp,
			}
			if err := output.Print(&describeOut); err != nil {
				return err
			}

			// The CRL content is multi-line PEM, so it is printed below the
			// table instead of inside it. JSON/file output already includes it.
			if len(f.Flags.Format) == 0 && len(f.Flags.Out) == 0 {
				if content := resp.GetCrl(); content != "" {
					logger.FInfo(f.IOStreams.Out, fmt.Sprintf("\nCRL:\n%s\n", content))
				}
			}

			return nil
		},
	}

	cobraCmd.Flags().Int64Var(&crlID, "crl-id", 0, msg.FlagId)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}

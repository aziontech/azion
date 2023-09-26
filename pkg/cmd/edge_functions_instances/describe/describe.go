package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	"github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/edge_functions_instances"

	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	applicationID int64
	instanceID    int64
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           msg.EdgeFuncInstanceDescribeUsage,
		Short:         msg.EdgeFuncInstanceDescribeShortDescription,
		Long:          msg.EdgeFuncInstanceDescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
      $ azion edge_functions_instances describe --application-id 1674767911 --instance-id 31223
      $ azion edge_functions_instances describe --application-id 1674767911 --instance-id 31223 --format json
      $ azion edge_functions_instances describe --application-id 1674767911 --instance-id 31223 --out "./tmp/test.json"
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("instance-id") {
				return msg.ErrorMandatoryFlags
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			ctx := context.Background()
			instance, err := client.GetFuncInstance(ctx, applicationID, instanceID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetEdgeFuncInstances.Error(), err)
			}

			out := f.IOStreams.Out
			formattedFuction, err := format(cmd, instance)
			if err != nil {
				return utils.ErrorFormatOut
			}

			if cmd.Flags().Changed("out") {
				err := cmdutil.WriteDetailsToFile(formattedFuction, opts.OutPath, out)
				if err != nil {
					return fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
				}
				fmt.Fprintf(out, msg.EdgeFuncInstanceFileWritten, filepath.Clean(opts.OutPath))
			} else {
				_, err := out.Write(formattedFuction[:])
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().Int64VarP(&applicationID, "application-id", "a", 0, msg.ApplicationFlagId)
	cmd.Flags().Int64VarP(&instanceID, "instance-id", "i", 0, msg.EdgeFuncInstanceFlagId)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.EdgeFuncInstanceDescribeFlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.EdgeFuncInstanceDescribeFlagFormat)
	cmd.Flags().BoolP("help", "h", false, msg.EdgeFuncInstanceDescribeHelpFlag)

	return cmd
}

func format(cmd *cobra.Command, instance api.FunctionsInstancesResponse) ([]byte, error) {
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		return json.MarshalIndent(instance, "", " ")
	}

	tbl := tablecli.New("", "")
	tbl.WithFirstColumnFormatter(color.New(color.FgGreen).SprintfFunc())
	tbl.AddRow("Edge Function Instance ID: ", instance.GetId())
	tbl.AddRow("Instance Name: ", instance.GetName())
	tbl.AddRow("Edge Function ID: ", instance.GetEdgeFunctionId())
	tbl.AddRow("Args: ", instance.GetArgs())
	return tbl.GetByteFormat(), nil
}

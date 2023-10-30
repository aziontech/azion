package origin

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/list/origin"
	api "github.com/aziontech/azion-cli/pkg/api/origin"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var edgeApplicationID int64 = 0

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.OriginsListUsage,
		Short:         msg.OriginsListShortDescription,
		Long:          msg.OriginsListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
        $ azion list origin  --application-id 16736354321
        $ azion list origin  --application-id 16736354321 --details
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {

				answer, err := utils.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				edgeApplicationID = num
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			if err := PrintTable(client, f, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetOrigins.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	flags := cmd.Flags()
	flags.Int64Var(&edgeApplicationID, "application-id", 0, msg.OriginsListFlagEdgeApplicationID)
	flags.BoolP("help", "h", false, msg.OriginsListHelpFlag)
	return cmd
}

func PrintTable(client *api.Client, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	c := context.Background()

	for {
		resp, err := client.ListOrigins(c, opts, edgeApplicationID)
		if err != nil {
			return err
		}

		tbl := table.New("ID", "NAME")
		tbl.WithWriter(f.IOStreams.Out)
		if opts.Details {
			tbl = table.New("ID", "NAME", "ORIGIN KEY", "ORIGIN TYPE", "ORIGIN PATH", "ADDRESSES", "CONNECTION TIMEOUT")
		}

		headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgGreen).SprintfFunc()
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, v := range resp.Results {
			tbl.AddRow(
				v.OriginId,
				utils.TruncateString(v.Name),
				v.OriginKey,
				v.OriginType,
				v.OriginPath,
				v.Addresses,
				v.ConnectionTimeout,
			)
		}

		format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
		tbl.CalculateWidths([]string{})

		// print the header only in the first flow
		if opts.Page == 1 {
			logger.PrintHeader(tbl, format)
		}

		for _, row := range tbl.GetRows() {
			logger.PrintRow(tbl, format, row)
		}

		if opts.Page >= resp.TotalPages {
			break
		}
		opts.Page++
	}

	return nil
}

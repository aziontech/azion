package list

import (
    "context"
    "fmt"
    "strings"

    "github.com/fatih/color"

    "github.com/MakeNowJust/heredoc"
    table "github.com/MaxwelMazur/tablecli"
    msg "github.com/aziontech/azion-cli/messages/origins"
    api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
    "github.com/aziontech/azion-cli/pkg/cmdutil"
    "github.com/aziontech/azion-cli/pkg/contracts"
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
        $ azioncli domains list
        $ azioncli origins list -a 16736354321
        $ azioncli origins list --application-id 16736354321
        $ azioncli origins list --application-id 16736354321 --details
        `),
        RunE: func(cmd *cobra.Command, args []string) error {
            if !cmd.Flags().Changed("application-id") {
                return msg.ErrorMissingApplicationIDArgument
            } 
            
            if err := PrintTable(cmd, f, opts); err != nil {
                return fmt.Errorf(msg.ErrorGetOrigins.Error(), err)
            }
            return nil
        },
    }

    cmdutil.AddAzionApiFlags(cmd, opts)
    flags := cmd.Flags()
    flags.Int64VarP(&edgeApplicationID, "application-id", "a", 0, msg.OriginsListFlagEdgeApplicationID)
    flags.BoolP("help", "h", false, msg.OriginsListHelpFlag)
    return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions) error {
    client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
    ctx := context.Background()

    response, err := client.ListOrigins(ctx, opts, edgeApplicationID)
    if err != nil {
        return fmt.Errorf(msg.ErrorGetOrigins.Error(), err)
    }
    
    tbl := table.New("ID", "NAME")
    tbl.WithWriter(f.IOStreams.Out)
    if cmd.Flags().Changed("details") {
        tbl = table.New("ID", "NAME", "ORIGIN TYPE", "ORIGIN PATH", "ADDRESSES", "CONNECTION TIMEOUT")
    }

    headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
    columnFmt := color.New(color.FgGreen).SprintfFunc()
    tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
    if cmd.Flags().Changed("details") {
        for _, v := range response.Results {
            tbl.AddRow(v.OriginId, v.Name, v.OriginType, v.OriginPath, v.Addresses, v.ConnectionTimeout)
        }
    } else {
        for _, v := range response.Results {
            tbl.AddRow(v.OriginId, v.Name)
        }
    }

    format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
    tbl.CalculateWidths([]string{})
    tbl.PrintHeader(format)
    for _, row := range tbl.GetRows() {
        tbl.PrintRow(format, row)
    }  
    return nil
}

package edge_storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/briandowns/spinner"
	"github.com/nsf/termbox-go"

	table "github.com/MaxwelMazur/tablecli"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	"github.com/aziontech/azion-cli/messages/general"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
)

func NewObject(f *cmdutil.Factory) *cobra.Command {
	object := &Objects{
		Factory: f,
		Options: &contracts.ListOptions{},
	}
	cmd := &cobra.Command{
		Use:           msg.USAGE_OBJECTS,
		Short:         msg.SHORT_DESCRIPTION_LIST_OBJECT,
		Long:          msg.LONG_DESCRIPTION_LIST_OBJECT,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          object.RunE,
		Example:       heredoc.Doc(msg.EXAMPLE_LIST_OBJECT),
	}
	object.AddFlags(cmd.Flags())
	return cmd
}

func (b *Objects) RunE(cmd *cobra.Command, args []string) error {
	if !cmd.Flags().Changed("bucket-name") {
		answer, err := utils.AskInput(msg.ASK_NAME_CREATE_BUCKET)
		if err != nil {
			return err
		}
		b.BucketName = answer
	}
	client := api.NewClient(
		b.Factory.HttpClient,
		b.Factory.Config.GetString("storage_url"),
		b.Factory.Config.GetString("token"))
	return b.PrintTable(client)
}

func (b *Objects) PrintTable(client *api.Client) error {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	printHeader := true
	count := 0
	for {
		c := context.Background()

		settings, err := token.ReadSettings()
		if err != nil {
			return err
		}

		if count > 0 && len(settings.ContinuationToken) == 0 {
			return nil
		}

		b.Options.ContinuationToken = settings.ContinuationToken
		count = count + 1

		resp, err := client.ListObject(c, b.BucketName, b.Options)
		if err != nil {
			return fmt.Errorf(msg.ERROR_LIST_BUCKET, err)
		}

		settings.ContinuationToken = resp.GetContinuationToken()
		err = token.WriteSettings(settings)
		if err != nil {
			return err
		}

		tbl := table.New("KEY", "LAST MODIFIED")
		tbl.WithWriter(b.Factory.IOStreams.Out)
		if b.Options.Details {
			tbl = table.New("KEY", "LAST MODIFIED", "SIZE", "ETAG")
		}
		headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgGreen).SprintfFunc()
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		for _, v := range resp.Results {
			tbl.AddRow(v.GetKey(), v.GetLastModified(), v.GetSize(), v.GetEtag())
		}
		format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
		tbl.CalculateWidths([]string{})
		if printHeader {
			logger.PrintHeader(tbl, format)
			printHeader = false
		}
		for _, row := range tbl.GetRows() {
			logger.PrintRow(tbl, format, row)
		}

		s := spinner.New(spinner.CharSets[26], 150*time.Millisecond)
		s.Prefix = "Press 'q' to exit, Enter or Space to continue"
		s.Start()

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc || ev.Ch == 'q' {
				s.Stop()
				return nil
			}

			if ev.Key == termbox.KeySpace || ev.Key == termbox.KeyEnter {
				s.Stop()
				continue
			}
		}
	}
}

func (b *Objects) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&b.BucketName, "bucket-name", "", msg.FLAG_NAME_BUCKET)
	flags.BoolVar(&b.Options.Details, "details", false, msg.FLAG_HELP_DETAILS_OBJECTS)
	flags.Int64Var(&b.Options.PageSize, "page-size", 50, general.ApiListFlagPageSize)
	flags.BoolP("help", "h", false, msg.FLAG_HELP_LIST_OBJECT)
}

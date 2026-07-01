package streams

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/streams"
	api "github.com/aziontech/azion-cli/pkg/api/data_stream"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/data-stream-api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var streamID int64

type DescribeCmd struct {
	Io       *iostreams.IOStreams
	AskInput func(string) (string, error)
	Get      func(context.Context, int64) (sdk.DataStream, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		Get: func(ctx context.Context, streamID int64) (sdk.DataStream, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.Get(ctx, streamID)
		},
	}
}

func NewCobraCmd(describe *DescribeCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cobraCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.DescribeShortDescription,
		Long:          msg.DescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion describe streams --stream-id 4312
        $ azion describe streams --stream-id 1337 --out "./tmp/test.json" --format json
        $ azion describe streams --stream-id 1337 --format json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("stream-id") {
				answer, err := describe.AskInput(msg.AskStreamID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertStreamId
				}

				streamID = num
			}

			ctx := context.Background()
			resp, err := describe.Get(ctx, streamID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetStream.Error(), err)
			}

			fields := map[string]string{
				"Id":             "ID",
				"Name":           "Name",
				"Active":         "Active",
				"ProductVersion": "Product Version",
				"LastModified":   "Modified at",
				"LastEditor":     "Last Editor",
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

			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().Int64Var(&streamID, "stream-id", 0, msg.FlagID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}

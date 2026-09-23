package applications

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/applications"
	api "github.com/aziontech/azion-cli/pkg/api/applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DescribeCmd struct {
	Io            *iostreams.IOStreams
	AskInput      func(string) (string, error)
	Get           func(context.Context, int64) (api.ApplicationResponse, error)
	ApplicationID int64
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		Get: func(ctx context.Context, applicationID int64) (api.ApplicationResponse, error) {
			client := registry.Applications(f)
			return client.Get(ctx, applicationID)
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
		$ azion describe application --application-id 4312
		$ azion describe application --application-id 1337 --out "./tmp/test.json"
		$ azion describe application --application-id 1337 --format json
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := describe.AskInput(msg.AskInputApplicationID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}
				describe.ApplicationID = num
			}

			ctx := context.Background()
			resp, err := describe.Get(ctx, describe.ApplicationID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetApplication.Error(), err)
			}

			fields := map[string]string{
				"Id":     "ID",
				"Name":   "Name",
				"Active": "Active",
				"Debug":  "Debug",
			}

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Out:   f.IOStreams.Out,
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
				},
				Fields: fields,
				Values: resp,
			}
			return output.Print(&describeOut)

		},
	}

	cobraCmd.Flags().Int64Var(&describe.ApplicationID, "application-id", 0, msg.FlagId)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}

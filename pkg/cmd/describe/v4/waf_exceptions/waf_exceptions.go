package wafexceptions

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/waf_exceptions"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DescribeCmd struct {
	Io              *iostreams.IOStreams
	AskInput        func(string) (string, error)
	GetWafException func(ctx context.Context, wafId, exceptionId int64) (sdk.WAFRule, error)
	ExceptionID     int64
	WafID           int64
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		GetWafException: func(ctx context.Context, wafId, exceptionId int64) (sdk.WAFRule, error) {
			client := registry.WAFExceptions(f)
			return client.Get(ctx, wafId, exceptionId)
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
		$ azion describe waf-exceptions --waf-id 4312 --exception-id 42069
		$ azion describe waf-exceptions --waf-id 1337 --exception-id 42069 --out "./wafexception.json" --format json
		$ azion describe waf-exceptions --waf-id 1337 --exception-id 42069 --format json
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !cmd.Flags().Changed("waf-id") {
				answer, err := describe.AskInput(msg.AskInputWafID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertWafID
				}

				describe.WafID = num
			}

			if !cmd.Flags().Changed("exception-id") {
				answer, err := describe.AskInput(msg.AskInputExceptionID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertExceptionID
				}

				describe.ExceptionID = num
			}

			ctx := context.Background()
			exception, err := describe.GetWafException(ctx, describe.WafID, describe.ExceptionID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetWafException, err.Error())
			}

			fields := make(map[string]string)
			fields["Id"] = "ID"
			fields["Name"] = "Name"
			fields["RuleId"] = "Rule ID"
			fields["Path"] = "Path"
			fields["Operator"] = "Operator"
			fields["Active"] = "Active"
			fields["LastEditor"] = "Last Editor"
			fields["LastModified"] = "Last Modified"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
					Out:   f.IOStreams.Out,
				},
				Fields: fields,
				Values: &exception,
			}
			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().Int64Var(&describe.WafID, "waf-id", 0, msg.FlagWafID)
	cobraCmd.Flags().Int64Var(&describe.ExceptionID, "exception-id", 0, msg.FlagExceptionID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}

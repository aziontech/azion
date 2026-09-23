package workloaddeployment

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/workload_deployment"
	api "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type DescribeCmd struct {
	Io            *iostreams.IOStreams
	AskInput      func(string) (string, error)
	GetDeployment func(ctx context.Context, id, deploymentid int64) (api.DeploymentResponse, error)
	DeploymentID  int64
	WorkloadID    int64
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		GetDeployment: func(ctx context.Context, id, deploymentid int64) (api.DeploymentResponse, error) {
			client := registry.Workloads(f)
			return client.GetDeployment(ctx, id, deploymentid)
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
		$ azion describe workload-deployment --workload-id 4312 --deployment-id 666
		$ azion describe workload-deployment --workload-id 1337 --deployment-id 42069 --out "./tmp/test.json" --format json
		$ azion describe workload-deployment --workload-id 1337 --deployment-id 6669 --format json
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !cmd.Flags().Changed("workload-id") {
				answer, err := describe.AskInput(msg.AskInputWorkloadID)
				if err != nil {
					return err
				}

				describe.WorkloadID, err = strconv.ParseInt(answer, 10, 64)
				if err != nil {
					return msg.ErrorConvertWorkloadId
				}
			}

			if !cmd.Flags().Changed("deployment-id") {
				answer, err := describe.AskInput(msg.AskInputDeploymentID)
				if err != nil {
					return err
				}

				describe.DeploymentID, err = strconv.ParseInt(answer, 10, 64)
				if err != nil {
					return msg.ErrorConvertDeploymentId
				}
			}

			ctx := context.Background()
			workload, err := describe.GetDeployment(ctx, describe.WorkloadID, describe.DeploymentID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetDeployment.Error(), err.Error())
			}

			fields := make(map[string]string)
			fields["Id"] = "ID"
			fields["Tag"] = "Tag"
			fields["Current"] = "Current"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
					Out:   f.IOStreams.Out,
				},
				Fields: fields,
				Values: workload,
			}
			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().Int64Var(&describe.WorkloadID, "workload-id", 0, msg.FlagWorkloadID)
	cobraCmd.Flags().Int64Var(&describe.DeploymentID, "deployment-id", 0, msg.FlagDeploymentID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}

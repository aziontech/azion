package kv

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/create/kv"
	api "github.com/aziontech/azion-cli/pkg/api/kv"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type Fields struct {
	Namespace string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create kv --namespace "my-namespace"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.CreateRequest{}

			if !cmd.Flags().Changed("namespace") {
				answer, err := utils.AskInput(msg.AskNamespace)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}
				fields.Namespace = answer
			}

			request.SetName(fields.Namespace)

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()
			_, err := client.Create(ctx, request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateNamespace, err.Error())
			}

			creatOut := output.GeneralOutput{
				Msg: fmt.Sprintf(msg.CreateOutputSuccess, fields.Namespace),
				Out: f.IOStreams.Out,
			}
			return output.Print(&creatOut)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)

	return cmd
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.Namespace, "namespace", "", msg.FlagNamespace)
	flags.BoolP("help", "h", false, msg.HelpFlag)
}

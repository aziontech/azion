package edgeconnector

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_connector"
	api "github.com/aziontech/azion-cli/pkg/api/edge_connector"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type Fields struct {
	Name                 string
	Type                 string
	LoadBalancer         string
	OriginShield         string
	Active               string
	Addresses            []string
	TlsPolicy            string
	LoadBalanceMethod    string
	ConnectionPreference []string
	ConnectionTimeout    int64
	ReadWriteTimeout     int64
	MaxRetries           int64
	Versions             []string
	Host                 string
	Path                 string
	FollowingRedirect    string
	RealIpHeader         string
	RealPortHeader       string
	Bucket               string
	Prefix               string
	InPath               string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create edge-connector --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.NewCreateRequest()

			if !cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &request)
				if err != nil {
					logger.Debug("Failed to unmarshal file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Create(ctx, request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateConnector.Error(), err)
			}

			creatOut := output.GeneralOutput{
				Msg: fmt.Sprintf(msg.CreateOutputSuccess, response.GetId()),
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
	flags.StringVar(&fields.InPath, "file", "", msg.FlagIn)
	flags.BoolP("help", "h", false, msg.CreateFlagHelp)
}

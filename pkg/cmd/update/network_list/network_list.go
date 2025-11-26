package networklist

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/network_list"
	api "github.com/aziontech/azion-cli/pkg/api/network_list"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type Fields struct {
	ID     string
	Name   string
	Type   string
	Items  string
	Active string
	InPath string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.UpdateShortDescription,
		Long:          msg.UpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion update network-list --network-list-id 1234 --name "Updated List"
		$ azion update network-list --network-list-id 4185 --type ip_cidr --items "192.168.1.0/24,10.0.0.0/8"
		$ azion update network-list --network-list-id 9123 --active true
		$ azion update network-list --network-list-id 9123 --active false
		$ azion update network-list --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("network-list-id") {
				answer, err := utils.AskInput(msg.UpdateAskNetworkListID)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				fields.ID = answer
			}

			request := api.UpdateRequest{}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {
				err := createRequestFromFlags(cmd, fields, &request)
				if err != nil {
					return err
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, &request, fields.ID)

			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateNetworkList.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.UpdateOutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&updateOut)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)

	return cmd
}

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.UpdateRequest) error {
	if cmd.Flags().Changed("name") {
		request.SetName(fields.Name)
	}

	if cmd.Flags().Changed("type") {
		request.SetType(fields.Type)
	}

	if cmd.Flags().Changed("items") {
		items := strings.Split(fields.Items, ",")
		for i := range items {
			items[i] = strings.TrimSpace(items[i])
		}
		request.SetItems(items)
	}

	if cmd.Flags().Changed("active") {
		active, err := strconv.ParseBool(fields.Active)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorActiveFlag, fields.Active)
		}
		request.SetActive(active)
	}

	return nil
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.ID, "network-list-id", "", msg.FlagID)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.Type, "type", "", msg.FlagType)
	flags.StringVar(&fields.Items, "items", "", msg.FlagItems)
	flags.StringVar(&fields.Active, "active", "", msg.FlagActive)
	flags.StringVar(&fields.InPath, "file", "", msg.FlagIn)
	flags.BoolP("help", "h", false, msg.UpdateHelpFlag)
}

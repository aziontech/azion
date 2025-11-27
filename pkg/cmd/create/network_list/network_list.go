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
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create network-list --name "My IP List" --type ip_cidr --items "192.168.0.1/32,10.0.0.0/8"
        $ azion create network-list --name "ASN List" --type asn --items "1234,5678" --active true
        $ azion create network-list --file "./network_list.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.NewCreateRequest()

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &request)
				if err != nil {
					logger.Debug("Failed to unmarshal file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				err := createRequestFromFlags(cmd, fields, request)
				if err != nil {
					return err
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Create(ctx, request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateNetworkList.Error(), err)
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
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.Type, "type", "", msg.FlagType)
	flags.StringVar(&fields.Items, "items", "", msg.FlagItems)
	flags.StringVar(&fields.Active, "active", "", msg.FlagActive)
	flags.StringVar(&fields.InPath, "file", "", msg.FlagIn)
	flags.BoolP("help", "h", false, msg.CreateHelpFlag)
}

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.CreateRequest) error {

	if !cmd.Flags().Changed("name") {
		answers, err := utils.AskInput(msg.AskName)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		fields.Name = answers
	}

	if !cmd.Flags().Changed("type") {
		answers, err := utils.Select(utils.NewSelectPrompter(msg.AskType, []string{"asn", "countries", "ip_cidr"}))
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		fields.Type = answers
	}

	if !cmd.Flags().Changed("items") {
		answers, err := utils.AskInput(msg.AskItems)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		fields.Items = answers
	}

	request.SetName(fields.Name)
	request.SetType(fields.Type)

	items := strings.Split(fields.Items, ",")
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
	}
	request.SetItems(items)

	if cmd.Flags().Changed("active") {
		isActive, err := strconv.ParseBool(fields.Active)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorActiveFlag, fields.Active)
		}
		request.SetActive(isActive)
	}

	return nil
}

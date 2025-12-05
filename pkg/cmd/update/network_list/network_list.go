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
	ID         int64
	Name       string
	Type       string
	Items      string
	AddItem    string
	RemoveItem string
	Active     string
	InPath     string
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
		$ azion update network-list --network-list-id 1 --add-item "1.1.1.1"
		$ azion update network-list --network-list-id 1 --add-item "1.1.1.1,2.2.2.2,3.3.3.3"
		$ azion update network-list --network-list-id 1 --remove-item "1.1.1.1"
		$ azion update network-list --network-list-id 1 --remove-item "1.1.1.1,2.2.2.2"
		$ azion update network-list --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("network-list-id") {
				answer, err := utils.AskInput(msg.UpdateAskNetworkListID)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertNetworkListId
				}

				fields.ID = num
			}

			request := api.UpdateRequest{}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			ctx := context.Background()

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {
				if cmd.Flags().Changed("add-item") || cmd.Flags().Changed("remove-item") {
					current, err := client.Get(ctx, fields.ID)
					if err != nil {
						return fmt.Errorf(msg.ErrorGetNetworkList.Error(), err)
					}

					currentItems := current.GetItems()
					modifiedItems := make([]string, len(currentItems))
					copy(modifiedItems, currentItems)

					if cmd.Flags().Changed("add-item") {
						for newItem := range strings.SplitSeq(fields.AddItem, ",") {
							newItem = strings.TrimSpace(newItem)
							if newItem == "" {
								continue
							}

							exists := false
							for _, item := range modifiedItems {
								if item == newItem {
									exists = true
									break
								}
							}
							if !exists {
								modifiedItems = append(modifiedItems, newItem)
							}
						}
					}

					if cmd.Flags().Changed("remove-item") {
						removeMap := make(map[string]bool)
						for removeItem := range strings.SplitSeq(fields.RemoveItem, ",") {
							removeItem = strings.TrimSpace(removeItem)
							if removeItem != "" {
								removeMap[removeItem] = true
							}
						}
						filteredItems := []string{}
						for _, item := range modifiedItems {
							if !removeMap[item] {
								filteredItems = append(filteredItems, item)
							}
						}
						modifiedItems = filteredItems
					}

					request.SetItems(modifiedItems)
				}

				err := createRequestFromFlags(cmd, fields, &request)
				if err != nil {
					return err
				}
			}

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
	flags.Int64Var(&fields.ID, "network-list-id", 0, msg.FlagID)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.Type, "type", "", msg.FlagType)
	flags.StringVar(&fields.Items, "items", "", msg.FlagItems)
	flags.StringVar(&fields.AddItem, "add-item", "", msg.FlagAddItem)
	flags.StringVar(&fields.RemoveItem, "remove-item", "", msg.FlagRemoveItem)
	flags.StringVar(&fields.Active, "active", "", msg.FlagActive)
	flags.StringVar(&fields.InPath, "file", "", msg.FlagIn)
	flags.BoolP("help", "h", false, msg.UpdateHelpFlag)
}

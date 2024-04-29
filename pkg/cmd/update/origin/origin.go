package origin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/origin"
	api "github.com/aziontech/azion-cli/pkg/api/origin"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Fields struct {
	OriginKey            string
	ApplicationID        int64
	Name                 string
	OriginType           string
	Addresses            []string
	OriginProtocolPolicy string
	HostHeader           string
	OriginPath           string
	HmacAuthentication   string
	HmacRegionName       string
	HmacAccessKey        string
	HmacSecretKey        string
	Path                 string
	Bucket               string
	Prefix               string
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
        $ azion update origin --application-id 1673635839 --origin-key "58755fef-e830-4ea4-b9e0-6481f1ef496d" --name "ffcafe222sdsdffdf" --addresses "httpbin.org" --host-header "\${host}" --origin-type "single_origin" --origin-protocol-policy "http" --origin-path "/requests" --hmac-authentication "false"
        $ azion update origin --application-id 1673635839 --origin-key "58755fef-e830-4ea4-b9e0-6481f1ef496d" --name "drink coffe" --addresses "asdfg.asd" --host-header "\${host}"
        $ azion update origin --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.UpdateRequest{}

			if !cmd.Flags().Changed("application-id") {
				answers, err := utils.AskInput(msg.AskAppID)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				applicationID, err := strconv.Atoi(answers)
				if err != nil {
					logger.Debug("Error while parsing string to integer", zap.Error(err))
					return utils.ErrorConvertingStringToInt
				}

				fields.ApplicationID = int64(applicationID)
			}

			if !cmd.Flags().Changed("origin-key") {
				answers, err := utils.AskInput(msg.AskOriginKey)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				fields.OriginKey = answers
			}

			if cmd.Flags().Changed("file") {
				if err := utils.FlagFileUnmarshalJSON(fields.Path, &request); err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {
				err := createRequestFromFlags(cmd, fields, &request)
				if err != nil {
					return err
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.Update(context.Background(), fields.ApplicationID, fields.OriginKey, &request)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateOrigin.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg: fmt.Sprintf(msg.UpdateOutputSuccess, response.GetOriginKey()),
				Out: f.IOStreams.Out,
			}
			return output.Print(&updateOut)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)
	return cmd
}

func prepareAddresses(addrs []string) (addresses []sdk.CreateOriginsRequestAddresses) {
	var addr sdk.CreateOriginsRequestAddresses
	for _, v := range addrs {
		addr.Address = v
		addresses = append(addresses, addr)
	}
	return
}

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.UpdateRequest) error {
	if cmd.Flags().Changed("name") {
		request.SetName(fields.Name)
	}
	if cmd.Flags().Changed("addresses") {
		request.SetAddresses(prepareAddresses(fields.Addresses))
	}
	if cmd.Flags().Changed("host-header") {
		request.SetHostHeader(fields.HostHeader)
	}
	if cmd.Flags().Changed("origin-type") {
		request.SetOriginType(fields.OriginType)
	}
	if cmd.Flags().Changed("origin-protocol-policy") {
		request.SetOriginProtocolPolicy(fields.OriginProtocolPolicy)
	}
	if cmd.Flags().Changed("origin-path") {
		request.SetOriginPath(fields.OriginPath)
	}
	if cmd.Flags().Changed("bucket") {
		request.SetBucket(fields.Bucket)
	}
	if cmd.Flags().Changed("prefix") {
		request.SetPrefix(fields.Prefix)
	}

	if cmd.Flags().Changed("hmac-authentication") {
		hmacAuth, err := strconv.ParseBool(fields.HmacAuthentication)
		if err != nil {
			return msg.ErrorHmacAuthenticationFlag
		}
		request.SetHmacAuthentication(hmacAuth)
	}

	if cmd.Flags().Changed("hmac-region-name") {
		request.SetHmacRegionName(fields.HmacRegionName)
	}

	if cmd.Flags().Changed("hmac-access-key") {
		request.SetHmacAccessKey(fields.HmacAccessKey)
	}

	if cmd.Flags().Changed("hmac-secret-key") {
		request.SetHmacSecretKey(fields.HmacSecretKey)
	}

	return nil
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.OriginKey, "origin-key", "", msg.FlagOriginKey)
	flags.Int64Var(&fields.ApplicationID, "application-id", 0, msg.FlagEdgeApplicationID)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.OriginType, "origin-type", "", msg.FlagOriginType)
	flags.StringSliceVar(&fields.Addresses, "addresses", []string{}, msg.FlagAddresses)
	flags.StringVar(&fields.OriginProtocolPolicy, "origin-protocol-policy", "", msg.FlagOriginProtocolPolicy)
	flags.StringVar(&fields.HostHeader, "host-header", "", msg.FlagHostHeader)
	flags.StringVar(&fields.OriginPath, "origin-path", "", msg.FlagOriginPath)
	flags.StringVar(&fields.HmacAuthentication, "hmac-authentication", "", msg.FlagHmacAuthentication)
	flags.StringVar(&fields.HmacRegionName, "hmac-region-name", "", msg.FlagHmacRegionName)
	flags.StringVar(&fields.HmacAccessKey, "hmac-access-key", "", msg.FlagHmacAccessKey)
	flags.StringVar(&fields.HmacSecretKey, "hmac-secret-key", "", msg.FlagHmacSecretKey)
	flags.StringVar(&fields.Bucket, "bucket", "", msg.FlagBucketUpdate)
	flags.StringVar(&fields.Prefix, "prefix", "", msg.FlagPrefixUpdate)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.UpdateFlagHelp)
}

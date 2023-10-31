package origin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	api "github.com/aziontech/azion-cli/pkg/api/origins"
	"github.com/aziontech/azion-cli/pkg/logger"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var example = heredoc.Doc(`
	$ azion create origin --application-id 1673635839 --name "drink coffe" --addresses "asdfg.asd" --host-header "host"
	$ azion create origin --application-id 1673635839 --in "create.json"
	`)

type Fields struct {
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
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           usage,
		Short:         shortDescription,
		Long:          longDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       example,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.Request{}

			if cmd.Flags().Changed("in") {
				err := utils.FlagINUnmarshalFileJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				err := createRequestFromFlags(cmd, fields, &request)
				if err != nil {
					return err
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.Create(context.Background(), fields.ApplicationID, &request)
			if err != nil {
				return fmt.Errorf(errorCreateOrigins, err)
			}
			logger.LogSuccess(f.IOStreams.Out, fmt.Sprintf(outputSuccess, response.GetOriginId()))

			return nil
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

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.Request) error {
	if !cmd.Flags().Changed("application-id") {
		answers, err := utils.AskInput("What is the ID of the Edge Application?")
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

	if !cmd.Flags().Changed("name") {
		answers, err := utils.AskInput("What is the Name of the Origins?")
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}

		fields.Name = answers
	}

	if !cmd.Flags().Changed("addresses") {
		answers, err := utils.AskInput("What is the Addresses of the Origins?")
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}

		fields.Addresses = []string{answers}
	}

	if !cmd.Flags().Changed("host-header") {
		answers, err := utils.AskInput("What is the Host Header of the Origins?")
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}

		fields.HostHeader = answers
	}

	request.SetName(fields.Name)
	request.SetAddresses(prepareAddresses(fields.Addresses))
	request.SetHostHeader(fields.HostHeader)

	if cmd.Flags().Changed("origin-type") {
		request.SetOriginType(fields.OriginType)
	}

	if cmd.Flags().Changed("origin-protocol-policy") {
		request.SetOriginProtocolPolicy(fields.OriginProtocolPolicy)
	}

	if cmd.Flags().Changed("origin-path") {
		request.SetOriginPath(fields.OriginPath)
	}

	if cmd.Flags().Changed("hmac-authentication") {
		hmacAuth, err := strconv.ParseBool(fields.HmacAuthentication)
		if err != nil {
			return fmt.Errorf(errorHmacAuthenticationFlag)
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
	flags.Int64Var(&fields.ApplicationID, "application-id", 0, flagEdgeApplicationID)
	flags.StringVar(&fields.Name, "name", "", flagName)
	flags.StringVar(&fields.OriginType, "origin-type", "", flagOriginType)
	flags.StringSliceVar(&fields.Addresses, "addresses", []string{}, flagAddresses)
	flags.StringVar(&fields.OriginProtocolPolicy, "origin-protocol-policy", "", flagOriginProtocolPolicy)
	flags.StringVar(&fields.HostHeader, "host-header", "", flagHostHeader)
	flags.StringVar(&fields.OriginPath, "origin-path", "", flagOriginPath)
	flags.StringVar(&fields.HmacAuthentication, "hmac-authentication", "", flagHmacAuthentication)
	flags.StringVar(&fields.HmacRegionName, "hmac-region-name", "", flagHmacRegionName)
	flags.StringVar(&fields.HmacAccessKey, "hmac-access-key", "", flagHmacAccessKey)
	flags.StringVar(&fields.HmacSecretKey, "hmac-secret-key", "", flagHmacSecretKey)
	flags.StringVar(&fields.Path, "in", "", flagIn)
	flags.BoolP("help", "h", false, flagHelp)
}

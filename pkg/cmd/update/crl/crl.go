package crl

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/update/crl"
	api "github.com/aziontech/azion-cli/pkg/api/crl"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type Fields struct {
	ID     int64
	Name   string
	Issuer string
	CRL    string
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
		$ azion update crl --crl-id 1234 --name 'My CRL'
		$ azion update crl --crl-id 1234 --crl ./list.crl
		$ azion update crl --crl-id 1234 --active false
		$ azion update crl --crl-id 1234 --file "update.json"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("crl-id") {
				answer, err := utils.AskInput(msg.UpdateAskCRLID)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdCRL
				}

				fields.ID = num
			}

			request := api.UpdateRequest{}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {
				if err := updateRequestFromFlags(cmd, fields, &request); err != nil {
					return err
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, &request, fields.ID)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateCRL.Error(), err)
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

func updateRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.UpdateRequest) error {
	if cmd.Flags().Changed("name") {
		request.SetName(fields.Name)
	}

	if cmd.Flags().Changed("issuer") {
		request.SetIssuer(fields.Issuer)
	}

	if cmd.Flags().Changed("crl") {
		content, err := os.ReadFile(fields.CRL)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorReadCRLFile, fields.CRL)
		}
		request.SetCrl(string(content))
	}

	if cmd.Flags().Changed("active") {
		isActive, err := strconv.ParseBool(fields.Active)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorActiveFlag, fields.Active)
		}
		request.SetActive(isActive)
	}

	return nil
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.Int64Var(&fields.ID, "crl-id", 0, msg.FlagID)
	flags.StringVar(&fields.Name, "name", "", msg.UpdateFlagName)
	flags.StringVar(&fields.Issuer, "issuer", "", msg.UpdateFlagIssuer)
	flags.StringVar(&fields.CRL, "crl", "", msg.UpdateFlagCRL)
	flags.StringVar(&fields.Active, "active", "", msg.UpdateFlagActive)
	flags.StringVar(&fields.InPath, "file", "", msg.UpdateFlagFile)
	flags.BoolP("help", "h", false, msg.UpdateHelpFlag)
}

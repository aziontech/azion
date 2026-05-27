package digitalcertificate

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/update/digital_certificate"
	api "github.com/aziontech/azion-cli/pkg/api/digital_certificate"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type Fields struct {
	ID              int64
	Name            string
	Active          string
	Certificate     string
	PrivateKey      string
	CertificateType string
	Authority       string
	Challenge       string
	KeyAlgorithm    string
	InPath          string
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
		$ azion update digital-certificate --digital-certificate-id 1234 --name 'My Certificate'
		$ azion update digital-certificate --digital-certificate-id 1234 --certificate ./cert.pem --private-key ./key.pem
		$ azion update digital-certificate --digital-certificate-id 1234 --active false
		$ azion update digital-certificate --digital-certificate-id 1234 --file "update.json"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("digital-certificate-id") {
				answer, err := utils.AskInput(msg.UpdateAskDigitalCertificateID)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdDigitalCertificate
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
				err := updateRequestFromFlags(cmd, fields, &request)
				if err != nil {
					return err
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, &request, fields.ID)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateDigitalCertificate.Error(), err)
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

	if cmd.Flags().Changed("active") {
		isActive, err := strconv.ParseBool(fields.Active)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorActiveFlag, fields.Active)
		}
		request.SetActive(isActive)
	}

	if cmd.Flags().Changed("certificate") {
		content, err := os.ReadFile(fields.Certificate)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorReadCertificateFile, fields.Certificate)
		}
		request.SetCertificate(string(content))
	}

	if cmd.Flags().Changed("private-key") {
		content, err := os.ReadFile(fields.PrivateKey)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorReadPrivateKeyFile, fields.PrivateKey)
		}
		request.SetPrivateKey(string(content))
	}

	if cmd.Flags().Changed("certificate-type") {
		request.SetType(fields.CertificateType)
	}

	if cmd.Flags().Changed("authority") {
		request.SetAuthority(fields.Authority)
	}

	if cmd.Flags().Changed("challenge") {
		request.SetChallenge(fields.Challenge)
	}

	if cmd.Flags().Changed("key-algorithm") {
		request.SetKeyAlgorithm(fields.KeyAlgorithm)
	}

	return nil
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.Int64Var(&fields.ID, "digital-certificate-id", 0, msg.FlagID)
	flags.StringVar(&fields.Name, "name", "", msg.UpdateFlagName)
	flags.StringVar(&fields.Active, "active", "", msg.UpdateFlagActive)
	flags.StringVar(&fields.Certificate, "certificate", "", msg.UpdateFlagCertificate)
	flags.StringVar(&fields.PrivateKey, "private-key", "", msg.UpdateFlagPrivateKey)
	flags.StringVar(&fields.CertificateType, "certificate-type", "", msg.UpdateFlagCertificateType)
	flags.StringVar(&fields.Authority, "authority", "", msg.UpdateFlagAuthority)
	flags.StringVar(&fields.Challenge, "challenge", "", msg.UpdateFlagChallenge)
	flags.StringVar(&fields.KeyAlgorithm, "key-algorithm", "", msg.UpdateFlagKeyAlgorithm)
	flags.StringVar(&fields.InPath, "file", "", msg.UpdateFlagFile)
	flags.BoolP("help", "h", false, msg.UpdateHelpFlag)
}

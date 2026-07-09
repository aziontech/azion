package digitalcertificate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/create/digital_certificate"
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
	Name             string
	Certificate      string
	PrivateKey       string
	CertificateType  string
	Authority        string
	Challenge        string
	CommonName       string
	AlternativeNames string
	KeyAlgorithm     string
	InPath           string
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
        $ azion create digital-certificate --name "My Certificate" --certificate "./cert.pem" --private-key "./key.pem"
        $ azion create digital-certificate --name "azion.com" --authority lets_encrypt --challenge dns --common-name "azion.com"
        $ azion create digital-certificate --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			ctx := context.Background()

			if cmd.Flags().Changed("file") {
				return runFromFile(ctx, f, client, fields.InPath)
			}

			if cmd.Flags().Changed("authority") {
				return runRequest(ctx, cmd, f, client, fields)
			}

			return runCreate(ctx, cmd, f, client, fields)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)

	return cmd
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.Certificate, "certificate", "", msg.FlagCertificate)
	flags.StringVar(&fields.PrivateKey, "private-key", "", msg.FlagPrivateKey)
	flags.StringVar(&fields.CertificateType, "certificate-type", "", msg.FlagCertificateType)
	flags.StringVar(&fields.Authority, "authority", "", msg.FlagAuthority)
	flags.StringVar(&fields.Challenge, "challenge", "", msg.FlagChallenge)
	flags.StringVar(&fields.CommonName, "common-name", "", msg.FlagCommonName)
	flags.StringVar(&fields.AlternativeNames, "alternative-names", "", msg.FlagAlternativeNames)
	flags.StringVar(&fields.KeyAlgorithm, "key-algorithm", "", msg.FlagKeyAlgorithm)
	flags.StringVar(&fields.InPath, "file", "", msg.FlagIn)
	flags.BoolP("help", "h", false, msg.CreateFlagHelp)
}

type createInput struct {
	Name        string  `json:"name"`
	Certificate *string `json:"certificate,omitempty"`
	PrivateKey  *string `json:"private_key,omitempty"`
	Type        *string `json:"type,omitempty"`
}

func runFromFile(ctx context.Context, f *cmdutil.Factory, client *api.Client, path string) error {
	data, err := readInput(path)
	if err != nil {
		return err
	}

	var input createInput
	if err := strictUnmarshal(data, &input); err == nil {
		createReq := api.NewCreateRequest()
		createReq.SetName(input.Name)
		if input.Certificate != nil {
			createReq.SetCertificate(*input.Certificate)
		}
		if input.PrivateKey != nil {
			createReq.SetPrivateKey(*input.PrivateKey)
		}
		if input.Type != nil {
			createReq.SetType(*input.Type)
		}
		return doCreate(ctx, f, client, createReq)
	} else {
		logger.Debug("Input file does not match Digital Certificate schema, trying Certificate Request", zap.Error(err))
	}

	requestReq := api.NewRequestRequest()
	if err := strictUnmarshal(data, &requestReq.CertificateRequest); err == nil {
		return doRequest(ctx, f, client, requestReq)
	} else {
		logger.Debug("Input file does not match Certificate Request schema either", zap.Error(err))
	}

	return msg.ErrorInvalidJSONFile
}

func runCreate(ctx context.Context, cmd *cobra.Command, f *cmdutil.Factory, client *api.Client, fields *Fields) error {
	req := api.NewCreateRequest()

	if !cmd.Flags().Changed("name") {
		answer, err := utils.AskInput(msg.AskName)
		if err != nil {
			return err
		}
		fields.Name = answer
	}
	req.SetName(fields.Name)

	if cmd.Flags().Changed("certificate") {
		content, err := os.ReadFile(fields.Certificate)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorReadCertificateFile, fields.Certificate)
		}
		req.SetCertificate(string(content))
	}

	if cmd.Flags().Changed("private-key") {
		content, err := os.ReadFile(fields.PrivateKey)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorReadPrivateKeyFile, fields.PrivateKey)
		}
		req.SetPrivateKey(string(content))
	}

	if cmd.Flags().Changed("certificate-type") {
		req.SetType(fields.CertificateType)
	}

	return doCreate(ctx, f, client, req)
}

func runRequest(ctx context.Context, cmd *cobra.Command, f *cmdutil.Factory, client *api.Client, fields *Fields) error {
	req := api.NewRequestRequest()

	if !cmd.Flags().Changed("name") {
		answer, err := utils.AskInput(msg.AskName)
		if err != nil {
			return err
		}
		fields.Name = answer
	}
	req.SetName(fields.Name)

	req.SetAuthority(fields.Authority)

	if !cmd.Flags().Changed("challenge") {
		answer, err := utils.AskInput(msg.AskChallenge)
		if err != nil {
			return err
		}
		fields.Challenge = answer
	}
	req.SetChallenge(fields.Challenge)

	if !cmd.Flags().Changed("common-name") {
		answer, err := utils.AskInput(msg.AskCommonName)
		if err != nil {
			return err
		}
		fields.CommonName = answer
	}
	req.SetCommonName(fields.CommonName)

	if cmd.Flags().Changed("alternative-names") {
		names := strings.Split(fields.AlternativeNames, ",")
		for i := range names {
			names[i] = strings.TrimSpace(names[i])
		}
		req.SetAlternativeNames(names)
	}

	if cmd.Flags().Changed("key-algorithm") {
		req.SetKeyAlgorithm(fields.KeyAlgorithm)
	}

	return doRequest(ctx, f, client, req)
}

func doCreate(ctx context.Context, f *cmdutil.Factory, client *api.Client, req *api.CreateRequest) error {
	response, err := client.Create(ctx, req)
	if err != nil {
		return fmt.Errorf(msg.ErrorCreateDigitalCertificate.Error(), err)
	}

	out := output.GeneralOutput{
		Msg: fmt.Sprintf(msg.CreateOutputSuccess, response.GetId()),
		Out: f.IOStreams.Out,
	}
	return output.Print(&out)
}

func doRequest(ctx context.Context, f *cmdutil.Factory, client *api.Client, req *api.RequestRequest) error {
	response, err := client.Request(ctx, req)
	if err != nil {
		return fmt.Errorf(msg.ErrorRequestDigitalCertificate.Error(), err)
	}

	out := output.GeneralOutput{
		Msg: fmt.Sprintf(msg.RequestOutputSuccess, response.GetId()),
		Out: f.IOStreams.Out,
	}
	return output.Print(&out)
}

func readInput(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", utils.ErrorOpeningFile, path)
	}
	return data, nil
}

func strictUnmarshal(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

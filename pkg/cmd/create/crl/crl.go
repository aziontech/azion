package crl

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/create/crl"
	api "github.com/aziontech/azion-cli/pkg/api/crl"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Fields struct {
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
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create crl --name "My CRL" --issuer "My CA" --crl "./list.crl"
        $ azion create crl --name "My CRL" --issuer "My CA" --crl "./list.crl" --active true
        $ azion create crl --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			ctx := context.Background()

			req := api.NewCreateRequest()

			if cmd.Flags().Changed("file") {
				if err := requestFromFile(fields.InPath, req); err != nil {
					return err
				}
				return doCreate(ctx, f, client, req)
			}

			if err := requestFromFlags(cmd, fields, req); err != nil {
				return err
			}

			return doCreate(ctx, f, client, req)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)

	return cmd
}

func requestFromFlags(cmd *cobra.Command, fields *Fields, req *api.CreateRequest) error {
	if !cmd.Flags().Changed("name") {
		answer, err := utils.AskInput(msg.AskName)
		if err != nil {
			return err
		}
		fields.Name = answer
	}
	req.SetName(fields.Name)

	if !cmd.Flags().Changed("issuer") {
		answer, err := utils.AskInput(msg.AskIssuer)
		if err != nil {
			return err
		}
		fields.Issuer = answer
	}
	req.SetIssuer(fields.Issuer)

	if !cmd.Flags().Changed("crl") {
		answer, err := utils.AskInput(msg.AskCRL)
		if err != nil {
			return err
		}
		fields.CRL = answer
	}
	content, err := os.ReadFile(fields.CRL)
	if err != nil {
		return fmt.Errorf("%w: %s", msg.ErrorReadCRLFile, fields.CRL)
	}
	req.SetCrl(string(content))

	if cmd.Flags().Changed("active") {
		isActive, err := strconv.ParseBool(fields.Active)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorActiveFlag, fields.Active)
		}
		req.SetActive(isActive)
	}

	setComputedTimestamps(req)

	return nil
}

// createInput holds the user-provided attributes of a CRL. It is used instead
// of unmarshalling directly into the SDK model, whose strict UnmarshalJSON
// requires server-computed fields (id, last_editor, created_at, ...).
type createInput struct {
	Name   string `json:"name"`
	Issuer string `json:"issuer"`
	Crl    string `json:"crl"`
	Active *bool  `json:"active,omitempty"`
}

func requestFromFile(path string, req *api.CreateRequest) error {
	var input createInput
	if err := utils.FlagFileUnmarshalJSON(path, &input); err != nil {
		return msg.ErrorInvalidJSON
	}

	req.SetName(input.Name)
	req.SetIssuer(input.Issuer)
	req.SetCrl(input.Crl)
	if input.Active != nil {
		req.SetActive(*input.Active)
	}

	setComputedTimestamps(req)

	return nil
}

// setComputedTimestamps fills the last_modified, last_update and next_update
// fields, which the API requires on creation but recomputes server-side from
// the CRL content.
func setComputedTimestamps(req *api.CreateRequest) {
	now := time.Now()
	req.SetLastModified(now)
	req.SetLastUpdate(now)
	req.SetNextUpdate(now)
}

func doCreate(ctx context.Context, f *cmdutil.Factory, client *api.Client, req *api.CreateRequest) error {
	response, err := client.Create(ctx, req)
	if err != nil {
		return fmt.Errorf(msg.ErrorCreateCRL.Error(), err)
	}

	out := output.GeneralOutput{
		Msg:   fmt.Sprintf(msg.CreateOutputSuccess, response.GetId()),
		Out:   f.IOStreams.Out,
		Flags: f.Flags,
	}
	return output.Print(&out)
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.Issuer, "issuer", "", msg.FlagIssuer)
	flags.StringVar(&fields.CRL, "crl", "", msg.FlagCRL)
	flags.StringVar(&fields.Active, "active", "", msg.FlagActive)
	flags.StringVar(&fields.InPath, "file", "", msg.FlagIn)
	flags.BoolP("help", "h", false, msg.CreateFlagHelp)
}

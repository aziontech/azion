package update

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/domains"
	api "github.com/aziontech/azion-cli/pkg/api/domains"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Id                int64
	EdgeApplicationId int64
	Name              string
	CnameAccessOnly   string
	Active            string
	InPath            string
	Cnames            []string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.DomainUpdateUsage,
		Short:         msg.DomainUpdateShortDescription,
		Long:          msg.DomainUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azioncli domains update --domain-id 1234 --name 'Hello'
		$ azioncli domains update --domain-id 1234 --application-id 4321
		$ azioncli domains update --domain-id 1234 --cnames www.testhere.com,www.pudim.com
		$ azioncli domains update -d 9123 --cname-access-only true
		$ azioncli domains update -d 9123 --cname-access-only false
		$ azioncli domains update -d 9123 --application-id 192837
		$ azioncli domains update --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			// either domain-id or in path should be passed
			if !cmd.Flags().Changed("domain-id") && !cmd.Flags().Changed("in") {
				return msg.ErrorMissingApplicationIdArgument
			}

			request := api.UpdateRequest{}

			if cmd.Flags().Changed("in") {
				var (
					file *os.File
					err  error
				)
				if fields.InPath == "-" {
					file = os.Stdin
				} else {
					file, err = os.Open(fields.InPath)
					if err != nil {
						return fmt.Errorf("%w: %s", utils.ErrorOpeningFile, fields.InPath)
					}
				}
				err = cmdutil.UnmarshallJsonFromReader(file, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {

				request.Id = fields.Id

				if cmd.Flags().Changed("name") {
					request.SetName(fields.Name)
				}

				if cmd.Flags().Changed("application-id") {
					request.SetEdgeApplicationId(fields.EdgeApplicationId)
				}

				if cmd.Flags().Changed("cnames") {
					request.SetCnames(fields.Cnames)
				}

				if cmd.Flags().Changed("cname-access-only") {
					active, err := strconv.ParseBool(fields.CnameAccessOnly)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorActiveFlag, fields.CnameAccessOnly)
					}
					request.SetCnameAccessOnly(active)
				}

				if cmd.Flags().Changed("active") {
					active, err := strconv.ParseBool(fields.Active)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorActiveFlag, fields.Active)
					}
					request.SetIsActive(active)
				}

			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, &request)

			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateDomain.Error(), err)
			}

			fmt.Fprintf(f.IOStreams.Out, msg.DomainUpdateOutputSuccess, response.GetId())

			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.Id, "domain-id", "d", 0, msg.DomainFlagId)
	flags.Int64VarP(&fields.EdgeApplicationId, "application-id", "a", 0, msg.ApplicationFlagId)
	flags.StringVar(&fields.Name, "name", "", msg.DomainUpdateFlagName)
	flags.StringSliceVar(&fields.Cnames, "cnames", []string{}, msg.DomainUpdateFlagCnames)
	flags.StringVar(&fields.Active, "active", "", msg.DomainUpdateFlagActive)
	flags.StringVar(&fields.CnameAccessOnly, "cname-access-only", "", msg.DomainUpdateFlagCnameAccessOnly)
	flags.StringVar(&fields.InPath, "in", "", msg.DomainUpdateFlagIn)
	flags.BoolP("help", "h", false, msg.DomainUpdateHelpFlag)

	return cmd
}

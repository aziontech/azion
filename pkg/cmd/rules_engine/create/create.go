package create

import (
    "context"
    "fmt"
    "os"

    "github.com/MakeNowJust/heredoc"

    msg "github.com/aziontech/azion-cli/messages/rules_engine"
    api "github.com/aziontech/azion-cli/pkg/api/edge_applications"

    "github.com/aziontech/azion-cli/pkg/cmdutil"
    "github.com/aziontech/azion-cli/utils"
    "github.com/spf13/cobra"
)

type Criterio struct {
    Conditional string `json:"conditional"`
    Variable    string `json:"variable"`
    Operator    string `json:"operator"`
    InputValue  string `json:"input_value"`
}

type Behaviors struct {
    Name string `json:"name"`
}

type Fields struct {
    Name          string       `json:"name"`
    Criteria      [][]Criterio `json:"criteria"`
    Behaviors     []Behaviors  `json:"behaviors"`
}

var (
    ApplicationID int64 
    Phase         string
    Path          string      
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
    fields := &Fields{}

    cmd := &cobra.Command{
        Use:           msg.RulesEngineCreateUsage,
        Short:         msg.RulesEngineCreateShortDescription,
        Long:          msg.RulesEngineCreateLongDescription,
        SilenceUsage:  true,
        SilenceErrors: true,
        Example: heredoc.Doc(`
        `),
        RunE: func(cmd *cobra.Command, args []string) error {
            if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("phase") {
                return msg.ErrorMandatoryCreateFlags
            }

            request := api.CreateRulesEngineRequest{}
            if cmd.Flags().Changed("in") {
                var (
                    file *os.File
                    err  error
                )
                if Path == "-" {
                    file = os.Stdin
                } else {
                    file, err = os.Open(Path)
                    if err != nil {
                        return fmt.Errorf("%w: %s", utils.ErrorOpeningFile, Path)
                    }
                }
                err = cmdutil.UnmarshallJsonFromReader(file, &request)
                if err != nil {
                    return utils.ErrorUnmarshalReader
                }
            } else {
                request.SetName(fields.Name)
            }

            client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
            response, err := client.CreateRulesEngine(context.Background(), ApplicationID, Phase, &request)
            if err != nil {
                return fmt.Errorf(msg.ErrorCreateRulesEngine.Error(), err)
            }
            fmt.Fprintf(f.IOStreams.Out, msg.RulesEngineCreateOutputSuccess, response.GetId())
            return nil
        },
    }

    flags := cmd.Flags()
    flags.Int64VarP(&ApplicationID, "application-id", "a", 0, msg.RulesEngineCreateFlagEdgeApplicationID)
    flags.StringVar(&fields.Name, "name", "", msg.RulesEngineCreateFlagName)
    flags.StringVar(&Phase, "phase", "", msg.RulesEngineCreateFlagPhase)
    flags.StringVar(&Path, "in", "", msg.RulesEngineCreateFlagIn)
    flags.BoolP("help", "h", false, msg.RulesEngineCreateHelpFlag)
    return cmd
}

package template

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/spf13/cobra"
)

type TemplateCmd struct {
	WriteFile func(filename string, data []byte, perm fs.FileMode) error
	f         *cmdutil.Factory
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(newTemplateCmd(f))
}

func newTemplateCmd(f *cmdutil.Factory) *TemplateCmd {
	return &TemplateCmd{
		WriteFile: os.WriteFile,
		f:         f,
	}
}

func NewCobraCmd(tempCmd *TemplateCmd) *cobra.Command {
	var out string
	template := api.UpdateRulesEngineRequest{}

	cmd := &cobra.Command{
		Use:           msg.RulesEngineTemplateUsage,
		Short:         msg.RulesEngineTemplateShortDescription,
		Long:          msg.RulesEngineTemplateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azioncli rules_engine template
		$ azioncli rules_engine template --out /path/to/your/file.json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			template.SetName("NewRulesEngine")

			beh := sdk.RulesEngineBehavior{}
			beh.SetName("run-function")
			beh.SetTarget(0)

			b := make([]sdk.RulesEngineBehavior, 1)
			b[0].SetName("run_function")
			b[0].SetTarget(0)
			template.SetBehaviors(b)

			c := make([][]sdk.RulesEngineCriteria, 1)
			for i := 0; i < 1; i++ {
				c[i] = make([]sdk.RulesEngineCriteria, 1)
			}

			c[0][0].SetConditional("if")
			c[0][0].SetVariable("${uri}")
			c[0][0].SetOperator("starts_with")
			c[0][0].SetInputValue("/")
			template.SetCriteria(c)

			marshalledJson, err := json.MarshalIndent(template, "", " ")
			if err != nil {
				return err
			}

			err = tempCmd.WriteFile(out, marshalledJson, 0644)
			if err != nil {
				return msg.ErrorWriteTemplate
			}

			fmt.Fprintf(tempCmd.f.IOStreams.Out, msg.RulesEngineFileWritten, filepath.Clean(out))

			return nil
		},
	}

	flags := cmd.Flags()
	cmd.Flags().StringVar(&out, "out", "./rule.json", msg.RulesEngineTemplateFlagOut)

	flags.BoolP("help", "h", false, msg.RulesEngineTemplateHelpFlag)
	return cmd
}

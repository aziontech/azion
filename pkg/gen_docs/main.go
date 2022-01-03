package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aziontech/azion-cli/cmd"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	dir := flags.StringP("doc-path", "", "", "Path directory where you want generate doc files")
	// filetype := flags.Int("file-type", 0, "File type for generating the documentation. 1 - Yaml; 2 -ReST") // we could create this flag to be able to decide the format in which the docs will be generated
	help := flags.BoolP("help", "h", false, "Help about any command")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *help {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n%s", filepath.Base(args[0]), flags.FlagUsages())
		return nil
	}

	if *dir == "" {
		return fmt.Errorf("error: --doc-path not set")
	}

	factory := &cmdutil.Factory{
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Timeout: 10 * time.Second,
			}, nil
		},
		IOStreams: iostreams.System(),
	}

	rootCmd := cmd.NewRootCmd(factory)
	rootCmd.InitDefaultHelpCmd()

	if err := os.MkdirAll(*dir, 0755); err != nil {
		return err
	}

	err := doc.GenYamlTree(rootCmd, *dir)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

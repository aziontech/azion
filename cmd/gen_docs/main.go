package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	cmd "github.com/aziontech/azion-cli/pkg/cmd/root"
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
	filetype := flags.String("file-type", "", "File type for generating the documentation: <yaml|md>")
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

	rootCmd := cmd.NewCmd(&cmdutil.Factory{
		IOStreams: iostreams.System(),
	})
	rootCmd.InitDefaultHelpCmd()

	switch {
	case *filetype == "yaml":
		if err := os.MkdirAll(*dir, 0755); err != nil {
			return err
		}
		err := doc.GenYamlTree(rootCmd, *dir)
		if err != nil {
			log.Fatal(err)
		}
	case *filetype == "md":
		if err := os.MkdirAll(*dir, 0755); err != nil {
			return err
		}

		removeMdSuffix := func(s string) string { return strings.TrimRight(s, ".md") }
		dontPreprendFile := func(s string) string { return "" }

		err := doc.GenMarkdownTreeCustom(rootCmd, *dir, dontPreprendFile, removeMdSuffix)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal(errors.New("You must provide a valid file type"))
	}

	return nil
}

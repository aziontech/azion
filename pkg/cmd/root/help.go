package root

import (
	"bytes"
	"fmt"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/general"
	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/text"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var hasFailed bool

func HasFailed() bool {
	return hasFailed
}

func nestedSuggestFunc(command *cobra.Command, arg string) {
	command.Printf("unknown command %q for %q\n", arg, command.CommandPath())

	var candidates []string
	if arg == "help" {
		candidates = []string{"--help"}
	} else {
		if command.SuggestionsMinimumDistance <= 0 {
			command.SuggestionsMinimumDistance = 2
		}
		candidates = command.SuggestionsFor(arg)
	}

	if len(candidates) > 0 {
		command.Print("\nDid you mean this?\n")
		for _, c := range candidates {
			command.Printf("\t%s\n", c)
		}
	}

	command.Print("\n")
}

func isRootCmd(command *cobra.Command) bool {
	return command != nil && !command.HasParent()
}

func rootHelpFunc(f *cmdutil.Factory, command *cobra.Command, args []string) {
	if isRootCmd(command.Parent()) && len(args) >= 2 && args[1] != "--help" && args[1] != "-h" {
		nestedSuggestFunc(command, args[1])
		hasFailed = true
		return
	}

	// in case root command is run only with --token flag, we only show the token being saved
	if command.Flags().Changed("token") {
		return
	}

	var (
		baseCommands   []string
		subcmdCommands []string
		examples       []string
	)

	for _, c := range command.Commands() {
		if c.Short == "" {
			continue
		}
		if c.Hidden {
			continue
		}

		s := rpad(c.Name(), c.NamePadding()) + c.Short

		if c.Annotations["Category"] == "skip" {
			continue
		} else if !isRootCmd(c.Parent()) {
			// Help of subcommand
			subcmdCommands = append(subcmdCommands, s)
			continue
		} else {
			baseCommands = append(baseCommands, s)
			continue
		}
	}

	type helpEntry struct {
		Title string
		Body  string
	}

	if len(command.Example) > 0 {
		examples = append(examples, command.Example)
	}

	longText := command.Long
	if longText == "" {
		longText = command.Short
	}

	styleTitle := color.Bold
	styleBody := color.FgHiWhite

	helpEntries := []helpEntry{}

	helpEntries = append(helpEntries, helpEntry{"", color.New(color.Bold).Sprint(fmt.Sprintf(msg.CliVersion, version.BinVersion))})

	if longText != "" {
		helpEntries = append(helpEntries, helpEntry{
			Title: color.New(styleTitle).Sprint("DESCRIPTION"),
			Body:  color.New(styleBody).Sprint(longText),
		})
	}

	helpEntries = append(helpEntries, helpEntry{
		Title: color.New(styleTitle).Sprint("SYNOPSIS"),
		Body:  color.New(styleBody).Sprint(command.UseLine()),
	})

	if len(examples) > 0 {
		helpEntries = append(helpEntries, helpEntry{
			Title: color.New(styleTitle).Sprint("EXAMPLES"),
			Body:  color.New(color.FgGreen).Sprint(strings.Join(examples, "\n")),
		})
	}

	if len(baseCommands) > 0 {
		helpEntries = append(helpEntries, helpEntry{
			Title: color.New(styleTitle).Sprint("AVAILABLE COMMANDS"),
			Body:  color.New(styleBody).Sprint(strings.Join(baseCommands, "\n")),
		})
	}

	if len(subcmdCommands) > 0 {
		helpEntries = append(helpEntries, helpEntry{
			Title: color.New(styleTitle).Sprint("AVAILABLE SUBCOMMANDS"),
			Body:  color.New(styleBody).Sprint(strings.Join(subcmdCommands, "\n")),
		})
	}

	flagUsages := command.LocalFlags().FlagUsages()
	if flagUsages != "" {
		helpEntries = append(helpEntries, helpEntry{
			Title: color.New(styleTitle).Sprint("LOCAL OPTIONS"),
			Body:  color.New(styleBody).Sprint(dedent(flagUsages)),
		})
	}

	inheritedFlagUsages := command.InheritedFlags().FlagUsages()
	if inheritedFlagUsages != "" {
		helpEntries = append(helpEntries, helpEntry{
			Title: color.New(styleTitle).Sprint("GLOBAL OPTIONS"),
			Body:  color.New(styleBody).Sprint(dedent(inheritedFlagUsages)),
		})
	}

	helpEntries = append(helpEntries, helpEntry{
		Title: color.New(styleTitle).Sprint("LEARN MORE"),
		Body:  color.New(styleBody).Sprint("\nUse 'azion <command> <subcommand> --help' for more information about a command"),
	})

	out := command.OutOrStdout()
	for _, e := range helpEntries {
		if e.Title != "" {
			// If there is a title, add indentation to each line in the body
			fmt.Fprintln(out, e.Title)
			fmt.Fprintln(out, text.Indent(strings.Trim(e.Body, "\r\n"), "  "))
		} else {
			// If there is no title print the body as is
			fmt.Fprintln(out, e.Body)
		}
		fmt.Fprintln(out)
	}
}

func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds ", padding)
	return fmt.Sprintf(template, s)
}

func dedent(s string) string {
	lines := strings.Split(s, "\n")
	minIndent := -1

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}

		indent := len(l) - len(strings.TrimLeft(l, " "))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return s
	}

	var buf bytes.Buffer
	for _, l := range lines {
		fmt.Fprintln(&buf, strings.TrimPrefix(l, strings.Repeat(" ", minIndent)))
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

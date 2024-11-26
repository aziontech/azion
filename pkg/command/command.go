//go:build !windows

package command

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
)

// RunCommandWithOutput returns the stringified command output, it's exit code and any errors
// Commands that exit with exit codes > 0 will return a non-nil error
func RunCommandWithOutput(envVars []string, comm string) (string, int, error) {
	command := exec.Command(SHELL, "-c", comm)
	if len(envVars) > 0 {
		command.Env = os.Environ()
		command.Env = append(command.Env, envVars...)
	}

	out, err := command.CombinedOutput()
	exitCode := command.ProcessState.ExitCode()

	return string(out), exitCode, err
}

// CommandRunInteractive runs a command interactively.
func CommandRunInteractiveWithOutput(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
	cmd := exec.Command(SHELL, "-c", comm)
	if len(envVars) > 0 {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, envVars...)
	}
	var stdoutBuffer bytes.Buffer

	if !f.Silent {
		cmd.Stdin = f.IOStreams.In
		cmd.Stdout = &stdoutBuffer
	}

	cmd.Stderr = f.IOStreams.Err

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	output := stdoutBuffer.String()

	return output, nil
}

// CommandRunInteractive runs a command interactively.
func CommandRunInteractive(f *cmdutil.Factory, comm string) error {
	cmd := exec.Command(SHELL, "-c", comm)

	if !f.Silent && !(len(f.Flags.Format) > 0) && !(len(f.Flags.Out) > 0) {
		cmd.Stdin = f.IOStreams.In
		cmd.Stdout = f.IOStreams.Out
	}

	cmd.Stderr = f.IOStreams.Err
	return cmd.Run()
}

// RunCommandStreamOutput executes the provived command while streaming its logs (stdout+stderr) directly to terminal
func RunCommandStreamOutput(out io.Writer, envVars []string, comm string) error {
	command := exec.Command(SHELL, "-c", comm)
	if len(envVars) > 0 {
		command.Env = os.Environ()
		command.Env = append(command.Env, envVars...)
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	multi := io.MultiReader(stdout, stderr)

	// start the command after having set up the pipe
	if err := command.Start(); err != nil {
		return fmt.Errorf(utils.ErrorRunningCommandStream.Error(), err)
	}

	// read command's stdout line by line
	in := bufio.NewScanner(multi)

	for in.Scan() {
		fmt.Fprintf(out, "%s\n", in.Text())
	}
	if err := in.Err(); err != nil {
		return fmt.Errorf(utils.ErrorRunningCommandStream.Error(), err)
	}

	return nil
}

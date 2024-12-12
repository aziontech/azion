//go:build windows

package command

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
)

// RunCommandStreamOutput executes the provided PowerShell command while streaming its logs (stdout+stderr) directly to the terminal
func RunCommandStreamOutput(out io.Writer, envVars []string, comm string) error {
	// Use PowerShell to execute the command
	command := exec.Command(POWERSHELL, "-NoProfile", "-Command", comm)

	// Set environment variables if provided
	if len(envVars) > 0 {
		command.Env = os.Environ()
		command.Env = append(command.Env, envVars...)
	}

	// Setup pipes for stdout and stderr
	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	multi := io.MultiReader(stdout, stderr)

	// Start the command after setting up pipes
	if err := command.Start(); err != nil {
		return fmt.Errorf("error starting PowerShell command: %w", err)
	}

	// Read the command's output line by line
	in := bufio.NewScanner(multi)

	for in.Scan() {
		fmt.Fprintf(out, "%s\n", in.Text())
	}

	// Check for scanning errors
	if err := in.Err(); err != nil {
		return fmt.Errorf("error reading command output: %w", err)
	}

	// Wait for the command to complete
	if err := command.Wait(); err != nil {
		return fmt.Errorf("PowerShell command execution failed: %w", err)
	}

	return nil
}

func RunCommandWithOutput(envVars []string, comm string) (string, int, error) {
	// Use PowerShell to execute the command
	command := exec.Command(POWERSHELL, "-NoProfile", "-Command", comm)

	// Set environment variables if provided
	if len(envVars) > 0 {
		command.Env = os.Environ()
		command.Env = append(command.Env, envVars...)
	}

	// Run the command and capture combined stdout and stderr
	out, err := command.CombinedOutput()

	// Retrieve the exit code
	var exitCode int
	if command.ProcessState != nil {
		exitCode = command.ProcessState.ExitCode()
	} else {
		exitCode = -1 // Use -1 if the process state is unavailable
	}

	return string(out), exitCode, err
}

// CommandRunInteractiveWithOutput runs a PowerShell command interactively and captures its output.
func CommandRunInteractiveWithOutput(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
	// Use PowerShell to execute the command
	cmd := exec.Command(POWERSHELL, "-NoProfile", "-Command", comm)

	// Set environment variables if provided
	if len(envVars) > 0 {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, envVars...)
	}

	var stdoutBuffer bytes.Buffer

	// Configure input/output streams based on factory settings
	if !f.Silent {
		cmd.Stdin = f.IOStreams.In // Forward stdin
		cmd.Stdout = &stdoutBuffer
	}

	cmd.Stderr = f.IOStreams.Err // Forward stderr

	// Run the command interactively
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	output := stdoutBuffer.String()
	return output, nil
}

// CommandRunInteractive runs a PowerShell command interactively.
func CommandRunInteractive(f *cmdutil.Factory, comm string) error {
	// Use PowerShell to execute the command
	cmd := exec.Command(POWERSHELL, "-NoProfile", "-Command", comm)

	// Configure input/output streams for interactivity if not in silent mode and no specific flags are set
	if !f.Silent && !(len(f.Flags.Format) > 0) && !(len(f.Flags.Out) > 0) {
		cmd.Stdin = f.IOStreams.In
		cmd.Stdout = f.IOStreams.Out
	}

	cmd.Stderr = f.IOStreams.Err

	return cmd.Run()
}

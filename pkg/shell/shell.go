package shell

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func Get() (string, error) {
	if runtime.GOOS == "windows" {

		// Check if PowerShell is in use
		cmdPowerShell := exec.Command("powershell", "-Command", "$PSVersionTable.PSVersion.Major")
		cmdPowerShell.Stderr = os.Stderr
		outputPowerShell, errPowerShell := cmdPowerShell.Output()
		if errPowerShell == nil && strings.TrimSpace(string(outputPowerShell)) != "" {
			return "powershell", nil
		} else {
			// Checks if WSL is in use
			cmdWSL := exec.Command("cmd", "/C", "echo", "%WSL_DISTRO_NAME%")
			cmdWSL.Stderr = os.Stderr
			outputWSL, errWSL := cmdWSL.Output()
			if errWSL == nil && strings.TrimSpace(string(outputWSL)) != "" {
				outputUnix, errUnix := echoShell()
				if errUnix == nil {
					return strings.TrimSpace(string(outputUnix)), nil
				}
				// WSL is in use
				return "wsl", nil
			} else {
				// Otherwise, assume cmd
				return "cmd", nil
			}
		}
	} else {
		// If it's not on Windows, it assumes a Unix/Linux environment
		outputUnix, errUnix := echoShell()
		if errUnix == nil {
			return strings.TrimSpace(string(outputUnix)), nil
		}
	}

	return "/bin/bash", errors.New("shell not found")
}

func echoShell() ([]byte, error) {
	cmdUnix := exec.Command("sh", "-c", "echo $SHELL")
	cmdUnix.Stderr = os.Stderr
	return cmdUnix.Output()
}

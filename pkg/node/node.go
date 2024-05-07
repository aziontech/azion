package node

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func NodeVersion() error {
	cmd := exec.Command("node", "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return errors.New(NODE_NOT_INSTALLED)
	}
	return checkNode(out.String())
}

func checkNode(str string) error {
	versionOutput := strings.TrimSpace(str)
	if len(versionOutput) > 0 && versionOutput[0] == 'v' {
		versionOutput = versionOutput[1:]
	}

	versionParts := strings.Split(versionOutput, ".")
	if len(versionParts) > 0 {
		majorVersion := versionParts[0]

		var major int
		fmt.Sscanf(majorVersion, "%d", &major)

		if major < 18 {
			return errors.New(NODE_OLDER_VERSION)
		}
	}
	return nil
}

package node

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aziontech/azion-cli/utils"
)

type NodePkg struct {
	NodeVer    func(node *NodePkg) error
	CheckNode  func(str string) error
	CmdBuilder func(name string, arg ...string) *exec.Cmd
}

func NewNode() *NodePkg {
	return &NodePkg{
		NodeVer:    nodeVersion,
		CheckNode:  checkNode,
		CmdBuilder: exec.Command,
	}
}

func nodeVersion(node *NodePkg) error {
	cmd := node.CmdBuilder("node", "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return errors.New(NODE_NOT_INSTALLED)
	}
	return node.CheckNode(out.String())
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
		_, err := fmt.Sscanf(majorVersion, "%d", &major)
		if err != nil {
			return err
		}

		if major < 18 {
			return errors.New(NODE_OLDER_VERSION)
		}
	}
	return nil
}

func DetectPackageManager(pathWorkDir string) string {
	npmLockFile := filepath.Join(pathWorkDir, "package-lock.json")
	yarnLockFile := filepath.Join(pathWorkDir, "yarn.lock")
	pnpmLockFile := filepath.Join(pathWorkDir, "pnpm-lock.yaml")

	switch {
	case utils.FileExists(npmLockFile):
		return "npm"
	case utils.FileExists(yarnLockFile):
		return "yarn"
	case utils.FileExists(pnpmLockFile):
		return "pnpm"
	default:
		return "npm"
	}
}

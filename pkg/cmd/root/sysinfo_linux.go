//go:build linux

package root

import (
	"encoding/json"
	"os/user"

	msg "github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/zcalusic/sysinfo"
	"go.uber.org/zap"
)

func showUpdadeMessageSystem(f *cmdutil.Factory) error {
	current, err := user.Current()
	if err != nil {
		logger.Debug("Error while getting current user's information", zap.Error(err))
		return msg.ErrorCurrentUser
	}

	if current.Uid != "0" {
		logger.FInfo(f.IOStreams.Out, msg.CouldNotGetUser)
		return nil
	}

	var si sysinfo.SysInfo

	si.GetSysInfo()

	data, err := json.MarshalIndent(&si, "", "  ")
	if err != nil {
		logger.Debug("Error while marshaling current user's information", zap.Error(err))
		return msg.ErrorMarshalUserInfo
	}

	var osInfo OSInfo
	err = json.Unmarshal(data, &osInfo)
	if err != nil {
		logger.Debug("Error while unmarshaling current user's information", zap.Error(err))
		return msg.ErrorUnmarshalUserInfo
	}

	logger.FInfo(f.IOStreams.Out, msg.DownloadRelease)
	switch osInfo.OS.Vendor {
	case "debian":
		logger.FInfo(f.IOStreams.Out, msg.DpkgUpdate)
	case "alpine":
		logger.FInfo(f.IOStreams.Out, msg.ApkUpdate)
	case "centos", "fedora", "opensuse", "mageia", "mandriva":
		logger.FInfo(f.IOStreams.Out, msg.RpmUpdate)
	}

	return nil
}

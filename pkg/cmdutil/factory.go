package cmdutil

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
)

type Factory struct {
	HttpClient *http.Client
	IOStreams  *iostreams.IOStreams
	Config     config.Config
	Flags
}

type Flags struct {
	logger.Logger
	GlobalFlagAll bool   `json:"-" yaml:"-" toml:"-"`
	Out           string `json:"-" yaml:"-" toml:"-"`
	Format        string `json:"-" yaml:"-" toml:"-"`
	NoColor       bool   `json:"-" yaml:"-" toml:"-"`
}

// GetActiveProfile returns the active profile name from config, defaulting to "default" if empty
func (f *Factory) GetActiveProfile() string {
	if f.Config == nil {
		return "default"
	}

	activeProfile := f.Config.GetString("active_profile")

	if activeProfile == "" {
		return "default"
	}

	// If the value looks like a path, extract just the profile name (last component)
	if strings.Contains(activeProfile, "/") {
		return filepath.Base(activeProfile)
	}

	return activeProfile
}

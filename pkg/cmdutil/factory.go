package cmdutil

import (
	"net/http"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
)

type Factory struct {
	HttpClient *http.Client
	IOStreams  *iostreams.IOStreams
	Config     config.Config
	logger.Logger
}

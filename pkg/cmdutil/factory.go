package cmdutil

import (
	"net/http"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/iostreams"
)

type Factory struct {
	HttpClient *http.Client
	IOStreams  *iostreams.IOStreams
	Config     config.Config
}

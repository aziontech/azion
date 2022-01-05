package cmdutil

import (
	"net/http"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/iostreams"
)

type Factory struct {
	HttpClient func() (*http.Client, error)
	IOStreams  *iostreams.IOStreams
	Config     config.Config
}

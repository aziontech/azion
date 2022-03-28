package testutils

import (
	"bytes"
	"net/http"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/spf13/viper"
)

func NewFactory(mock *httpmock.Registry) (factory *cmdutil.Factory, out *bytes.Buffer, err *bytes.Buffer) {
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	f := &cmdutil.Factory{
		HttpClient: &http.Client{Transport: mock},
		IOStreams: &iostreams.IOStreams{
			Out: stdout,
			Err: stderr,
		},
		Config: viper.New(),
	}
	return f, stdout, stderr
}

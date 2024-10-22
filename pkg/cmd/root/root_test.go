package root

import (
	"net/http"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/constants"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/viper"
)

func TestExecute(t *testing.T) {
	streams := iostreams.System()
	httpClient := &http.Client{
		Timeout: 50 * time.Second, // TODO: Configure this somewhere
	}

	tok, _ := token.ReadSettings()
	viper.SetEnvPrefix("AZIONCLI")
	viper.AutomaticEnv()
	viper.SetDefault("token", tok.Token)
	viper.SetDefault("api_url", constants.ApiURL)
	viper.SetDefault("storage_url", constants.StorageApiURL)

	factory := &cmdutil.Factory{
		HttpClient: httpClient,
		IOStreams:  streams,
		Config:     viper.GetViper(),
	}

	Execute(factory)
}

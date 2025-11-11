package main

import (
	"net/http"
	"path"
	"time"

	cmd "github.com/aziontech/azion-cli/pkg/cmd/root"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/constants"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/viper"
)

func main() {
	streams := iostreams.System()
	httpClient := &http.Client{
		Timeout: 50 * time.Second,
	}

	profiles, settingsPath, _ := token.ReadProfiles()
	activeProfile := path.Join(settingsPath, profiles.Name)
	tok, _ := token.ReadSettings(profiles.Name)
	viper.SetEnvPrefix("AZIONCLI")
	viper.AutomaticEnv()
	viper.SetDefault("token", tok.Token)
	viper.SetDefault("api_url", constants.ApiURL)
	viper.SetDefault("api_v4_url", constants.ApiV4URL)
	viper.SetDefault("storage_url", constants.StorageApiURL)
	viper.SetDefault("active_profile", activeProfile)

	factory := &cmdutil.Factory{
		HttpClient: httpClient,
		IOStreams:  streams,
		Config:     viper.GetViper(),
	}

	cmd.Execute(cmd.NewFactoryRoot(factory))
}

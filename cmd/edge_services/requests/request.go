package requests

import (
	"fmt"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/token"
	sdk "github.com/aziontech/edgeservices-go-sdk"
	"github.com/spf13/cobra"
)

var ApiUrl string

func CreateClient(f *cmdutil.Factory, cmd *cobra.Command) (*sdk.APIClient, error) {
	var (
		tok string
		err error
	)

	if cmd.Flags().Changed("token") {
		tok, err = cmd.Flags().GetString("token")
	} else {
		tok, err = token.ReadFromDisk()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get api token: %w", err)
	}

	conf := sdk.NewConfiguration()

	httpClient, err := f.HttpClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get http client: %w", err)
	}

	conf.HTTPClient = httpClient

	conf.AddDefaultHeader("Authorization", "token "+tok)
	conf.Servers = sdk.ServerConfigurations{
		{
			URL: ApiUrl,
		},
	}

	return sdk.NewAPIClient(conf), nil
}

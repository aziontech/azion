package deploy

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
)

func checkToken(f *cmdutil.Factory) error {
	configureToken := f.Config.GetString("token")

	t, err := token.New(&token.Config{
		Client: f.HttpClient,
		Out:    f.IOStreams.Out,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", utils.ErrorTokenManager, err)
	}

	if configureToken == "" {
		return utils.ErrorTokenNotProvided
	}

	valid, _, err := t.Validate(&configureToken)
	if err != nil {
		return err
	}
	if !valid {
		return msg.ErrorInvalidToken
	}

	return nil
}

func openBrowser(f *cmdutil.Factory, urlConsoleDeploy string, cmd *DeployCmd) error {
	logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msg.VisitMsg, urlConsoleDeploy))
	err := cmd.OpenBrowserFunc(urlConsoleDeploy)
	if err != nil {
		return err
	}
	return nil
}

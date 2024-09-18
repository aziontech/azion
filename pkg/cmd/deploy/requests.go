package deploy

import (
	"fmt"

	"github.com/skratchdot/open-golang/open"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
)

func checkToken(f *cmdutil.Factory) error {
	configureToken := f.Config.GetString("token")

	t := token.New(&token.Config{
		Client: f.HttpClient,
		Out:    f.IOStreams.Out,
	})

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

func openBrowser(f *cmdutil.Factory, urlConsoleDeploy string) error {
	logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msg.VisitMsg, urlConsoleDeploy))
	err := open.Run(urlConsoleDeploy)
	if err != nil {
		return err
	}
	return nil
}

package purge

import (
	"context"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/purge"
	apipurge "github.com/aziontech/azion-cli/pkg/api/realtime_purge"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"

	"go.uber.org/zap"
)

func purgeWildcard(urls []string, f *cmdutil.Factory) error {
	if len(urls) > 1 {
		return msg.ErrorTooManyUrls
	}
	ctx := context.Background()

	clipurge := apipurge.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	err := clipurge.PurgeWildcard(ctx, urls)
	if err != nil {
		logger.Debug("Error while purging domains", zap.Error(err))
		return err
	}

	return nil
}

func purgeUrls(urls []string, f *cmdutil.Factory) error {
	ctx := context.Background()

	clipurge := apipurge.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	err := clipurge.PurgeUrls(ctx, urls)
	if err != nil {
		logger.Debug("Error while purging URLs", zap.Error(err))
		return err
	}

	return nil
}

func purgeCacheKeys(urls []string, f *cmdutil.Factory) error {
	ctx := context.Background()

	clipurge := apipurge.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	err := clipurge.PurgeCacheKey(ctx, urls, Layer)
	if err != nil {
		logger.Debug("Error while purging domains", zap.Error(err))
		return err
	}

	return nil
}

func getPurgeType() (string, error) {
	opts := []string{"URLs", "Wildcard", "Cache-Key"}
	answer := ""
	prompt := &survey.Select{
		Message: "Choose a purge type:",
		Options: opts,
	}
	err := survey.AskOne(prompt, &answer)
	if err != nil {
		return "", err
	}
	return answer, nil
}

func askForInput() ([]string, error) {
	var userInput string
	prompt := &survey.Input{
		Message: msg.AskForInput,
	}

	// Prompt the user for input
	err := survey.AskOne(prompt, &userInput, survey.WithKeepFilter(true))
	if err != nil {
		return []string{}, err
	}

	listOfUrls := strings.Split(userInput, ",")

	return listOfUrls, nil
}

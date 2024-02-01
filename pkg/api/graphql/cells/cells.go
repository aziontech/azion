package cells

import (
	"context"
	"fmt"
	"time"

	msg "github.com/aziontech/azion-cli/messages/logs/cells"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/machinebox/graphql"
	"go.uber.org/zap"
)

type CellsConsoleEvent struct {
	Ts              time.Time `json:"ts"`
	SolutionId      string    `json:"solutionId"`
	ConfigurationId string    `json:"configurationId"`
	FunctionId      string    `json:"functionId"`
	ID              string    `json:"id"`
	LineSource      string    `json:"lineSource"`
	Level           string    `json:"level"`
	Line            string    `json:"line"`
}

const query string = `
query ConsoleLog {
	cellsConsoleEvents(
	  %s
	  filter: {
		%s
		tsGt: "%s"
		
	  }
	  orderBy: [ts_ASC]
	) {
	  ts
	  solutionId
	  configurationId
	  functionId
	  id
	  lineSource
	  level
	  line
	}
  }	  
`

type CellsConsoleEventsResponse struct {
	CellsConsoleEvents []CellsConsoleEvent `json:"cellsConsoleEvents"`
}

func CellsConsoleLogs(f *cmdutil.Factory, functionId string, currentTime time.Time, limitFlag string) (CellsConsoleEventsResponse, error) {
	graphqlClient := graphql.NewClient("https://api.azionapi.net/events/graphql")

	formattedTime := currentTime.Format("2006-01-02T15:04:05")

	filter := ""
	if functionId != "" {
		filter = fmt.Sprintf(`functionId: "%s"`, functionId)
	}

	limit := "limit: " + limitFlag

	//prepare query
	formattedQuery := fmt.Sprintf(query, limit, filter, formattedTime)

	graphqlRequest := graphql.NewRequest(formattedQuery)

	tokenvalue := f.Config.GetString("token")
	token := "Token " + tokenvalue

	graphqlRequest.Header.Set("Authorization", token)

	var response CellsConsoleEventsResponse
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &response); err != nil {
		logger.Debug("", zap.Any("Error", err.Error()))
		return CellsConsoleEventsResponse{}, msg.ErrorRequest
	}

	return response, nil
}

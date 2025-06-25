package http

import (
	"context"
	"fmt"
	"time"

	msg "github.com/aziontech/azion-cli/messages/logs/http"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/machinebox/graphql"
	"go.uber.org/zap"
)

type HTTPEvent struct {
	Host              string    `json:"host"`
	GeolocCountry     string    `json:"geolocCountryName"`
	GeolocRegion      string    `json:"geolocRegionName"`
	HTTPUserAgent     string    `json:"httpUserAgent"`
	RequestURI        string    `json:"requestUri"`
	Status            int       `json:"status"`
	Ts                time.Time `json:"ts"`
	UpstreamBytesSent int       `json:"upstreamBytesSent"`
	RequestTime       string    `json:"requestTime"`
	RequestMethod     string    `json:"requestMethod"`
}

type HTTPEventsResponse struct {
	HTTPEvents []HTTPEvent `json:"httpEvents"`
}

const query string = `
query HttpEventsLogs {
	httpEvents(
	  %s
	  filter: {
	   tsGt: "%s"
	  }
	  orderBy: [ts_ASC]
	)
	{
	  host
	  httpUserAgent
	  geolocRegionName
	  requestUri
	  status
	  ts
	  upstreamBytesSent
	  requestTime
	  requestMethod
	}
  }
`

func HttpEvents(f *cmdutil.Factory, currentTime time.Time, limitFlag string) (HTTPEventsResponse, error) {
	graphqlClient := graphql.NewClient("https://api.azionapi.net/events/graphql")

	formattedTime := currentTime.Format("2006-01-02T15:04:05")

	limit := "limit: " + limitFlag

	//prepare query
	formattedQuery := fmt.Sprintf(query, limit, formattedTime)

	graphqlRequest := graphql.NewRequest(formattedQuery)

	tokenvalue := f.Config.GetString("token")
	token := "Token " + tokenvalue

	graphqlRequest.Header.Set("Authorization", token)

	var response HTTPEventsResponse
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &response); err != nil {
		logger.Debug("", zap.Any("Error", err.Error()))
		return HTTPEventsResponse{}, msg.ErrorRequest
	}

	return response, nil
}

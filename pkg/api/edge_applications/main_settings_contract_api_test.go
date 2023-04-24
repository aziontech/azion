package edgeapplications

import (
	"bytes"
	"encoding/json"
	"fmt"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/go-faker/faker/v4"
	"net/http"
	"testing"
)

var (
	name                                 = faker.Name()
	deliveryProtocol                     = "http"
	originType                           = "single_origin"
	address                              = "www.new.api"
	originProtocolPolicy                 = "preserve"
	hostHeader                           = "www.new.api"
	browserCacheSettings                 = "override"
	browserCacheSettingsMaximumTTL int64 = 20
	cdnCacheSettings                     = "honor"
	cdnCacheSettingsMaximumTTL     int64 = 60
)

var request = sdk.CreateApplicationRequest{
	Name:                           name,
	DeliveryProtocol:               &deliveryProtocol,
	OriginType:                     &originType,
	Address:                        &address,
	OriginProtocolPolicy:           &originProtocolPolicy,
	HostHeader:                     &hostHeader,
	BrowserCacheSettings:           &browserCacheSettings,
	BrowserCacheSettingsMaximumTtl: &browserCacheSettingsMaximumTTL,
	CdnCacheSettings:               &cdnCacheSettings,
	CdnCacheSettingsMaximumTtl:     &cdnCacheSettingsMaximumTTL,
}

func getRequest() []byte {
	b, _ := json.Marshal(request)
	return b
}

func token() string {
	return "Token <token>"
}

func TestEdgeApplicationsAPIContract(t *testing.T) {
	const host string = "https://api.azionapi.net"

	tests := []struct {
		name           string
		headers        map[string]string
		methodHttp     string
		path           string
		bodyReq        []byte
		want           sdk.GetApplicationResponse
		wantErr        bool
		StatusExpected int
	}{
		{
			name:       "main settings /edge_applications get list, status OK",
			methodHttp: http.MethodGet,
			headers: map[string]string{
				"Accept":        "application/json; version=3",
				"Authorization": token(),
				"Content-Type":  "application/json",
			},
			path:           "/edge_applications",
			bodyReq:        []byte(""),
			wantErr:        false,
			StatusExpected: http.StatusOK,
		},
		{
			name:       "main settings /edge_applications/:id get item, status OK",
			methodHttp: http.MethodGet,
			headers: map[string]string{
				"Accept":        "application/json; version=3",
				"Authorization": token(),
				"Content-Type":  "application/json",
			},
			path:           "/edge_applications/1673635839",
			bodyReq:        []byte(""),
			wantErr:        false,
			StatusExpected: http.StatusOK,
		},
		{
			name:       "main settings /edge_applications post item, status OK",
			methodHttp: http.MethodPost,
			headers: map[string]string{
				"Accept":        "application/json; version=3",
				"Authorization": token(),
				"Content-Type":  "application/json",
			},
			path:           "/edge_applications",
			bodyReq:        getRequest(),
			wantErr:        false,
			StatusExpected: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%v%v", host, tt.path)
			req, err := http.NewRequest(tt.methodHttp, url, bytes.NewBuffer(tt.bodyReq))
			if (err != nil) != tt.wantErr {
				t.Errorf("error on request %v", err)
				return
			}

			for key, value := range tt.headers {
				req.Header.Set(key, fmt.Sprintf("%v", value))
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("error on request %v", err)
				return
			}
			defer resp.Body.Close()
		})
	}
}

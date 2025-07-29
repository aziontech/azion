package root

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	msg "github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/constants"
	"github.com/aziontech/azion-cli/pkg/logger"
)

// Response structure with only the needed field
type AccountInfo struct {
	ClientFlags []string `json:"client_flags"`
}

func HasBlockAPIV4Flag(token string, f *factoryRoot) (bool, error) {
	if token == "" {
		logger.FInfoFlags(f.factory.IOStreams.Out, msg.LoginMessage, f.factory.Format, f.factory.Out)
		return false, nil
	}
	url := constants.AuthURL + "/account/info"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json; version=1")
	req.Header.Set("Authorization", "Token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("non-200 response: %d, body: %s", resp.StatusCode, string(body))
	}

	var info AccountInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return false, err
	}

	for _, flag := range info.ClientFlags {
		if flag == "block_apiv4_incompatible_endpoints" {
			return true, nil
		}
	}
	return false, nil
}

package deploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

var ResponseMap map[string]string

// Response represents the structure of the response from the API.
type Response struct {
	UUID  string    `json:"uuid"`
	Image string    `json:"image"`
	Start time.Time `json:"start"`
}

func callScript(token, id, secret, prefix, name, confDir string, cmd *DeployCmd) (string, error) {
	logger.Debug("Calling script runner api")
	instantiateURL := fmt.Sprintf("%s/api/template-engine/templates/%s/instantiate", DeployURL, ScriptID)

	// Define the request payload
	payload := []map[string]string{
		{
			"field":                   "AZCLI_PROJECT_NAME",
			"instantiation_data_path": "envs.[0].value",
			"value":                   name,
		},
		{
			"field":                   "AZCLI_B2_APP_KEY_ID",
			"instantiation_data_path": "envs.[1].value",
			"value":                   id,
		},
		{
			"field":                   "AZCLI_B2_APP_KEY",
			"instantiation_data_path": "envs.[2].value",
			"value":                   secret,
		},
		{
			"field":                   "AZCLI_PREFIX",
			"instantiation_data_path": "envs.[3].value",
			"value":                   prefix,
		},
		{
			"field":                   "AZCLI_TOKEN",
			"instantiation_data_path": "envs.[4].value",
			"value":                   token,
		},
		{
			"field":                   "AZCLI_CONFDIR",
			"instantiation_data_path": "envs.[5].value",
			"value":                   confDir,
		},
	}

	// Marshal the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Debug("Error marshalling payload", zap.Error(err))
		return "", err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", instantiateURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		logger.Debug("Error creating request", zap.Error(err))
		return "", err
	}

	// Set headers
	req.Header.Set("accept", "application/json; version=3")
	req.Header.Set("content-type", "application/json; version=3")
	req.Header.Set("Authorization", "Token "+token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Debug("Error sending request", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Unmarshal the response body into the Response struct
	if err := cmd.Unmarshal(body, &ResponseMap); err != nil {
		logger.Debug("Error unmarshalling response", zap.Error(err))
		return "", err
	}

	return ResponseMap["uuid"], nil
}

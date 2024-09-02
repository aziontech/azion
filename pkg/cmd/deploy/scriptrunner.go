package deploy

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

// Response represents the structure of the response from the API.
type Response struct {
	UUID  string    `json:"uuid"`
	Image string    `json:"image"`
	Start time.Time `json:"start"`
}

func callScript(token, id, secret, prefix, name string) (string, error) {
	logger.Debug("Calling script runner api")
	url := "https://stage-console.azion.com/api/template-engine/templates/92480a31-b88b-495b-8615-3ed5eff6314e/instantiate"

	// Define the request payload
	payload := []map[string]string{
		{
			"field":                   "PROJECT_NAME",
			"instantiation_data_path": "envs.[0].value",
			"value":                   name,
		},
		{
			"field":                   "B2_APP_KEY_ID",
			"instantiation_data_path": "envs.[1].value",
			"value":                   id,
		},
		{
			"field":                   "B2_APP_KEY",
			"instantiation_data_path": "envs.[2].value",
			"value":                   secret,
		},
		{
			"field":                   "PREFIX",
			"instantiation_data_path": "envs.[3].value",
			"value":                   prefix,
		},
	}

	// Marshal the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Debug("Error marshalling payload", zap.Error(err))
		return "", err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
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
	var responseMap map[string]string
	if err := json.Unmarshal(body, &responseMap); err != nil {
		logger.Debug("Error unmarshalling response", zap.Error(err))
		return "", err
	}

	return responseMap["uuid"], nil
}

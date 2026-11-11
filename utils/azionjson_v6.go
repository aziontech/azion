package utils

// ============================================================================
// v6 — owned by API v6. Same file and errors as the v3 and v4 helpers in
// helpers.go; only the contract differs.
// ============================================================================

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

// GetAzionJsonContentV6 reads azion.json as a v6 project. A missing file is
// not an error: a v6 deploy can start from nothing and writes the file as it
// goes, so it returns an empty contract instead.
func GetAzionJsonContentV6(confPath string) (*contracts.AzionApplicationOptionsV6, error) {
	wd, err := GetWorkingDir()
	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(path.Join(wd, confPath, "azion.json"))
	if errors.Is(err, os.ErrNotExist) {
		return &contracts.AzionApplicationOptionsV6{}, nil
	}
	if err != nil {
		logger.Debug("Error reading azion.json file", zap.Error(err))
		return nil, ErrorOpeningAzionJsonFile
	}

	if isV4AzionJson(file) {
		return fromV4AzionJson(file)
	}

	conf := &contracts.AzionApplicationOptionsV6{}
	if err := json.Unmarshal(file, conf); err != nil {
		logger.Debug("Error unmarshalling azion.json file", zap.Error(err))
		return nil, ErrorUnmarshalAzionJsonFile
	}
	return conf, nil
}

// v4OnlyKeys are azion.json keys only the v4 format has. "workloads" is in
// both, but as an object in v4 and a list in v6.
var v4OnlyKeys = []string{"application", "domain", "rules-engine", "rt-purge", "function", "cache-settings", "origin"}

// isV4AzionJson reports whether azion.json was written by the v4 commands, as
// "azion init" did for every account before it learned about v6.
func isV4AzionJson(file []byte) bool {
	var keys map[string]json.RawMessage
	if json.Unmarshal(file, &keys) != nil {
		return false
	}
	for _, k := range v4OnlyKeys {
		if _, ok := keys[k]; ok {
			return true
		}
	}
	w := bytes.TrimSpace(keys["workloads"])
	return len(w) > 0 && w[0] == '{'
}

// fromV4AzionJson starts a v6 project from a v4 azion.json. Only what both
// generations mean the same way is kept: the project's name, preset, bucket,
// prefix and env. The v4 resource ids do not exist in v6, so the v6 fields
// start empty, and not-first-run is dropped so the v6 placeholders get
// replaced and the deploy finds or creates the resources by name. The file
// is rewritten in the v6 format the first time the deploy saves it.
func fromV4AzionJson(file []byte) (*contracts.AzionApplicationOptionsV6, error) {
	var v4 struct {
		Name         string `json:"name"`
		Bucket       string `json:"bucket"`
		Preset       string `json:"preset"`
		Env          string `json:"env"`
		Prefix       string `json:"prefix"`
		RotatePrefix *bool  `json:"rotate-prefix"`
	}
	if err := json.Unmarshal(file, &v4); err != nil {
		logger.Debug("Error unmarshalling a v4 azion.json file", zap.Error(err))
		return nil, ErrorUnmarshalAzionJsonFile
	}
	logger.Debug("azion.json is in the v4 format; starting the v6 project from its name, preset, bucket and prefix")
	return &contracts.AzionApplicationOptionsV6{
		Name:         v4.Name,
		Bucket:       v4.Bucket,
		Preset:       v4.Preset,
		Env:          v4.Env,
		Prefix:       v4.Prefix,
		RotatePrefix: v4.RotatePrefix,
	}, nil
}

// WriteAzionJsonContentV6 writes azion.json for a v6 project.
func WriteAzionJsonContentV6(conf *contracts.AzionApplicationOptionsV6, confPath string) error {
	wd, err := GetWorkingDir()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		logger.Debug("Error marshalling azion.json", zap.Error(err))
		return ErrorMarshalAzionJsonFile
	}
	if err := os.MkdirAll(path.Join(wd, confPath), 0755); err != nil {
		logger.Debug("Error creating the config directory", zap.Error(err))
		return ErrorWritingAzionJsonFile
	}
	if err := os.WriteFile(path.Join(wd, confPath, "azion.json"), data, 0644); err != nil {
		logger.Debug("Error writing azion.json", zap.Error(err))
		return ErrorWritingAzionJsonFile
	}
	return nil
}

// IsAzionJsonV6 reports whether azion.json in confPath is in the v6 format,
// as a v6 deploy writes it. A missing, unreadable or v4 file is not.
func IsAzionJsonV6(confPath string) bool {
	wd, err := GetWorkingDir()
	if err != nil {
		return false
	}
	file, err := os.ReadFile(path.Join(wd, confPath, "azion.json"))
	if err != nil {
		return false
	}
	var keys map[string]json.RawMessage
	if json.Unmarshal(file, &keys) != nil {
		return false
	}
	return !isV4AzionJson(file)
}

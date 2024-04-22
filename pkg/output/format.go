package output

import (
	"encoding/json"
	"fmt"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

const (
	JSON = "json"
	YAML = "yaml"
	YML  = "yml"
	TOML = "toml"
)

func format(v any, g GeneralOutput) error {
	var b []byte
	var err error

	switch g.FlagFormat {
	case JSON:
		b, err = json.MarshalIndent(v, "", " ")
		if err != nil {
			return err
		}
	case YAML, YML:
		b, err = yaml.Marshal(v)
		if err != nil {
			return err
		}
	case TOML:
		b, err = toml.Marshal(v)
		if err != nil {
			return err
		}
	default:
		b, err = json.MarshalIndent(v, "", " ")
		if err != nil {
			return err
		}
	}

	if len(g.FlagOutPath) > 0 {
		err = cmdutil.WriteDetailsToFile(b, g.FlagOutPath, g.Out)
		if err != nil {
			return fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
		}
		logger.FInfo(g.Out, fmt.Sprintf(WRITE_SUCCESS, g.FlagOutPath))
		return nil
	}

	logger.FInfo(g.Out, string(b))
	return nil
}

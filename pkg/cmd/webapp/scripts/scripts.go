package scripts

import (
	"encoding/json"
	msg "github.com/aziontech/azion-cli/messages/webapp"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	"os"
)

func getConfig() (conf *contracts.AzionApplicationConfig, err error) {
	path, err := utils.GetWorkingDir()
	if err != nil {
		return conf, err
	}
	jsonConf := path + "/azion/config.json"
	file, err := os.ReadFile(jsonConf)
	if err != nil {
		return conf, msg.ErrorOpeningConfigFile
	}
	conf = &contracts.AzionApplicationConfig{}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		return conf, msg.ErrorUnmarshalConfigFile
	}
	if conf.InitData.Cmd == "" {
		return conf, msg.ErrorWebappInitCmdNotSpecified
	}
	return conf, nil
}

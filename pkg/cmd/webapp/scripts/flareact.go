package scripts

import (
	"github.com/aziontech/azion-cli/pkg/cmd/webapp/build"
	"github.com/aziontech/azion-cli/pkg/cmd/webapp/init"
)

func InitFlareact(info *init.InitInfo, cmd *init.InitCmd) (string, int, error) {
	//_, err := cmd.LookPath("npm")
	//if err != nil {
	//	return "", 0, errors.New("npm not found")
	//}
	//
	//conf, err := getConfig(cmd)
	//if err != nil {
	//	return "", 0, err
	//}
	//
	//envs, err := cmd.EnvLoader(conf.InitData.Env)
	//if err != nil {
	//	return "", 0, errors.New("failed load envs err: " + err.Error())
	//}
	//
	//fmt.Fprintf(cmd.Io.Out, msg.WebappInitRunningCmd)
	//fmt.Fprintf(cmd.Io.Out, "$ %s\n", conf.BuildData.Cmd)
	//
	//output, exitCode, err := cmd.CommandRunner(conf.BuildData.Cmd, envs)
	//
	//packageJsonPath := info.PathWorkingDir + "/package.json"
	//packageJson, err := cmd.FileReader(packageJsonPath)
	//if err != nil {
	//	return "", 0, errors.New("failed on read file")
	//}
	//
	//packJsonReplaceBuild, err := sjson.Set(string(packageJson), "scripts.build", "azioncli webapp build")
	//if err != nil {
	//	return "", 0, errors.New("failed replace scripts.build")
	//}
	//
	//packJsonReplaceDeploy, err := sjson.Set(packJsonReplaceBuild, "scripts.deploy", "azioncli webapp publish")
	//if err != nil {
	//	return "", 0, errors.New("failed replace scripts.deploy")
	//}
	//
	//cmd.WriteFile(packageJsonPath, []byte(packJsonReplaceDeploy), 0644)
	//
	return "", 0, nil
}

func BuildFlareact(cmd *build.BuildCmd) (string, int, error) {
	//conf, err := getConfig()
	//if err != nil {
	//	return "", 0, err
	//}
	//
	//envs, err := cmd.EnvLoader(conf.InitData.Env)
	//if err != nil {
	//	return "", 0, errors.New("failed load envs err: " + err.Error())
	//}
	//
	//workDirPath, err := cmd.GetWorkDir()
	//
	//workDirPath += "/args.json"
	//_, err = cmd.FileReader(workDirPath)
	//if err != nil {
	//	cmd.WriteFile(workDirPath, []byte("{}"), 0644)
	//}
	//
	//fmt.Fprintf(cmd.Io.Out, msg.WebappBuildRunningCmd)
	//fmt.Fprintf(cmd.Io.Out, "$ %s\n", conf.BuildData.Cmd)
	//
	//return cmd.CommandRunner(conf.BuildData.Cmd, envs)

	return "", 0, nil
}

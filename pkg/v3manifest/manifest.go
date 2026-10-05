package manifest

import (
	"encoding/json"
	"errors"
	"os"

	msg "github.com/aziontech/azion-cli/messages/manifest"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	"github.com/aziontech/azion-cli/utils"
)

var (
	CacheIds         map[string]int64
	CacheIdsBackup   map[string]int64
	RuleIds          map[string]contracts.RuleIdsStruct
	OriginKeys       map[string]string
	OriginIds        map[string]int64
	manifestFilePath = "/.edge/manifest.json"
)

type ManifestInterpreter struct {
	FileReader            func(path string) ([]byte, error)
	GetWorkDir            func() (string, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptionsV3, confPath string) error
}

func NewManifestInterpreter() *ManifestInterpreter {
	return &ManifestInterpreter{
		FileReader:            os.ReadFile,
		GetWorkDir:            utils.GetWorkingDir,
		WriteAzionJsonContent: utils.WriteAzionJsonContentV3,
	}
}

func (man *ManifestInterpreter) ManifestPath() (string, error) {
	pathWorkingDir, err := man.GetWorkDir()
	if err != nil {
		return "", err
	}

	return utils.Concat(pathWorkingDir, manifestFilePath), nil
}

func (man *ManifestInterpreter) ReadManifest(
	path string, f *cmdutil.Factory, msgs *[]string) (*contracts.Manifest, error) {
	logger.FInfoFlags(f.IOStreams.Out, msg.ReadingManifest, f.Format, f.Out)
	*msgs = append(*msgs, msg.ReadingManifest)
	manifest := &contracts.Manifest{}

	byteManifest, err := man.FileReader(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteManifest, &manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

// CreateResources applies the manifest by running the v3 pipeline.
//
// The order of the steps, and the one that is conditional, live in V3Pipeline.
// pkg/pipeline owns what used to sit around each block here: the announcement
// line, the debug log per step and stopping at the first error.
func (man *ManifestInterpreter) CreateResources(
	conf *contracts.AzionApplicationOptionsV3,
	manifest *contracts.Manifest,
	f *cmdutil.Factory,
	projectConf string,
	msgs *[]string) error {

	rc := newResourceContext(f, conf, manifest, projectConf, msgs, man.WriteAzionJsonContent)

	// No spinner and no timing callback: v3 announces with a plain line and has
	// never reported per-step timings, unlike v4.
	_, err := pipeline.Run(rc.ctx, V3Pipeline{}, rc, &pipeline.Plan{
		Factory:  f,
		Msgs:     msgs,
		Announce: msg.CreatingManifest,
	})
	if errors.Is(err, errPurgeAborted) {
		// A failed purge stops the run and reports success, skipping orphan
		// removal. See errPurgeAborted.
		return nil
	}
	return err
}

package contracts

import "os"

type FileOps struct {
	Path        string
	MimeType    string
	FileContent *os.File
	VersionID   string
}

type BuildInfo struct {
	Preset        string
	Mode          string
	Entry         string
	NodePolyfills string
	OwnWorker     string
	IsFirewall    bool
}

type DevInfo struct {
	IsFirewall string
}

type ListOptions struct {
	Details           bool
	OrderBy           string
	Sort              string
	Page              int64
	PageSize          int64
	Filter            string
	ContinuationToken string
}

type DescribeOptions struct {
	OutPath string
	Format  string
}

type AzionApplicationOptions struct {
	Test        func(path string) error  `json:"-"`
	Name        string                   `json:"name"`
	Bucket      string                   `json:"bucket"`
	Preset      string                   `json:"preset"` // framework: react, next, vue, angular and etc
	Mode        string                   `json:"mode"`   // deliver == ssg, compute == ssr
	Env         string                   `json:"env"`
	Prefix      string                   `json:"prefix"`
	Function    AzionJsonDataFunction    `json:"function"`
	Application AzionJsonDataApplication `json:"application"`
	Domain      AzionJsonDataDomain      `json:"domain"`
	RtPurge     AzionJsonDataPurge       `json:"rt-purge"`
	Origin      AzionJsonDataOrigin      `json:"origin"`
	RulesEngine AzionJsonDataRulesEngine `json:"rules-engine"`
}

type AzionApplicationSimple struct {
	Name        string                   `json:"name"`
	Type        string                   `json:"type"`
	Domain      AzionJsonDataDomain      `json:"domain"`
	Application AzionJsonDataApplication `json:"application"`
}

type AzionApplicationConfig struct {
	InitData    InitConf    `json:"init"`
	BuildData   BuildConf   `json:"build"`
	PublishData PublishConf `json:"publish"`
}

type InitConf struct {
	Cmd        string `json:"cmd"`
	Env        string `json:"env"`
	OutputCtrl string `json:"output-ctrl"`
	Default    string `json:"default"`
}

type BuildConf struct {
	Cmd        string `json:"cmd"`
	Env        string `json:"env"`
	OutputCtrl string `json:"output-ctrl"`
	Default    string `json:"default"`
}

type PublishConf struct {
	Cmd        string `json:"pre_cmd"`
	Env        string `json:"env"`
	OutputCtrl string `json:"output-ctrl"`
	Default    string `json:"default"`
}

type CacheConf struct {
	PurgeOnPublish bool `json:"purge_on_publish"`
}

type AzionJsonDataFunction struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	File         string `json:"file"`
	Args         string `json:"args"`
	InstanceID   int64  `json:"instance-id"`
	InstanceName string `json:"instance-name"`
}

type AzionJsonDataApplication struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type AzionJsonDataOrigin struct {
	SingleOriginID   int64    `json:"single-origin-id"`
	StorageOriginID  int64    `json:"storage-origin-id"`
	StorageOriginKey string   `json:"storage-origin-key"`
	Name             string   `json:"name"`
	Address          []string `json:"address"`
}

type AzionJsonDataDomain struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type AzionJsonDataPurge struct {
	PurgeOnPublish bool `json:"purge_on_publish"`
}

type AzionJsonDataRulesEngine struct {
	Created bool `json:"created"`
}

package contracts

type ListOptions struct {
	Details  bool
	OrderBy  string
	Sort     string
	Page     int64
	PageSize int64
	Filter   string
}

type DescribeOptions struct {
	OutPath string
	Format  string
}

type AzionApplicationOptions struct {
	Test        func(path string) error  `json:"-"`
	Name        string                   `json:"name"`
	Env         string                   `json:"env"`
	Function    AzionJsonDataFunction    `json:"function"`
	Application AzionJsonDataApplication `json:"Application"`
	Domain      AzionJsonDataDomain      `json:"Domain"`
	RtPurge     AzionJsonDataPurge       `json:"RTPurge"`
}

type AzionApplicationConfig struct {
	InitData    InitConf    `json:"init"`
	BuildData   BuildConf   `json:"build"`
	PublishData PublishConf `json:"publish"`
}

type InitConf struct {
	Cmd string `json:"cmd"`
	Env string `json:"env"`
}

type BuildConf struct {
	Cmd string `json:"cmd"`
	Env string `json:"env"`
}

type PublishConf struct {
	Cmd string `json:"pre_cmd"`
	Env string `json:"env"`
}

type CacheConf struct {
	PurgeOnPublish bool `json:"purge_on_publish"`
}

type AzionJsonDataFunction struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	File     string `json:"file"`
	Args     string `json:"args"`
	Language string `json:"language"`
}

type AzionJsonDataApplication struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type AzionJsonDataDomain struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type AzionJsonDataPurge struct {
	Name bool `json:"purge_on_publish"`
}

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
	Test         func(path string) error `json:"-"`
	Name         string                  `json:"name"`
	Language     string                  `json:"language"`
	Env          string                  `json:"env"`
	FunctionFile string                  `json:"function_file"`
	CacheData    CacheConf               `json:"cache"`
}

type AzionApplicationConfig struct {
	InitData    InitConf  `json:"init"`
	BuildData   BuildConf `json:"build"`
	PublishData InitConf  `json:"publish"`
}

type InitConf struct {
	Cmd string `json:"pre_cmd"`
	Env string `json:"env"`
}

type BuildConf struct {
	Cmd string `json:"cmd"`
	Env string `json:"env"`
}

type PublishConf struct {
	Cmd string `json:"cmd"`
	Env string `json:"env"`
}

type CacheConf struct {
	PurgeOnPublish bool `json:"purge_on_publish"`
}

type AzionJsonData struct {
	Name        string `json:"name"`
	Env         string `json:"env"`
	Function    AzionJsonDataFunction
	Application AzionJsonDataApplication
	Domain      AzionJsonDataDomain
	RTPurge     AzionJsonDataPurge
}

type AzionJsonDataFunction struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	File     string `json:"file"`
	Args     string `json:"args"`
	Active   bool   `json:"active"`
	Language string `json:"language"`
}

type AzionJsonDataApplication struct {
	Name string `json:"name"`
}

type AzionJsonDataDomain struct {
	Name string `json:"name"`
}

type AzionJsonDataPurge struct {
	Name bool `json:"purge_on_publish"`
}

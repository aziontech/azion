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
	CacheData    cacheConf               `json:"cache"`
}

type AzionApplicationConfig struct {
	InitData  InitConf  `json:"init"`
	BuildData buildConf `json:"build"`
}

type InitConf struct {
	Cmd string `json:"cmd"`
	Env string `json:"env"`
}

type buildConf struct {
	Cmd string `json:"cmd"`
	Env string `json:"env"`
}

type cacheConf struct {
	PurgeOnPublish bool `json:"purge_on_publish"`
}

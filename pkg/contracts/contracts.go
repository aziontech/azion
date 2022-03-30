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
	Test         func() error `json:"-"`
	Name         string       `json:"name"`
	Language     string       `json:"language"`
	Env          string       `json:"env"`
	FunctionFile string       `json:"function_file"`
	CacheData    cacheConf    `json:"cache"`
}

type AzionApplicationConfig struct {
	InitData  initConf
	BuildData buildConf
}

type initConf struct {
	Cmd string
}

type buildConf struct {
	Cmd string
}

type cacheConf struct {
	PurgeOnPublish bool `json:"purge_on_publish"`
}

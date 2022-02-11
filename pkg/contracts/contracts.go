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

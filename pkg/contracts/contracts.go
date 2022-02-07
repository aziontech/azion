package contracts

type ListOptions struct {
	Details   bool
	Order_by  string
	Sort      string
	Page      int64
	Page_size int64
}

type DescribeOptions struct {
	OutPath string
	Format  string
}

package output

import "io"

type DescribeOutput struct {
	Field []string
	Value []string
	Out   io.Writer
}

func (c *DescribeOutput) Output() {
}

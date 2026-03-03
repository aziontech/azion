package output

import (
	"io"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
)

type RawOutput struct {
	Bytes []byte        `json:"-" yaml:"-" toml:"-"`
	Out   io.Writer     `json:"-" yaml:"-" toml:"-"`
	Flags cmdutil.Flags `json:"-" yaml:"-" toml:"-"`
}

func (r *RawOutput) Format() (bool, error) {
	formatted := false
	if len(r.Flags.Format) > 0 || len(r.Flags.Out) > 0 {
		formatted = true
		if len(r.Flags.Out) > 0 {
			err := WriteDetailsToFile(r.Bytes, r.Flags.Out)
			if err != nil {
				return formatted, err
			}
			logger.FInfo(r.Out, WRITE_SUCCESS+": "+r.Flags.Out)
			return formatted, nil
		}
		logger.FInfo(r.Out, string(r.Bytes))
		return formatted, nil
	}
	return formatted, nil
}

func (r *RawOutput) Output() {
	logger.FInfo(r.Out, string(r.Bytes))
}

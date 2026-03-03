package output

import (
	"io"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
)

// RawOutput is used for outputting raw bytes without any wrapping or formatting.
// This is useful for storage objects where the content should be preserved exactly as returned.
type RawOutput struct {
	Bytes []byte        `json:"-" yaml:"-" toml:"-"`
	Out   io.Writer     `json:"-" yaml:"-" toml:"-"`
	Flags cmdutil.Flags `json:"-" yaml:"-" toml:"-"`
}

func (r *RawOutput) Format() (bool, error) {
	formatted := false
	if len(r.Flags.Format) > 0 || len(r.Flags.Out) > 0 {
		formatted = true
		// For raw output, we write the bytes directly without any marshaling
		if len(r.Flags.Out) > 0 {
			err := WriteDetailsToFile(r.Bytes, r.Flags.Out)
			if err != nil {
				return formatted, err
			}
			logger.FInfo(r.Out, WRITE_SUCCESS+": "+r.Flags.Out)
			return formatted, nil
		}
		// Write raw bytes to output
		logger.FInfo(r.Out, string(r.Bytes))
		return formatted, nil
	}
	return formatted, nil
}

func (r *RawOutput) Output() {
	// For raw output, just write the bytes directly
	logger.FInfo(r.Out, string(r.Bytes))
}

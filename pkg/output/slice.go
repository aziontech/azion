package output

type SliceOutput struct {
	Messages []string `json:"messages" yaml:"messages" toml:"messages"`
	GeneralOutput
}

func (i *SliceOutput) Format() (bool, error) {
	formatted := false
	if len(i.Flags.Format) > 0 || len(i.Flags.Out) > 0 {
		formatted = true
		err := format(i, i.GeneralOutput)
		if err != nil {
			return formatted, err
		}
	}
	return true, nil
}

func (i *SliceOutput) Output() {}

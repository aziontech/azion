package output

type SliceOutput struct {
	Messages []string `json:"messages" yaml:"messages" toml:"messages"`
	GeneralOutput
}

func (i *SliceOutput) Format() (bool, error) {
	formated := false
	if len(i.Flags.Format) > 0 || len(i.Flags.Out) > 0 {
		formated = true
		err := format(i, *&i.GeneralOutput)
		if err != nil {
			return formated, err
		}
	}
	return true, nil
}

func (i *SliceOutput) Output() {}

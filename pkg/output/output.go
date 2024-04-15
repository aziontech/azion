package output

type TypeOutputInterface interface {
	Output()
	Format() (bool, error)
}

type Output struct {
	Output TypeOutputInterface
}

func Print(out TypeOutputInterface) error {
	format, err := out.Format()
	if err != nil {
		return err
	}

	// Check if the format is true, if yes we should not print again.
	if !format {
		out.Output()
	}
	return nil
}

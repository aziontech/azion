package output

type TypeOutputInterface interface {
	Output()
}

type Output struct {
	Output TypeOutputInterface
}

func Print(out TypeOutputInterface) {
	out.Output()
}

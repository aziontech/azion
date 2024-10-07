package output

type MockTypeOutputInterface struct {
	formatBool bool
	formatErr  error
}

func (m *MockTypeOutputInterface) Format() (bool, error) { return m.formatBool, m.formatErr }
func (m *MockTypeOutputInterface) Output()               {}

package cmdutil

type Command interface {
	Execute() error
	InitDefaultHelpCmd()
}

package install

const (
	Usage             = "install"
	ShortDescription  = "Install bundled resources"
	LongDescription   = "Installs bundled resources such as skills to the local system"
	Example           = "$ azion install --skills"
	FlagSkills        = "Install all bundled skills to ~/.claude/skills/"
	FlagHelp          = "Displays more information about the install command"
	MsgResolveHome    = "Resolving home directory..."
	MsgValidateSource = "Validating source skills directory..."
	MsgCreateTarget   = "Creating target directory: %s"
	MsgRemoveExisting = "Removing existing skill: %s"
	MsgCopySkill      = "Copying skill: %s"
	MsgDone           = "Successfully installed %d skill(s) to ~/.claude/skills/"
	MsgNoSkills       = "No skills found to install"
)

package commands

import "github.com/PersonalTool/code/go/go_analyzer/utility"

type Table3D struct {
	*CommandStruct
}

func (command *Table3D) Execute() error {
	utility.TestOutput("This is 3D Table")
	return nil
}

func (command *Table3D) parseCommandParams() error {
	return nil
}

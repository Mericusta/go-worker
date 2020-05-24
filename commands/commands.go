package commands

import "github.com/go-worker/global"

var commandFactory *CommandFactory

func init() {
	commandFactory = &CommandFactory{
		commandNo: 0,
	}
}

type CommandInterface interface {
	Execute() error
	parseCommandParams() error
}

type CommandStruct struct {
	No          int
	Type        global.CommandEnum
	InputString string
	Param       CommandParam
}

type CommandParam struct {
}

func CreateCommand(commandEnum global.CommandEnum, inputString string) CommandInterface {
	return commandFactory.CreateCommand(commandEnum, inputString)
}

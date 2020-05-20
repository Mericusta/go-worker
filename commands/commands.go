package commands

var commandFactory *CommandFactory

// Command 指令枚举类型
type CommandEnum int

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
	Type        CommandEnum
	InputString string
	Param       CommandParam
}

type CommandParam struct {
}

func CreateCommand(commandEnum CommandEnum, inputString string) CommandInterface {
	return commandFactory.CreateCommand(commandEnum, inputString)
}

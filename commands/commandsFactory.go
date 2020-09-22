package commands

import "github.com/go-worker/global"

type CommandFactory struct {
	commandNo int
}

func (commandFactory *CommandFactory) CreateCommand(commandEnum global.CommandEnum, inputString string) CommandInterface {
	var commandInterface CommandInterface
	command := &CommandStruct{
		No:          commandFactory.commandNo,
		Type:        commandEnum,
		InputString: inputString,
	}
	switch commandEnum {
	case global.CmdCustom:
		commandInterface = &Custom{
			CommandStruct: command,
		}
	case global.CmdBind:
		commandInterface = &Bind{
			CommandStruct: command,
		}
	case global.CmdExit:
		commandInterface = &Exit{
			CommandStruct: command,
		}
	case global.CmdCreate:
		commandInterface = &Create{
			CommandStruct: command,
		}
	case global.CmdConvert:
		commandInterface = &Convert{
			CommandStruct: command,
		}
	case global.CmdAnalyze:
		commandInterface = &Analyze{
			CommandStruct: command,
		}
	case global.CmdRemove:
		commandInterface = &Remove{
			CommandStruct: command,
		}
	case global.CmdSpider:
		commandInterface = &Spider{
			CommandStruct: command,
		}
	case global.CmdCSVChecker:
		commandInterface = &CSVChecker{
			CommandStruct: command,
		}
	}

	commandFactory.commandNo++
	return commandInterface
}

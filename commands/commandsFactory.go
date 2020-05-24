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
	case global.CmdBind:
		cmdBind := &Bind{
			CommandStruct: command,
		}
		commandInterface = cmdBind
	case global.CmdExit:
		cmdExit := &Exit{
			CommandStruct: command,
		}
		commandInterface = cmdExit
	case global.CmdCreate:
		cmdCreate := &Create{
			CommandStruct: command,
		}
		commandInterface = cmdCreate
	case global.CmdConvert:
		cmdConvert := &Convert{
			CommandStruct: command,
		}
		commandInterface = cmdConvert
	case global.CmdAnalyze:
		cmdAnalyze := &Analyze{
			CommandStruct: command,
		}
		commandInterface = cmdAnalyze
	}

	commandFactory.commandNo++
	return commandInterface
}

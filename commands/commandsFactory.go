package commands

type CommandFactory struct {
	commandNo int
}

func (commandFactory *CommandFactory) CreateCommand(commandEnum CommandEnum, inputString string) CommandInterface {
	var commandInterface CommandInterface
	command := &CommandStruct{
		No:          commandFactory.commandNo,
		Type:        commandEnum,
		InputString: inputString,
	}
	switch commandEnum {
	case CmdBind:
		cmdBind := &Bind{
			CommandStruct: command,
		}
		commandInterface = cmdBind
	case CmdExit:
		cmdExit := &Exit{
			CommandStruct: command,
		}
		commandInterface = cmdExit
	case CmdCreate:
		cmdCreate := &Create{
			CommandStruct: command,
		}
		commandInterface = cmdCreate
	case CmdConvert:
		cmdConvert := &Convert{
			CommandStruct: command,
		}
		commandInterface = cmdConvert
	}

	commandFactory.commandNo++
	return commandInterface
}

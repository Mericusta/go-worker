package commandsfsm

import (
	"github.com/go-worker/commands"
	"github.com/go-worker/fsm"
)

var commandsFSMStateMap map[commands.Command]fsm.FSMState

func init() {
	commandsFSMStateMap = map[commands.Command]fsm.FSMState{
		commands.Unknown: fsm.Error,
		commands.Exit:    fsm.Exiting,
	}
}

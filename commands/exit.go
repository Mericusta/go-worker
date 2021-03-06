package commands

import (
	"github.com/go-worker/fsm"
	"github.com/go-worker/global"
	"github.com/go-worker/ui"
)

type Exit struct {
	*CommandStruct
}

func (command *Exit) Execute() error {
	ui.OutputNoteInfoWithOutFormat(ui.CMDExit)
	global.FsmState = fsm.Exiting
	return nil
}

func (command *Exit) parseCommandParams() error {
	return nil
}

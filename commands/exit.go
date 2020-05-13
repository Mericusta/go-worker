package commands

import (
	"github.com/go-worker/fsm"
	"github.com/go-worker/global"
)

const exitParamIndex = 1

func init() {

}

func ExitExecutor(inputString string) error {
	global.FsmState = fsm.Exiting
	return nil
}

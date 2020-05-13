package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/go-worker/commands"
	"github.com/go-worker/fsm"
	"github.com/go-worker/global"
	"github.com/go-worker/regexpscommands"
	"github.com/go-worker/ui"
	"github.com/go-worker/utility"
)

func init() {
	global.FsmState = fsm.Waitting
}

func main() {
	for {
		fmt.Print(ui.FSMWaitting)
		input := bufio.NewScanner(os.Stdin)

		hasInput := input.Scan()
		if hasInput {
			global.FsmState = fsm.Waitting
			utility.TestOutput("inputString = %v", input.Text())

			// TODO: fsm 作为独立携程存在

			global.FsmState = fsm.Executing
			command, parseCommandError := regexpscommands.ParseCommandByRegexp(input.Text())
			if parseCommandError != nil {
				utility.ErrorOutput("%v", parseCommandError)
				continue
			}

			utility.TestOutput("command = %v", command)
			commandExecuteError := commands.Execute(input.Text(), command)
			if commandExecuteError != nil {
				ui.OutputErrorInfo(commandExecuteError)
			}

			if global.FsmState == fsm.Exiting {
				utility.NoteOutput("Thanks for using!")
				break
			} else if global.FsmState == fsm.Error {
				utility.TestOutput("FSM Error")
			} else {
				utility.TestOutput("FSM state = %v", global.FsmState)
			}
		}
	}
}

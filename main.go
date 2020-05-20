package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	var input *bufio.Scanner
	for {
		// 等待输入
		global.FsmState = fsm.Waitting
		fmt.Print(ui.FSMWaitting)
		input = bufio.NewScanner(os.Stdin)
		if input.Scan() {
			inputTextWithoutTrimSpace := strings.TrimSpace(input.Text())
			utility.TestOutput("inputString = %v", inputTextWithoutTrimSpace)

			// 解析输入
			global.FsmState = fsm.Parsing
			command, parseCommandError := regexpscommands.ParseCommandByRegexp(inputTextWithoutTrimSpace)
			if parseCommandError != nil {
				ui.OutputErrorInfo("%v", parseCommandError)
				continue
			}

			// 执行指令
			global.FsmState = fsm.Executing
			commandExecuteError := command.Execute()
			if commandExecuteError != nil {
				ui.OutputErrorInfo("%v", commandExecuteError)
			}

			// 状态机结果
			if global.FsmState == fsm.Exiting {
				break
			} else if global.FsmState == fsm.Error {
				utility.TestOutput("FSM Error")
			} else {
				utility.TestOutput("FSM state = %v", global.FsmState)
			}
		}
	}
}

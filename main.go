package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-worker/fsm"
	"github.com/go-worker/global"
	"github.com/go-worker/logger"
	"github.com/go-worker/regexpscommands"
	"github.com/go-worker/ui"
)

func init() {
	global.FsmState = fsm.Waitting
}

func main() {
	go logger.Run()

	var input *bufio.Scanner
	for {
		// 等待输入
		global.FsmState = fsm.Waitting
		fmt.Print(ui.FSMWaitting)
		input = bufio.NewScanner(os.Stdin)
		if input.Scan() {
			inputTextWithoutTrimSpace := strings.TrimSpace(input.Text())

			// 解析输入
			global.FsmState = fsm.Parsing
			command, parseCommandError := regexpscommands.ParseCommandByRegexp(inputTextWithoutTrimSpace)
			if parseCommandError != nil {
				ui.OutputErrorInfo("%v", parseCommandError)
				continue
			}

			// 执行指令
			global.FsmState = fsm.Executing
			logger.OutputNoteInfo(ui.CommonNote3, inputTextWithoutTrimSpace)
			commandExecuteError := command.Execute()
			if commandExecuteError != nil {
				ui.OutputErrorInfo("%v", commandExecuteError)
			}
			logger.OutputNoteInfo(ui.CommonNote2)

			// 状态机结果
			if global.FsmState == fsm.Exiting {
				break
			} else if global.FsmState == fsm.Error {
			} else {
			}
		}
	}

	logger.LogControl <- true
}

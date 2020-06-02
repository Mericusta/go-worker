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
			} else {
			}
		}
	}
}

// func ExampleFunc1(str string) {

// }

// func ExampleFunc2(str string) int {
// 	return 0
// }

// func ExampleFunc3(str string) (int, int) {
// 	return 0, 0
// }

// type ExampleStruct struct {
// 	value int
// }

// func (example ExampleStruct) ExampleFunc4() {

// }

// func (example ExampleStruct) ExampleFunc5(str string) {

// }

// func (example ExampleStruct) ExampleFunc6(str string) int {
// 	return 0
// }

// func (example ExampleStruct) ExampleFunc7(ex *ExampleStruct) *ExampleStruct {
// 	return nil
// }

// func (example ExampleStruct) ExampleFunc8(str string) (int, int) {
// 	return 0, 0
// }

// func (example ExampleStruct) ExampleFunc9(
// 	str string,
// 	value int,
// 	example2 ExampleStruct,
// ) (int, int) {
// 	return 0, 0
// }

// func (example *ExampleStruct) ExampleFunc10(
// 	str string,
// 	value int,
// 	example2 ExampleStruct,
// 	example3 *ExampleStruct,
// ) (int, int, *ExampleStruct) {
// 	return 0, 0, nil
// }

// func (example *ExampleStruct) ExampleFunc11(
// 	str string,
// 	value int,
// 	example2 ExampleStruct,
// 	example3 *ExampleStruct,
// ) (a int, b int, c *ExampleStruct) {
// 	return 0, 0, nil
// }

// func (example *ExampleStruct) ExampleFunc12(
// 	str string,
// 	value int,
// 	example2 ExampleStruct,
// 	example3 *ExampleStruct,
// 	list []string,
// 	asd ...string,
// ) (a int, b int, c *ExampleStruct) {
// 	return 0, 0, nil
// }

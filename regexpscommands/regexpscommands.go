package regexpscommands

import (
	"fmt"

	"github.com/go-worker/commands"
	"github.com/go-worker/regexps"
	"github.com/go-worker/ui"
)

var expressionCmdMap map[regexps.Expression]commands.CommandEnum

func init() {
	expressionCmdMap = map[regexps.Expression]commands.CommandEnum{
		regexps.ExpExit:    commands.CmdExit,
		regexps.ExpBind:    commands.CmdBind,
		regexps.ExpCreate:  commands.CmdCreate,
		regexps.ExpConvert: commands.CmdConvert,
	}
}

// ParseCommandByRegexp 根据输入使用正则表达式解析指令
// @param  inputString 待解析的
// @return 指令枚举，错误
func ParseCommandByRegexp(inputString string) (commands.CommandInterface, error) {
	for expression, regexp := range regexps.ExpressionRegexpMap {
		if regexp.MatchString(inputString) {
			if commandEnum, hasCommand := expressionCmdMap[regexps.Expression(expression)]; hasCommand {
				commandObject := commands.CreateCommand(commandEnum, inputString)
				if commandObject == nil {
					return nil, fmt.Errorf("create command[%v] but get nil", commandEnum)
				}
				return commandObject, nil
			}
		}
	}
	return nil, fmt.Errorf("%v", ui.FSMUnknownCommand)
}

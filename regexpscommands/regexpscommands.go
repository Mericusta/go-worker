package regexpscommands

import (
	"fmt"
	// "github.com/go-worker/fsm"
	"github.com/go-worker/commands"
	"github.com/go-worker/errors"
	"github.com/go-worker/regexps"
	"github.com/go-worker/utility"
)

var expressionCmdMap map[regexps.Expression]commands.Command

func init() {
	expressionCmdMap = map[regexps.Expression]commands.Command{
		regexps.Exit:      commands.Exit,
		regexps.Create:    commands.Create,
		regexps.Spider:    commands.Spider,
		regexps.Simulator: commands.Simulator,
		regexps.Worker:    commands.Worker,
	}
}

// ParseCommandByRegexp 根据输入使用正则表达式解析指令
// @param  inputString 待解析的
// @return 指令枚举，错误
func ParseCommandByRegexp(inputString string) (commands.Command, error) {
	utility.TestOutput("func ParseCommandByRegexp")
	for expression, regexp := range regexps.ExpressionRegexpMap {
		utility.TestOutput("expression = %v", expression)
		if regexp.MatchString(inputString) {
			if command, hasCommand := expressionCmdMap[expression]; hasCommand {
				return command, nil
			}
		}
	}
	return commands.Unknown, fmt.Errorf("%v", errors.UnknownCommand)
}

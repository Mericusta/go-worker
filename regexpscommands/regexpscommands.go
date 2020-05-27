package regexpscommands

import (
	"fmt"

	"github.com/go-worker/commands"
	"github.com/go-worker/regexps"
	"github.com/go-worker/ui"
)

// ParseCommandByRegexp 根据输入使用正则表达式解析指令
// @param  inputString 待解析的
// @return
func ParseCommandByRegexp(inputString string) (commands.CommandInterface, error) {
	for regexp, commandEnum := range regexps.RegexpCommandEnumMap {
		if regexp.MatchString(inputString) {
			commandObject := commands.CreateCommand(commandEnum, inputString)
			if commandObject == nil {
				return nil, fmt.Errorf("create command[%v] but get nil", commandEnum)
			}
			return commandObject, nil
		}
	}
	return nil, fmt.Errorf("%v", ui.FSMUnknownCommand)
}

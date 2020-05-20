package commands

import (
	"fmt"
	"strings"

	"github.com/go-worker/config"
	"github.com/go-worker/ui"
	"github.com/go-worker/utility"
)

var optionUIMap map[string]string

func init() {
	optionUIMap = make(map[string]string)

	optionUIMap["project"] = ui.CMDBindShowCurrentProject
	optionUIMap["syntax"] = ui.CMDBindShowCurrentSyntax
}

type Bind struct {
	*CommandStruct
	BindParams *bindParam
}

func (command Bind) Execute() error {
	// 解析指令的选项和参数
	parseCommandParamsError := command.parseCommandParams()
	if parseCommandParamsError != nil {
		return parseCommandParamsError
	}

	config.WorkerConfig.ProjectPath = command.BindParams.value

	optionUI, hasOptionUI := optionUIMap[command.BindParams.option]
	if !hasOptionUI {
		return fmt.Errorf(ui.CommonError2)
	}
	if command.BindParams.option == "project" && !utility.IsExist(command.BindParams.value) {
		ui.OutputNoteInfo(ui.CMDBindProjectNotExist)
		return nil
	}
	ui.OutputNoteInfoWithFormat(optionUI, command.BindParams.value)

	return nil
}

const (
	bindCommandIndex = 0
	bindOptionIndex  = 1
	bindValueIndex   = 2
)

type bindParam struct {
	option string
	value  string
}

func (command *Bind) parseCommandParams() error {
	inputStringList := strings.Split(command.InputString, " ")
	if bindValueIndex >= len(inputStringList) {
		return fmt.Errorf(ui.CommonError1)
	}
	command.BindParams = &bindParam{
		option: inputStringList[bindOptionIndex],
		value:  inputStringList[bindValueIndex],
	}
	return nil
}

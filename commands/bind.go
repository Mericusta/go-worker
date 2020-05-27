package commands

import (
	"fmt"
	"strings"

	"github.com/go-worker/config"
	"github.com/go-worker/global"
	"github.com/go-worker/regexps"
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

	if command.BindParams == nil {
		return fmt.Errorf("parse command bind, but param is nil")
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

type bindParam struct {
	option string
	value  string
}

func (command *Bind) parseCommandParams() error {
	optionValueString := ""
	if optionValueRegexp, hasOptionValueRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEBindOptionValue]; hasOptionValueRegexp {
		optionValueString = optionValueRegexp.FindString(command.CommandStruct.InputString)
	} else {
		ui.OutputWarnInfo(ui.CommonWarn2, "bind", "option")
	}
	if optionValueString == "" {
		return fmt.Errorf(ui.CommonError1)
	}
	optionValueList := strings.Split(optionValueString, " ")
	command.BindParams = &bindParam{
		option: optionValueList[0],
		value:  optionValueList[1],
	}
	return nil
}

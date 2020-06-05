package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-worker/config"
	"github.com/go-worker/global"
	"github.com/go-worker/regexps"
	"github.com/go-worker/ui"
	"github.com/go-worker/utility"
	"github.com/go-worker/utility3"
)

type Remove struct {
	*CommandStruct
	Params *removeParam
}

func (command *Remove) Execute() error {
	// 解析指令的选项和参数
	parseCommandParamsError := command.parseCommandParams()
	if parseCommandParamsError != nil {
		return parseCommandParamsError
	}

	// 获取全局配置
	projectPath := config.GetCurrentProjectPath()
	fileType := config.WorkerConfig.ProjectSyntax
	if fileType == "" {
		ui.OutputWarnInfo(ui.CommonWarn1)
	}

	if command.Params.optionValue == "" {
		return fmt.Errorf(ui.CommonError1)
	}

	if command.Params.parentValue != "" && command.Params.parentValue != "." {
		projectPath = fmt.Sprintf("%v/%v", projectPath, command.Params.parentValue)
	}

	ignorePath := ""
	if command.Params.ignoreValue != "" {
		ignorePath = fmt.Sprintf("%v/%v", projectPath, command.Params.ignoreValue)
	}

	toRemoveFileList := make([]string, 0)
	switch command.Params.option {
	case "file":
		toRemoveFileList = append(toRemoveFileList, fmt.Sprintf("%v/%v", projectPath, command.Params.optionValue))
	case "type":
		toRemoveFileList = utility3.TraverseDirectorySpecificFile(projectPath, command.Params.optionValue)
	}

	for _, toRemoveFile := range toRemoveFileList {
		if ignorePath == "" || strings.Index(toRemoveFile, ignorePath) == -1 {
			if utility.IsExist(toRemoveFile) {
				removeError := os.Remove(toRemoveFile)
				if removeError != nil {
					ui.OutputWarnInfo(ui.CommonError10, toRemoveFile)
				}
			}
		}
	}

	return nil
}

type removeParam struct {
	option      string
	optionValue string
	parentValue string
	ignoreValue string
}

func (command *Remove) parseCommandParams() error {
	optionValueString := ""
	if optionValueRegexp, hasOptionValueRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AERemoveOptionValue]; hasOptionValueRegexp {
		optionValueString = optionValueRegexp.FindString(command.CommandStruct.InputString)
	} else {
		ui.OutputWarnInfo(ui.CommonWarn2, "remove", "file|type")
	}
	if optionValueString == "" {
		return fmt.Errorf(ui.CommonError1)
	}
	optionValueList := strings.Split(optionValueString, " ")
	command.Params = &removeParam{
		option:      optionValueList[0],
		optionValue: optionValueList[1],
	}
	parentValue := ""
	if parentValueRegexp := regexps.GetRegexpByTemplateEnum(global.OptionParentValueTemplate); parentValueRegexp != nil {
		parentValue = parentValueRegexp.FindString(command.CommandStruct.InputString)
	} else {
		ui.OutputWarnInfo(ui.CommonWarn2, "remove", "parent")
	}
	if parentValue != "" {
		parentValueList := strings.Split(parentValue, " ")
		command.Params.parentValue = parentValueList[1]
	}
	ignoreValue := ""
	if ignoreValueRegexp := regexps.GetRegexpByTemplateEnum(global.OptionIgnoreValueTemplate); ignoreValueRegexp != nil {
		ignoreValue = ignoreValueRegexp.FindString(command.CommandStruct.InputString)
	} else {
		ui.OutputWarnInfo(ui.CommonWarn2, "remove", "ignore")
	}
	if ignoreValue != "" {
		ignoreValueList := strings.Split(ignoreValue, " ")
		command.Params.ignoreValue = ignoreValueList[1]
	}
	return nil
}

package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-worker/config"
	"github.com/go-worker/global"
	"github.com/go-worker/regexps"
	"github.com/go-worker/ui"
	"github.com/go-worker/utility"
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
		ui.OutputWarnInfo(ui.CommonError14)
	}

	if command.Params.optionValue == "" {
		return fmt.Errorf(ui.CommonError1)
	}

	if command.Params.parentValue != "" && command.Params.parentValue != "." {
		projectPath = filepath.Join(projectPath, command.Params.parentValue)
	}

	ignorePath := ""
	if command.Params.ignoreValue != "" {
		ignorePath = filepath.Join(projectPath, command.Params.ignoreValue)
	}

	toRemoveFileList := make([]string, 0)
	switch command.Params.option {
	case "file":
		toRemoveFilePath := filepath.Join(projectPath, command.Params.optionValue)
		if utility.IsExist(toRemoveFilePath) {
			toRemoveFileList = append(toRemoveFileList, toRemoveFilePath)
		} else {
			ui.OutputWarnInfo(ui.CommonError17, toRemoveFilePath)
		}
	case "type":
		toRemoveFileList = utility.TraverseDirectorySpecificFile(projectPath, command.Params.optionValue)
	}

	for _, toRemoveFile := range toRemoveFileList {
		if ignorePath == "" || filepath.Dir(toRemoveFile) != ignorePath {
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
		ui.OutputWarnInfo(ui.CommonError15, "remove", "file|type")
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
		ui.OutputWarnInfo(ui.CommonError15, "remove", "parent")
	}
	if parentValue != "" {
		parentValueList := strings.Split(parentValue, " ")
		command.Params.parentValue = parentValueList[1]
	}
	ignoreValue := ""
	if ignoreValueRegexp := regexps.GetRegexpByTemplateEnum(global.OptionIgnoreValueTemplate); ignoreValueRegexp != nil {
		ignoreValue = ignoreValueRegexp.FindString(command.CommandStruct.InputString)
	} else {
		ui.OutputWarnInfo(ui.CommonError15, "remove", "ignore")
	}
	if ignoreValue != "" {
		ignoreValueList := strings.Split(ignoreValue, " ")
		command.Params.ignoreValue = ignoreValueList[1]
	}
	return nil
}

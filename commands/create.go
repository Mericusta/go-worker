package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-worker/config"
	"github.com/go-worker/global"
	"github.com/go-worker/regexps"
	"github.com/go-worker/ui"
	"github.com/go-worker/utility"
)

type Create struct {
	*CommandStruct
	Params *createParam
}

func (command *Create) Execute() error {
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

	if command.Params.parentValue != "" {
		projectPath = filepath.Join(projectPath, command.Params.parentValue)
	}

	createFilePath, filePackage := command.Params.optionValue, command.Params.optionValue
	switch command.Params.option {
	case "package":
		packagePath := filepath.Join(projectPath, command.Params.optionValue)
		createError := utility.CreateDir(packagePath)
		if createError != nil {
			return fmt.Errorf(ui.CommonError3, createError)
		}
		createFilePath = filepath.Join(packagePath, fmt.Sprintf("%v.%v", command.Params.optionValue, fileType))
	case "file":
		createFilePath = filepath.Join(projectPath, fmt.Sprintf("%v.%v", command.Params.optionValue, fileType))
		filePackage = command.Params.parentValue
		if filePackage == "" {
			filePackage = "main"
		}
	default:
		ui.OutputNoteInfo(ui.CommonNote1, command.Params.option)
		return nil
	}

	file, createFileError := utility.CreateFile(createFilePath)
	defer file.Close()
	if createFileError != nil {
		return fmt.Errorf(ui.CommonError4, createFileError)
	}
	fileContent := fmt.Sprintf("package %v", filePackage)
	file.WriteString(fileContent)

	return nil
}

type createParam struct {
	option      string
	optionValue string
	parentValue string
}

func (command *Create) parseCommandParams() error {
	optionValueString := ""
	if optionValueRegexp, hasOptionValueRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AECreateOptionValue]; hasOptionValueRegexp {
		optionValueString = optionValueRegexp.FindString(command.CommandStruct.InputString)
	} else {
		ui.OutputWarnInfo(ui.CommonError15, "create", "package|file")
	}
	if optionValueString == "" {
		return fmt.Errorf(ui.CommonError1)
	}
	optionValueList := strings.Split(optionValueString, " ")
	command.Params = &createParam{
		option:      optionValueList[0],
		optionValue: optionValueList[1],
	}
	parentValue := ""
	if parentValueRegexp := regexps.GetRegexpByTemplateEnum(global.OptionParentValueTemplate); parentValueRegexp != nil {
		parentValue = parentValueRegexp.FindString(command.CommandStruct.InputString)
	} else {
		ui.OutputWarnInfo(ui.CommonError15, "create", "parent")
	}
	if parentValue != "" {
		parentValueList := strings.Split(parentValue, " ")
		command.Params.parentValue = parentValueList[1]
	}
	return nil
}

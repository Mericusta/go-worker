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
		ui.OutputWarnInfo(ui.CommonWarn1)
	}

	if command.Params.parent != "" {
		projectPath = fmt.Sprintf("%v/%v", projectPath, command.Params.parent)
	}

	createFilePath, filePackage := command.Params.value, command.Params.value
	switch command.Params.option {
	case "package":
		packagePath := fmt.Sprintf("%v/%v", projectPath, command.Params.value)
		createError := utility.CreateDir(packagePath)
		if createError != nil {
			return fmt.Errorf(ui.CommonError3, createError)
		}
		createFilePath = fmt.Sprintf("%v/%v.%v", packagePath, command.Params.value, fileType)
	case "file":
		createFilePath = fmt.Sprintf("%v/%v.%v", projectPath, command.Params.value, fileType)
		utility.TestOutput("command.Params.parent = %v", command.Params.parent)
		filePackage = command.Params.parent
		if filePackage == "" {
			filePackage = "main"
		}
	default:
		ui.OutputNoteInfo(ui.CommonNote1)
		return nil
	}

	utility.TestOutput("create file %v", createFilePath)
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
	option string
	value  string
	parent string
}

func (command *Create) parseCommandParams() error {
	optionValueString := ""
	if optionValueRegexp, hasOptionValueRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AECreateOptionValue]; hasOptionValueRegexp {
		optionValueString = optionValueRegexp.FindString(command.CommandStruct.InputString)
	} else {
		ui.OutputWarnInfo(ui.CommonWarn2, "create", "package|file")
	}
	if optionValueString == "" {
		return fmt.Errorf(ui.CommonError1)
	}
	optionValueList := strings.Split(optionValueString, " ")
	command.Params = &createParam{
		option: optionValueList[0],
		value:  optionValueList[1],
	}
	parentValue := ""
	if parentValueRegexp := regexps.GetRegexpByTemplateEnum(global.OptionParentValueTemplate); parentValueRegexp != nil {
		parentValue = parentValueRegexp.FindString(command.CommandStruct.InputString)
	} else {
		ui.OutputWarnInfo(ui.CommonWarn2, "create", "parent")
	}
	if parentValue != "" {
		parentValueList := strings.Split(parentValue, " ")
		command.Params.parent = parentValueList[1]
	}
	return nil
}

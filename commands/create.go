package commands

import (
	"fmt"

	"github.com/go-worker/config"
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

const (
	createCommandIndex = 0
	createOptionIndex  = 1
	createValueIndex   = 2
	createParentValue  = 4
)

type createParam struct {
	option string
	value  string
	parent string
}

func (command *Create) parseCommandParams() error {
	// optionValueString := ""
	// if optionValueRegexp, hasOptionValueRegexp := regexps.ExpressionEnumRegexpMap[global.AECreateOptionValue]; hasOptionValueRegexp {
	// 	optionValueString = optionValueRegexp.FindString(command.CommandStruct.InputString)
	// } else {
	// 	ui.OutputWarnInfo(ui.CommonWarn2, "create", "option")
	// }
	// if optionValueString == "" {
	// 	return fmt.Errorf(ui.CommonError1)
	// }
	// optionValueList := strings.Split(optionValueString, " ")
	// command.Params = &createParam{
	// 	option: optionValueList[0],
	// 	value:  optionValueList[1],
	// }
	// parentValue := ""
	// if parentValueRegexp, hasParentValueRegexp := regexps.MatchTemplateRegexpMap[global.OptionParentValueTemplate]; hasParentValueRegexp {
	// 	parentValue = parentValueRegexp.FindString(command.CommandStruct.InputString)
	// } else {
	// 	ui.OutputWarnInfo(ui.CommonWarn2, "create", "parent")
	// }
	// utility.TestOutput("parentValue = %v", parentValue)
	// if parentValue != "" {
	// 	parentValueList := strings.Split(parentValue, " ")
	// 	command.Params.parent = parentValueList[1]
	// }
	return nil
}

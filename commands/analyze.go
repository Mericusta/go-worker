package commands

import (
	"github.com/go-worker/utility"
)

type Analyze struct {
	*CommandStruct
	Params *analyzeParam
}

func (command *Analyze) Execute() error {
	utility.TestOutput("analyze execute")
	// 解析指令的选项和参数
	parseCommandParamsError := command.parseCommandParams()
	if parseCommandParamsError != nil {
		return parseCommandParamsError
	}
	return nil
}

const (
	analyzeCommandIndex     = 0
	analyzeSourceTypeIndex  = 1
	analyzeSourceValueIndex = 2
)

type analyzeParam struct {
	sourceType  string
	sourceValue string
	parentValue string
	targetValue string
}

func (command *Analyze) parseCommandParams() error {
	// inputStringList := strings.Split(command.CommandStruct.InputString, " ")
	// if analyzeSourceValueIndex >= len(inputStringList) {
	// 	return fmt.Errorf(ui.CommonError1)
	// }
	// _, inputStringList = utility.SlicePop(inputStringList)
	// pop := func() string {
	// 	var element string
	// 	element, inputStringList = utility.SlicePop(inputStringList)
	// 	return element
	// }
	// parentValueRegexp, hasParentValueRegexp := regexps.ExpressionRegexpMap[global.OptionParentValueExpression]
	// parentValue := ""
	// if hasParentValueRegexp {
	// 	parentValue = parentValueRegexp.FindString(command.CommandStruct.InputString)
	// } else {
	// 	ui.OutputWarnInfo(ui.CommonWarn2, "analyze", "parent")
	// }
	// utility.TestOutput("parent = %v", parentValue)
	// command.Params = &analyzeParam{
	// 	sourceType:  pop(),
	// 	sourceValue: pop(),
	// 	parentValue: parentValue,
	// 	targetValue: func() string {
	// 		if len(inputStringList) > 0 {
	// 			return pop()
	// 		}
	// 		return ""
	// 	}(),
	// }
	return nil
}

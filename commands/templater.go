package commands

import (
	"github.com/PersonalTool/code/go/go_analyzer/utility"
)

type Templater struct {
	*CommandStruct
}

func (command *Templater) Execute() error {
	utility.TestOutput("This is Templater")
	return nil
}

type templaterParam struct {
	value       string
	parentValue string
	outputValue string
}

func (command *Templater) parseCommandParams() error {
	// optionValueString := ""
	// if optionValueRegexp, hasOptionValueRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AERemoveOptionValue]; hasOptionValueRegexp {
	// 	optionValueString = optionValueRegexp.FindString(command.CommandStruct.InputString)
	// } else {
	// 	ui.OutputWarnInfo(ui.CommonError15, "remove", "file|type")
	// }
	// if optionValueString == "" {
	// 	return fmt.Errorf(ui.CommonError1)
	// }
	// optionValueList := strings.Split(optionValueString, " ")
	// command.Params = &templaterParam{
	// 	option:      optionValueList[0],
	// 	optionValue: optionValueList[1],
	// }
	// parentValue := ""
	// if parentValueRegexp := regexps.GetRegexpByTemplateEnum(global.OptionParentValueTemplate); parentValueRegexp != nil {
	// 	parentValue = parentValueRegexp.FindString(command.CommandStruct.InputString)
	// } else {
	// 	ui.OutputWarnInfo(ui.CommonError15, "remove", "parent")
	// }
	// if parentValue != "" {
	// 	parentValueList := strings.Split(parentValue, " ")
	// 	command.Params.parentValue = parentValueList[1]
	// }
	// outputValue := ""
	// if outputValueRegexp := regexps.GetRegexpByTemplateEnum(global.OptionOutputValueTemplate); outputValueRegexp != nil {
	// 	outputValue = outputValueRegexp.FindString(command.CommandStruct.InputString)
	// } else {
	// 	ui.OutputWarnInfo(ui.CommonError15, "analyze", "output")
	// }
	// if outputValue != "" {
	// 	outputValueList := strings.Split(outputValue, " ")
	// 	command.Params.outputValue = outputValueList[1]
	// }
	return nil
}

package commands

import (
	"fmt"
)

type Spider struct {
	*CommandStruct
	Params *spiderParam
}

func (command *Spider) Execute() error {
	// 解析指令的选项和参数
	parseCommandParamsError := command.parseCommandParams()
	if parseCommandParamsError != nil {
		return parseCommandParamsError
	}
	fmt.Println("params = %+v", command.Params)
	return nil
}

type spiderParam struct {
	option      string
	optionValue string
}

func (command *Spider) parseCommandParams() error {
	// optionValue := ""
	// if optionValueRegexp := regexps.GetRegexpByTemplateEnum(global.OptionWebAddressValueTemplate); optionValueRegexp != nil {
	// 	optionValue = optionValueRegexp.FindString(command.CommandStruct.InputString)
	// } else {
	// 	ui.OutputWarnInfo(ui.CommonWarn2, "spider", "website|page")
	// }
	// if optionValue != "" {
	// 	return fmt.Errorf(ui.CommonError1)
	// }
	// optionValueList := strings.Split(optionValue, " ")
	// command.Params = &spiderParam{
	// 	option:      optionValueList[0],
	// 	optionValue: optionValueList[1],
	// }
	return nil
}

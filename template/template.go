package template

import "github.com/go-worker/global"

// TemplateString 模板字符串类型
type TemplateString string

// TemplateExpression 模板表达式类型
type TemplateExpression string

// TemplateEnumExpressionMap 模板与模板表达式映射
var TemplateEnumExpressionMap map[global.TemplateEnum]TemplateExpression

// CommandTemplateExpressionMap 指令与模板表达式映射
var CommandTemplateExpressionMap map[global.CommandEnum]TemplateExpression

func init() {
	TemplateEnumExpressionMap = map[global.TemplateEnum]TemplateExpression{
		global.OptionParentValueTemplate: TEOptionParentValue,
		global.OptionOutputValueTemplate: TEOptionOutputValue,
	}
	CommandTemplateExpressionMap = map[global.CommandEnum]TemplateExpression{
		global.CmdCreate:  TECreate,
		global.CmdConvert: TEConvert,
		global.CmdAnalyze: TEAnalyze,
	}
}

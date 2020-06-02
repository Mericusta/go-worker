package template

import "github.com/go-worker/global"

// TemplateString 模板字符串类型
type TemplateString string

// TemplateExpression 模板表达式类型
type TemplateExpression string

// CommandTemplateExpressionMap 指令枚举与模板表达式的映射
var CommandTemplateExpressionMap map[global.CommandEnum]TemplateExpression

// TemplateExpressionMap 模板枚举与模板表达式的映射
var TemplateExpressionMap map[global.TemplateEnum]TemplateExpression

func init() {
	CommandTemplateExpressionMap = map[global.CommandEnum]TemplateExpression{
		global.CmdBind:    TECommandBind,
		global.CmdCreate:  TECommandCreate,
		global.CmdConvert: TECommandConvert,
		global.CmdAnalyze: TECommandAnalyze,
	}
	TemplateExpressionMap = map[global.TemplateEnum]TemplateExpression{
		global.OptionParentValueTemplate:    TEOptionParentValue,
		global.OptionOutputValueTemplate:    TEOptionOutputValue,
		global.GoKeywordImportValueTemplate: TEGoKeywordImportValue,
		global.GoFunctionDefinitionTemplate: TEGoFunctionDefinition,
		global.GoFunctionCallTemplate:       TEGoFunctionCall,
		global.CommandBindTemplate:          TECommandBind,
		global.CommandCreateTemplate:        TECommandCreate,
		global.CommandConverTemplate:        TECommandConvert,
	}
}

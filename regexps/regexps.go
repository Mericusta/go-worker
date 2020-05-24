package regexps

import (
	"regexp"

	"github.com/go-worker/global"
	"github.com/go-worker/template"
	"github.com/go-worker/ui"
)

// Expression 正则表达式类型
type Expression string

// templateEnumAtomicExpressionMap 模板与原子表达式映射
var templateEnumAtomicExpressionMap map[global.TemplateEnum]Expression

// matchTemplateExpressionMap 模板与模板匹配式映射
var matchTemplateExpressionMap map[global.TemplateEnum]Expression

// MatchTemplateRegexpMap 模板与模板匹配解析式映射
var MatchTemplateRegexpMap map[global.TemplateEnum]*regexp.Regexp

// commandExpressionMap 指令枚举与正则表达式映射
var commandExpressionMap map[global.CommandEnum]string

// RegexpCommandEnumMap 解析式与指令映射
var RegexpCommandEnumMap map[*regexp.Regexp]global.CommandEnum

// // LogicRegexpMap 内部逻辑表达式枚举
// var LogicRegexpMap map[string]*regexp.Regexp

func init() {
	MatchTemplateRegexpMap = make(map[global.TemplateEnum]*regexp.Regexp)
	commandExpressionMap = make(map[global.CommandEnum]string)
	RegexpCommandEnumMap = make(map[*regexp.Regexp]global.CommandEnum)

	// 注册模板与原子表达式
	registTemplateEnumAtomicExpressionMap()

	// 注册模板与模板匹配式
	registTemplateMatchTemplateExpression()

	// 编译模板匹配式
	complieMatchTemplateRegexp()

	// 解析指令的模板表达式
	parseCommandTemplateExpression()

	// 注册指令的原子表达式
	registCommandAtomicExpression()

	// 编译指令正则表达式
	complieCommandRegexp()
}

func registTemplateEnumAtomicExpressionMap() {
	templateEnumAtomicExpressionMap = map[global.TemplateEnum]Expression{
		global.PathTemplate: AEPath,
	}
}

func registTemplateMatchTemplateExpression() {
	matchTemplateExpressionMap = map[global.TemplateEnum]Expression{
		global.PathTemplate:              MTEPath,
		global.OptionParentValueTemplate: MTEOptionParentValue,
		global.OptionOutputValueTemplate: MTEOptionOutputValue,
	}
}

func complieMatchTemplateRegexp() {
	for templateEnum, matchTemplateExpression := range matchTemplateExpressionMap {
		matchTemplateRegexp := regexp.MustCompile(string(matchTemplateExpression))
		if matchTemplateRegexp == nil {
			ui.OutputErrorInfo("complie match template %v expression %v, but get nil", templateEnum, matchTemplateExpression)
			continue
		}
		MatchTemplateRegexpMap[templateEnum] = matchTemplateRegexp
	}
}

func parseCommandTemplateExpression() {
	for commandEnum, templateExpression := range template.CommandTemplateExpressionMap {
		replaceExpression := string(templateExpression)
		index := 0
		for {
			hasTemplate := false
			for templateEnum, matchTemplateRegexp := range MatchTemplateRegexpMap {
				replaceString := ""
				if expression, hasExpression := template.TemplateEnumExpressionMap[templateEnum]; hasExpression {
					replaceString = string(expression)
				} else if atomicExpression, hasAtomicExpression := templateEnumAtomicExpressionMap[templateEnum]; hasAtomicExpression {
					replaceString = string(atomicExpression)
				}
				if matchTemplateRegexp.MatchString(replaceExpression) && replaceString != "" {
					replaceExpression = matchTemplateRegexp.ReplaceAllString(replaceExpression, replaceString)
					hasTemplate = true
					break
				}
			}
			if !hasTemplate {
				break
			}
			index++
		}
		commandExpressionMap[commandEnum] = replaceExpression
	}
}

func registCommandAtomicExpression() {
	commandExpressionMap[global.CmdExit] = string(AECmdExit)
	commandExpressionMap[global.CmdBind] = string(AECmdBind)
}

func complieCommandRegexp() {
	for commandEnum, commandExpression := range commandExpressionMap {
		commandRegexp := regexp.MustCompile(string(commandExpression))
		if commandRegexp == nil {
			ui.OutputErrorInfo("complie command %v expression %v, but get nil", commandEnum, commandExpression)
			continue
		}
		RegexpCommandEnumMap[commandRegexp] = commandEnum
	}
}

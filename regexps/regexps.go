package regexps

import (
	"regexp"

	"github.com/go-worker/global"
	"github.com/go-worker/template"
	"github.com/go-worker/ui"
	"github.com/go-worker/utility"
)

// Expression 正则表达式类型
type Expression string

// expressionEnumAtomicExpressionMap 正则表达式枚举与原子表达式映射
var expressionEnumAtomicExpressionMap map[global.RegexpEnum]Expression

// ExpressionEnumRegexoMap 正则表达式枚举与解析式映射
var ExpressionEnumRegexoMap map[global.RegexpEnum]*regexp.Regexp

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

func init() {
	ExpressionEnumRegexoMap = make(map[global.RegexpEnum]*regexp.Regexp)
	MatchTemplateRegexpMap = make(map[global.TemplateEnum]*regexp.Regexp)
	commandExpressionMap = make(map[global.CommandEnum]string)
	RegexpCommandEnumMap = make(map[*regexp.Regexp]global.CommandEnum)

	// 注册原子表达式
	registAtomicExpression()

	// 编译原子表达式的解析式
	complieAtomicRegexp()

	// 注册模板与原子表达式
	registTemplateEnumAtomicExpressionMap()

	// 注册模板与模板匹配式
	registTemplateMatchTemplateExpression()

	// 编译模板匹配式的解析式
	complieMatchTemplateRegexp()

	// 解析指令的模板表达式
	parseCommandTemplateExpression()

	// 注册指令的原子表达式
	registCommandAtomicExpression()

	// 编译指令正则表达式的解析式
	complieCommandRegexp()
}

func registAtomicExpression() {
	expressionEnumAtomicExpressionMap = map[global.RegexpEnum]Expression{
		global.AETemplateCommonKeyword: AETemplateCommonKeyword,
		global.AEPath:                  AEPath,
		global.AEBindOptionValue:       AEBindOptionValue,
		global.AECreateOptionValue:     AECreateOptionValue,
		global.AEConvertOptionValue:    AEConvertOptionValue,
	}
}

func complieAtomicRegexp() {
	for expressionEnum, atomicExpression := range expressionEnumAtomicExpressionMap {
		atomicRegexp := regexp.MustCompile(string(atomicExpression))
		if atomicRegexp == nil {
			ui.OutputWarnInfo("complie atomic expression %v %v, but get nil", expressionEnum, atomicExpression)
		}
		ExpressionEnumRegexoMap[expressionEnum] = atomicRegexp
	}
}

func registTemplateEnumAtomicExpressionMap() {
	templateEnumAtomicExpressionMap = map[global.TemplateEnum]Expression{
		global.PathTemplate:               AEPath,
		global.BindOptionValueTemplate:    AEBindOptionValue,
		global.CreateOptionValueTemplate:  AECreateOptionValue,
		global.ConvertOptionValueTemplate: AEConvertOptionValue,
	}
}

func registTemplateMatchTemplateExpression() {
	matchTemplateExpressionMap = map[global.TemplateEnum]Expression{
		global.PathTemplate:               MTEPath,
		global.OptionParentValueTemplate:  MTEOptionParentValue,
		global.OptionOutputValueTemplate:  MTEOptionOutputValue,
		global.BindOptionValueTemplate:    MTEBindOption,
		global.CreateOptionValueTemplate:  MTECreateOption,
		global.ConvertOptionValueTemplate: MTEConvertOption,
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
	for commandEnum, expression := range commandExpressionMap {
		utility.TestOutput("commandEnum = %v, expression = %v", commandEnum, expression)
	}
}

func registCommandAtomicExpression() {
	commandExpressionMap[global.CmdExit] = string(AECmdExit)
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

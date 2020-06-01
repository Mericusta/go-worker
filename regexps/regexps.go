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

// atomicExpressionMap 原子表达式枚举类型与原子表达式的映射
var atomicExpressionMap map[global.AtomicExpressionEnum]AtomicExpression

// AtomicExpressionEnumRegexpMap 原子表达式枚举类型与原子表达式的解析式的映射
var AtomicExpressionEnumRegexpMap map[global.AtomicExpressionEnum]*regexp.Regexp

// templateKeywordRegexpMap 模板关键词与模板关键词匹配式的映射
var templateKeywordRegexpMap map[template.TemplateKeyword]*regexp.Regexp

// templateKeywordReplaceStringMap 模板关键词与替换文本的映射
var templateKeywordReplaceStringMap map[template.TemplateKeyword]string

// commandAtomicExpressionMap 指令枚举与原子表达式的映射
var commandAtomicExpressionMap map[global.CommandEnum]AtomicExpression

// commandTemplateExpressionMap 指令枚举与模板表达式解析之后的正则表达式的映射
var commandTemplateExpressionMap map[global.CommandEnum]Expression

// RegexpCommandEnumMap 指令枚举与指令匹配式的映射
var RegexpCommandEnumMap map[*regexp.Regexp]global.CommandEnum

func init() {
	AtomicExpressionEnumRegexpMap = make(map[global.AtomicExpressionEnum]*regexp.Regexp)
	templateKeywordRegexpMap = make(map[template.TemplateKeyword]*regexp.Regexp)
	commandTemplateExpressionMap = make(map[global.CommandEnum]Expression)
	RegexpCommandEnumMap = make(map[*regexp.Regexp]global.CommandEnum)

	// 注册原子表达式
	registAtomicExpression()

	// 编译原子表达式的解析式
	complieAtomicRegexp()

	// 注册模板关键词的替换文本
	registTemplateKeywordReplaceString()

	// 编译模板关键词匹配式的解析式
	complieTemplateKeywordRegexp()

	// 解析指令的模板表达式
	parseTemplateExpressionCommand()

	// 注册指令的原子表达式
	regisCommandAtomicExpression()

	// 编译指令的解析式
	complieCommandRegexp()
}

func registAtomicExpression() {
	atomicExpressionMap = map[global.AtomicExpressionEnum]AtomicExpression{
		global.AETemplateCommonKeyword: AETemplateCommonKeyword,
		global.AEPath:                  AEPath,
		global.AEDoubleQuotesContent:   AEDoubleQuotesContent,
		global.AEBracketsContent:       AEBracketsContent,
		global.AESquareBracketsContent: AESquareBracketsContent,
		global.AECurlyBracesContent:    AECurlyBracesContent,
		global.AEBindOptionValue:       AEBindOptionValue,
		global.AECreateOptionValue:     AECreateOptionValue,
		global.AEConvertOptionValue:    AEConvertOptionValue,
		global.AEConvertACOptionValue:  AEConvertACOptionValue,
		global.AEAnalyzeOptionValue:    AEAnalyzeOptionValue,
		global.AEGoKeywordPackageValue: AEGoKeywordPackageValue,
		global.AETemplateStyle:         AETemplateStyle,
	}
}

func complieAtomicRegexp() {
	for atomicExpressionEnum, atomicExpression := range atomicExpressionMap {
		atomicRegexp := regexp.MustCompile(string(atomicExpression))
		if atomicRegexp == nil {
			ui.OutputWarnInfo("complie atomic expression[%v], but get nil", atomicExpression)
			continue
		}
		AtomicExpressionEnumRegexpMap[atomicExpressionEnum] = atomicRegexp
	}
}

func registTemplateKeywordReplaceString() {
	templateKeywordReplaceStringMap = map[template.TemplateKeyword]string{
		template.TK_Path:                      string(AEPath),
		template.TK_OPVExpression:             string(template.TEOptionParentValue),
		template.TK_OOVExpression:             string(template.TEOptionOutputValue),
		template.TK_DoubleQuotesContent:       string(AEDoubleQuotesContent),
		template.TK_BOVExpression:             string(AEBindOptionValue),
		template.TK_CreateOVExpression:        string(AECreateOptionValue),
		template.TK_ConvertOVExpression:       string(AEConvertOptionValue),
		template.TK_ConvertACOptionExpression: string(AEConvertACOptionValue),
		template.TK_AnalyzeOVExpression:       string(AEAnalyzeOptionValue),
	}
}

func complieTemplateKeywordRegexp() {
	for TemplateKeyword := range templateKeywordReplaceStringMap {
		if templateKeywordRegexp := regexp.MustCompile(string(TemplateKeyword)); templateKeywordRegexp != nil {
			templateKeywordRegexpMap[TemplateKeyword] = templateKeywordRegexp
		} else {
			ui.OutputErrorInfo("template keyword %v cannot complie to regexp", TemplateKeyword)
		}
	}
}

func parseTemplateExpressionCommand() {
	templateCommonKeywordRegexp, hasTemplateCommonKeywordRegexp := AtomicExpressionEnumRegexpMap[global.AETemplateCommonKeyword]
	if !hasTemplateCommonKeywordRegexp {
		ui.OutputErrorInfo("template common keyword regexp does not exist")
		return
	}
	for commandEnum, TECommand := range template.CommandTemplateExpressionMap {
		parsedCommand := parseTemplateExpression(templateCommonKeywordRegexp, TECommand)
		if parsedCommand == "" {
			ui.OutputWarnInfo("parse command[%v] template expression[%v] but get empty", commandEnum, TECommand)
			continue
		}
		commandTemplateExpressionMap[commandEnum] = parsedCommand
	}
}

func regisCommandAtomicExpression() {
	commandAtomicExpressionMap = map[global.CommandEnum]AtomicExpression{
		global.CmdExit: AECommandExit,
	}
}

func complieCommandRegexp() {
	// 编译模板表达式解析之后的原子表达式指令
	for commandEnum, commandExpression := range commandTemplateExpressionMap {
		commandRegexp := regexp.MustCompile(string(commandExpression))
		if commandRegexp == nil {
			ui.OutputWarnInfo("complie command[%v] by expression[%v], but get nil", commandEnum, commandExpression)
		}
		RegexpCommandEnumMap[commandRegexp] = commandEnum
	}
	// 编译纯原子表达式指令
	for commandEnum, commandExpression := range commandAtomicExpressionMap {
		commandRegexp := regexp.MustCompile(string(commandExpression))
		if commandRegexp == nil {
			ui.OutputWarnInfo("complie command[%v] by expression[%v], but get nil", commandEnum, commandExpression)
		}
		RegexpCommandEnumMap[commandRegexp] = commandEnum
	}
}

// GetRegexpByTemplateEnum 根据模板枚举获得其原子解析式
func GetRegexpByTemplateEnum(templateEnum global.TemplateEnum) *regexp.Regexp {
	templateExpression, hasTemplateExpression := template.TemplateExpressionMap[templateEnum]
	if !hasTemplateExpression {
		ui.OutputErrorInfo("template expression[%v] does not regist expression", templateEnum)
		return nil
	}
	templateCommonKeywordRegexp, hasTemplateCommonKeywordRegexp := AtomicExpressionEnumRegexpMap[global.AETemplateCommonKeyword]
	if !hasTemplateCommonKeywordRegexp {
		ui.OutputErrorInfo("template common keyword regexp does not exist")
		return nil
	}
	parsedCommand := parseTemplateExpression(templateCommonKeywordRegexp, templateExpression)
	if parsedCommand == "" {
		ui.OutputWarnInfo("parse template expression[%v] but get empty", templateExpression)
		return nil
	}
	// utility.TestOutput("templateExpression = %v", parsedCommand)
	templateExpressionRegexp := regexp.MustCompile(string(parsedCommand))
	if templateExpressionRegexp == nil {
		ui.OutputWarnInfo("complie template expression[%v], but get nil", parsedCommand)
	}
	return templateExpressionRegexp
}

func parseTemplateExpression(templateCommonKeywordRegexp *regexp.Regexp, templateExpression template.TemplateExpression) Expression {
	// 查找模板关键词
	templateKeywordExpressionList := templateCommonKeywordRegexp.FindAllString(string(templateExpression), -1)
	utility.TestOutput("templateExpression = %v, templateKeywordExpressionList = %v", templateExpression, templateKeywordExpressionList)
	for len(templateKeywordExpressionList) != 0 {
		for _, templateKeywordExpression := range templateKeywordExpressionList {
			// utility.TestOutput("templateKeywordExpression = %v", templateKeywordExpression)
			if templateKeywordRegexp, hasKeywordRegexp := templateKeywordRegexpMap[template.TemplateKeyword(templateKeywordExpression)]; hasKeywordRegexp {
				if toReplaceString, hasToReplaceString := templateKeywordReplaceStringMap[template.TemplateKeyword(templateKeywordExpression)]; hasToReplaceString {
					// 替换模板关键词
					// utility.TestOutput("to replace %v", toReplaceString)
					templateExpression = template.TemplateExpression(templateKeywordRegexp.ReplaceAllString(string(templateExpression), toReplaceString))
				} else {
					// utility.TestOutput("%v does not have to replace string", templateKeywordExpression)
					return ""
				}
			} else {
				// utility.TestOutput("%v does not have regexp", templateKeywordExpression)
				return ""
			}
		}

		utility.TestOutput("after replace, templateExpression = %v", templateExpression)

		// 查找模板关键词
		templateKeywordExpressionList = templateCommonKeywordRegexp.FindAllString(string(templateExpression), -1)
		// utility.TestOutput("templateKeywordExpressionList = %v", templateKeywordExpressionList)
	}
	return Expression(templateExpression)
}

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
	templateKeywordRegexpMap = make(map[template.TemplateKeyword]*regexp.Regexp)
	commandTemplateExpressionMap = make(map[global.CommandEnum]Expression)
	RegexpCommandEnumMap = make(map[*regexp.Regexp]global.CommandEnum)

	// 注册模板关键词的替换文本
	registTemplateKeywordReplaceString()

	// 注册指令的原子表达式
	regisCommandAtomicExpression()

	// 模板关键词匹配式
	complieTemplateKeywordRegexp()

	// 编译通用模板关键词匹配式
	complieTemplateCommonKeywordRegexp()
	if templateCommonKeywordRegexp == nil {
		return
	}

	// 解析指令的模板表达式
	parseTemplateExpressionCommand()

	// 编译指令的解析式
	complieCommandRegexp()
}

func registTemplateKeywordReplaceString() {
	templateKeywordReplaceStringMap = map[template.TemplateKeyword]string{
		template.TK_BOVExpression:       string(AEBindOptionValue),
		template.TK_CreateOVExpression:  string(AECreateOptionValue),
		template.TK_ConvertOVExpression: string(AEConvertOptionValue),
		template.TK_Path:                string(AEPath),
		template.TK_OPVExpression:       string(template.TEOptionParentValue),
		template.TK_OOVExpression:       string(template.TEOptionOutputValue),
	}
}

func regisCommandAtomicExpression() {
	commandAtomicExpressionMap = map[global.CommandEnum]AtomicExpression{
		global.CmdExit: AECommandExit,
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

func complieTemplateCommonKeywordRegexp() {
	templateCommonKeywordRegexp = regexp.MustCompile(string(AETemplateCommonKeyword))
	if templateCommonKeywordRegexp == nil {
		ui.OutputErrorInfo("complie template keyword common regexp, but get nil")
	}
}

func parseTemplateExpressionCommand() {
	for commandEnum, TECommand := range template.CommandTemplateExpressionMap {
		// 查找模板关键词
		templateKeywordExpressionList := templateCommonKeywordRegexp.FindAllString(string(TECommand), -1)
		utility.TestOutput("commandEnum %v, TECommand = %v, templateKeywordExpressionList = %v", commandEnum, TECommand, templateKeywordExpressionList)
		for len(templateKeywordExpressionList) != 0 {
			for _, templateKeywordExpression := range templateKeywordExpressionList {
				utility.TestOutput("templateKeywordExpression = %v", templateKeywordExpression)
				if templateKeywordRegexp, hasKeywordRegexp := templateKeywordRegexpMap[template.TemplateKeyword(templateKeywordExpression)]; hasKeywordRegexp {
					if toReplaceString, hasToReplaceString := templateKeywordReplaceStringMap[template.TemplateKeyword(templateKeywordExpression)]; hasToReplaceString {
						// 替换模板关键词
						utility.TestOutput("to replace %v", toReplaceString)
						TECommand = template.TemplateExpression(templateKeywordRegexp.ReplaceAllString(string(TECommand), toReplaceString))
					} else {
						utility.TestOutput("%v does not have to replace string", templateKeywordExpression)
						return
					}
				} else {
					utility.TestOutput("%v does not have regexp", templateKeywordExpression)
					return
				}
			}

			utility.TestOutput("after replace, TECommand = %v", TECommand)

			// 查找模板关键词
			templateKeywordExpressionList = templateCommonKeywordRegexp.FindAllString(string(TECommand), -1)
			utility.TestOutput("templateKeywordExpressionList = %v", templateKeywordExpressionList)
		}
		commandTemplateExpressionMap[commandEnum] = Expression(TECommand)
	}
}

func complieCommandRegexp() {
	for commandEnum, commandExpression := range commandAtomicExpressionMap {
		commandRegexp := regexp.MustCompile(string(commandExpression))
		if commandRegexp == nil {
			ui.OutputWarnInfo("complie command[%v] by expression[%v], but get nil", commandEnum, commandExpression)
		}
		RegexpCommandEnumMap[commandRegexp] = commandEnum
	}
	for commandEnum, commandExpression := range commandTemplateExpressionMap {
		commandRegexp := regexp.MustCompile(string(commandExpression))
		if commandRegexp == nil {
			ui.OutputWarnInfo("complie command[%v] by expression[%v], but get nil", commandEnum, commandExpression)
		}
		RegexpCommandEnumMap[commandRegexp] = commandEnum
	}
}

package regexps

import (
	"regexp"

	"github.com/go-worker/global"
)

// Expression 正则表达式类型
type Expression string

// CmdExpressionList 指令正则表达式列表
var CmdExpressionList []Expression

// LogicExpressionList 逻辑正则表达式列表
var LogicExpressionList map[string]Expression

// ExpressionRegexpMap 正则表达式和解析式
var ExpressionRegexpMap map[string]*regexp.Regexp

func init() {
	CmdExpressionList = []Expression{
		ExpExit,
		ExpBind,
		ExpCreate,
		ExpConvert,
	}
	LogicExpressionList = map[string]Expression{
		global.TemplateKeywordExpression: ExpTemplateKeyword,
	}
	ExpressionRegexpMap = make(map[string]*regexp.Regexp)
	for _, cmdExp := range CmdExpressionList {
		ExpressionRegexpMap[string(cmdExp)] = regexp.MustCompile(string(cmdExp))
	}
	for key, logicExp := range LogicExpressionList {
		ExpressionRegexpMap[key] = regexp.MustCompile(string(logicExp))
	}
}

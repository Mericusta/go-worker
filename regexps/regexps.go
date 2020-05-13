package regexps

import (
	"regexp"
)

// Expression 正则表达式类型
type Expression string

// CmdExpressionList 正则表达式列表
var CmdExpressionList []Expression

// ExpressionRegexpMap 正则表达式和解析式
var ExpressionRegexpMap map[Expression]*regexp.Regexp

func init() {
	CmdExpressionList = []Expression{
		Exit,
		Create,
		Spider,
		Worker,
	}
	ExpressionRegexpMap = make(map[Expression]*regexp.Regexp)
	for _, cmdExp := range CmdExpressionList {
		ExpressionRegexpMap[cmdExp] = regexp.MustCompile(string(cmdExp))
	}
}

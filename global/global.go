package global

import "github.com/go-worker/fsm"

var FsmState fsm.FSMState

var SyntaxFileSuffixMap map[string]string

func init() {
	SyntaxFileSuffixMap = make(map[string]string)
	SyntaxFileSuffixMap[SyntaxGo] = "go"
	SyntaxFileSuffixMap[SyntaxCSV] = "csv"
}

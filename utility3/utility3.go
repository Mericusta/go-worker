package utility3

import (
	"github.com/go-worker/global"
	"github.com/go-worker/regexps"
	"github.com/go-worker/ui"
)

// TrimSpaceLine 移除空白行
func TrimSpaceLine(content string) string {
	spaceLineRegexp, hasSpaceLineRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AESpaceLine]
	if !hasSpaceLineRegexp {
		ui.OutputWarnInfo(ui.CommonWarn3, global.AESpaceLine)
		return ""
	}
	replaceContent := spaceLineRegexp.ReplaceAllString(content, "\n")
	return replaceContent
}

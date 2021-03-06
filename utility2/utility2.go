package utility2

import (
	"fmt"
	"strings"

	"github.com/go-worker/global"
)

// FormatOutput 格式化输出
func FormatOutput(logMark, format string, content ...interface{}) {
	formatContent := fmt.Sprintf(format, content...)
	fmt.Printf("%v: %v\n", logMark, formatContent)
}

// TestOutput 测试输出
func TestOutput(format string, content ...interface{}) {
	FormatOutput(global.LogMarkTest, format, content...)
}

// GetPunctuationMark 获取标点符号
func GetPunctuationMark(punctuationMark int) (rune, rune) {
	switch punctuationMark {
	case global.PunctuationMarkCurlyBracket:
		return global.PunctuationMarkLeftCurlyBracket, global.PunctuationMarkRightCurlyBracket
	case global.PunctuationMarkSquareBracket:
		return global.PunctuationMarkLeftSquareBracket, global.PunctuationMarkRightSquareBracket
	default:
		return global.PunctuationMarkLeftBracket, global.PunctuationMarkRightBracket
	}
}

// FormatFilePathWithOS 根据操作系统格式化路径
func FormatFilePathWithOS(filePath string) string {
	beReplaced := "/"
	toReplace := "\\"
	if global.OperationSystem == global.OSLinux {
		beReplaced, toReplace = toReplace, beReplaced
	}
	return strings.ReplaceAll(filePath, beReplaced, toReplace)
}

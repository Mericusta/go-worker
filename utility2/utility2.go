package utility2

import (
	"fmt"

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

func GetPunctuationMark(punctuationMark int) (rune, rune) {
	switch punctuationMark {
	case global.PunctuationMarkBracket:
		return global.PunctuationMarkLeftBracket, global.PunctuationMarkRightBracket
	case global.PunctuationMarkCurlyBracket:
		return global.PunctuationMarkLeftCurlyBracket, global.PunctuationMarkRightCurlyBracket
	default:
		return global.PunctuationMarkLeftQuote, global.PunctuationMarkRightQuote
	}
}

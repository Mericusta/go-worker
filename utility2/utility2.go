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

// CalculatePunctuationMarksContentLength 计算成对标点符号的内容的长度
func CalculatePunctuationMarksContentLength(afterLeftContent string, punctuationMark int) int {
	leftCount := 1
	rightCount := 0
	leftPunctionMark := global.PunctuationMarkLeftQuote
	rightPunctionMark := global.PunctuationMarkRightQuote
	switch punctuationMark {
	case global.PunctuationMarkBracket:
		leftPunctionMark = global.PunctuationMarkLeftBracket
		rightPunctionMark = global.PunctuationMarkRightBracket
	case global.PunctuationMarkCurlyBraces:
		leftPunctionMark = global.PunctuationMarkLeftCurlyBraces
		rightPunctionMark = global.PunctuationMarkRightCurlyBraces
	}
	return strings.IndexFunc(afterLeftContent, func(r rune) bool {
		if r == leftPunctionMark {
			leftCount++
		} else if r == rightPunctionMark {
			rightCount++
		}
		return leftCount == rightCount
	})
}

// FixBracketMatchingResult 修正贪婪括号匹配的错误 type(x)( -> type(x)
func FixBracketMatchingResult(content string) string {
	leftBracketIndex := strings.Index(content, string(global.PunctuationMarkLeftBracket))
	if leftBracketIndex == -1 {
		return content
	}
	bracketContentLength := CalculatePunctuationMarksContentLength(content[leftBracketIndex+1:], global.PunctuationMarkBracket)
	// first 1 is because from 0
	// last 1 is right bracket
	return content[:leftBracketIndex+1+bracketContentLength+1]
}

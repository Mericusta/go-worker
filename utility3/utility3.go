package utility3

import (
	"strings"

	"github.com/go-worker/global"
	"github.com/go-worker/regexps"
	"github.com/go-worker/ui"
	"github.com/go-worker/utility"
	"github.com/go-worker/utility2"
)

// TrimSpaceLine 移除空白行
func TrimSpaceLine(content string) string {
	spaceLineRegexp, hasSpaceLineRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AESpaceLine]
	if !hasSpaceLineRegexp {
		ui.OutputWarnInfo(ui.CommonError16, global.AESpaceLine)
		return ""
	}
	replaceContent := spaceLineRegexp.ReplaceAllString(content, "\n")
	return replaceContent
}

// CalculatePunctuationMarksContentLength 计算成对标点符号的内容的长度
func CalculatePunctuationMarksContentLength(afterLeftContent string, punctuationMark int) int {
	leftCount := 1
	rightCount := 0
	leftPunctuationMark, rightPunctuationMark := utility2.GetPunctuationMark(punctuationMark)
	return strings.IndexFunc(afterLeftContent, func(r rune) bool {
		if r == leftPunctuationMark {
			leftCount++
		} else if r == rightPunctuationMark {
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

// TraitPunctuationMarksContent 成对标点符号的内容提取
func TraitPunctuationMarksContent(content string, punctuationMark int) *utility.PunctuationContent {
	leftPunctuationMark, rightPunctuationMark := utility2.GetPunctuationMark(punctuationMark)
	return utility.RecursiveTraitPunctuationContent(content, leftPunctuationMark, rightPunctuationMark)
}

// TraitMultiPunctuationMarksContent 混合成对标点符号的内容分类提取
func TraitMultiPunctuationMarksContent(content string, punctuationMarkList []int, maxDeep int) *utility.NewPunctuationContent {
	leftPunctuationMarkList := make([]rune, 0, len(punctuationMarkList))
	for _, punctuationMark := range punctuationMarkList {
		leftPunctuationMark, _ := utility2.GetPunctuationMark(punctuationMark)
		leftPunctuationMarkList = append(leftPunctuationMarkList, leftPunctuationMark)
	}
	return utility.RecursiveTraitMultiPunctuationMarksContent(content, 0, 0, leftPunctuationMarkList, maxDeep, 0)
}

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

// SplitContent 划分内容节点
type SplitContent struct {
	ContentList         []string
	SubSplitContentList []*SplitContent
}

// RecursiveSplitUnderSameDeepPunctuationMarksContent 相同深度的成对标点符号下的内容划分
func RecursiveSplitUnderSameDeepPunctuationMarksContent(content string, punctuationMarkList []int, splitter string, maxDeep int) *SplitContent {
	if punctuationContentNode := TraitMultiPunctuationMarksContent(content, punctuationMarkList, maxDeep); punctuationContentNode != nil {
		return splitUnderSameDeepPunctuationMarksContent(punctuationContentNode, splitter, maxDeep, 0)
	}
	return nil
}

// RecursiveSplitUnderSameDeepPunctuationMarksContentNode 相同深度的成对标点符号下的内容划分
func RecursiveSplitUnderSameDeepPunctuationMarksContentNode(punctuationContentNode *utility.NewPunctuationContent, splitter string, maxDeep int) *SplitContent {
	return splitUnderSameDeepPunctuationMarksContent(punctuationContentNode, splitter, maxDeep, 0)
}

func splitUnderSameDeepPunctuationMarksContent(punctuationContentNode *utility.NewPunctuationContent, splitter string, maxDeep, deep int) *SplitContent {
	splitContentNode := &SplitContent{
		ContentList:         make([]string, 0),
		SubSplitContentList: make([]*SplitContent, 0),
	}

	var offset int
	var leftIndex int
	cycle := 0
	maxCycle := len(strings.Split(punctuationContentNode.Content, splitter))
	for cycle != maxCycle {
		cycle++
		length := strings.Index(punctuationContentNode.Content[leftIndex+offset:], splitter)
		if length == -1 {
			splitContentNode.ContentList = append(splitContentNode.ContentList, punctuationContentNode.Content[leftIndex:])
			break
		}
		rightIndex := leftIndex + length + offset
		inner := false
		for subIndex := 0; subIndex != len(punctuationContentNode.SubPunctuationContentList); subIndex++ {
			if punctuationContentNode.SubPunctuationIndexMap[subIndex][0] <= rightIndex && rightIndex <= punctuationContentNode.SubPunctuationIndexMap[subIndex][1] {
				inner = true
				offset = punctuationContentNode.SubPunctuationIndexMap[subIndex][1] - leftIndex + 1
				break
			}
		}
		if inner {
			continue
		}
		splitContentNode.ContentList = append(splitContentNode.ContentList, punctuationContentNode.Content[leftIndex:rightIndex])
		offset = 0
		leftIndex = rightIndex + len(splitter)
	}

	if deep == maxDeep {
		return splitContentNode
	}

	for _, subPuncutationContentNode := range punctuationContentNode.SubPunctuationContentList {
		if len(subPuncutationContentNode.Content) != 0 {
			if subSplitContentNode := splitUnderSameDeepPunctuationMarksContent(subPuncutationContentNode, splitter, maxDeep, deep+1); subSplitContentNode != nil {
				splitContentNode.SubSplitContentList = append(splitContentNode.SubSplitContentList, subSplitContentNode)
			}
		}
	}

	return splitContentNode
}

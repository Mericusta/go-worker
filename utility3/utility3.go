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

// // CalculatePunctuationMarksContentLength 计算成对标点符号的内容的长度
// func CalculatePunctuationMarksContentLength(afterLeftContent string, leftPunctuationMark, rightPunctuationMark rune) int {
// 	leftCount := 1
// 	rightCount := 0
// 	length := -1
// 	strings.IndexFunc(afterLeftContent, func(r rune) bool {
// 		if r == leftPunctuationMark {
// 			leftCount++
// 		} else if r == rightPunctuationMark {
// 			rightCount++
// 		}
// 		length++
// 		return leftCount == rightCount
// 	})
// 	return length
// }

// CalculatePunctuationMarksContentLength 计算成对符号的内容长度
// @contentAfterLeftPunctuationMark 待计算的字符串，不包括起始符号
// @leftPunctuationMark 符号左边界字符
// @rightPunctuationMark 符号右边界字符
// @invalidScopePunctuationMarkMap 排除计算的边界符号
// @return
func CalculatePunctuationMarksContentLength(contentAfterLeftPunctuationMark string, leftPunctuationMark, rightPunctuationMark rune, invalidScopePunctuationMarkMap map[rune]rune) int {
	length := 0
	leftCount := 1
	rightCount := 0
	isValid := true
	var invalidScopePunctuationMark rune = -1
	strings.IndexFunc(contentAfterLeftPunctuationMark, func(r rune) bool {
		length++

		// end invalid scope
		if !isValid && r == invalidScopePunctuationMark {
			isValid = true
			invalidScopePunctuationMark = -1
			return false
		}

		// in invalid scope
		if !isValid {
			return false
		}

		// begin invalid scope
		if punctuationMark, isInvalidScopePunctuationMark := invalidScopePunctuationMarkMap[r]; isValid && isInvalidScopePunctuationMark {
			isValid = false
			invalidScopePunctuationMark = punctuationMark
			return false
		}

		// out invalid scope
		if r == leftPunctuationMark {
			leftCount++
		} else if r == rightPunctuationMark {
			rightCount++
		}

		if leftCount == rightCount {
			length-- // cut right punctuation mark len
		}
		return leftCount == rightCount
	})
	return length
}

// FixBracketMatchingResult 修正贪婪括号匹配的错误 type(x)( -> type(x)
func FixBracketMatchingResult(content string) string {
	leftBracketIndex := strings.Index(content, string(global.PunctuationMarkLeftBracket))
	if leftBracketIndex == -1 {
		return content
	}
	leftPunctuationMark, rightPunctuationMark := utility2.GetPunctuationMark(global.PunctuationMarkBracket)
	bracketContentLength := CalculatePunctuationMarksContentLength(content[leftBracketIndex+1:], leftPunctuationMark, rightPunctuationMark, global.GoAnalyzerInvalidScopePunctuationMarkMap)
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
	return RecursiveTraitMultiPunctuationMarksContent(content, &utility.PunctuationMarkInfo{
		PunctuationMark: 0,
		Index:           -1,
	}, &utility.PunctuationMarkInfo{
		PunctuationMark: 0,
		Index:           len(content),
	}, leftPunctuationMarkList, maxDeep, 0)
}

// RecursiveTraitMultiPunctuationMarksContent 混合成对标点符号的内容分类提取
// @content 待处理内容
// @leftPunctuationMarkInfo 根节点的左标点符号
// @rightPunctuationMarkInfo 根节点的右标点符号
// @scopeLeftPunctuationMarkList 所有作为划分区域的左标点符号
// @maxDeep 待处理的最大深度
// @deep 当前深度
// @return 根节点
func RecursiveTraitMultiPunctuationMarksContent(content string, leftPunctuationMarkInfo, rightPunctuationMarkInfo *utility.PunctuationMarkInfo, scopeLeftPunctuationMarkList []rune, maxDeep, deep int) *utility.NewPunctuationContent {
	punctuationContent := &utility.NewPunctuationContent{
		Content:                   content,
		LeftPunctuationMark:       leftPunctuationMarkInfo,
		RightPunctuationMark:      rightPunctuationMarkInfo,
		SubPunctuationContentList: make([]*utility.NewPunctuationContent, 0),
	}

	passLeftLength := 0
	for len(content) != 0 && deep != maxDeep {
		var leftPunctuationMark rune
		var rightPunctuationMark rune
		leftPunctuationMarkIndex := len(content) - 1

		for _, toSearchLeftPunctuationMark := range scopeLeftPunctuationMarkList {
			toSearchLeftPunctuationMarkIndex := strings.IndexRune(content, toSearchLeftPunctuationMark)
			if toSearchLeftPunctuationMarkIndex != -1 && toSearchLeftPunctuationMarkIndex < leftPunctuationMarkIndex {
				leftPunctuationMarkIndex = toSearchLeftPunctuationMarkIndex
				leftPunctuationMark = toSearchLeftPunctuationMark
			}
		}
		// fmt.Printf("relative leftPunctuationMarkIndex = %v, leftPunctuationMark = %v\n", leftPunctuationMarkIndex, string(rune(leftPunctuationMark)))

		rightPunctuationMark = utility.GetAnotherPunctuationMark(leftPunctuationMark)
		if leftPunctuationMark == 0 || rightPunctuationMark == 0 || leftPunctuationMarkIndex == len(content)-1 {
			break
		}

		afterLeftPunctuationMarkContentIndex := leftPunctuationMarkIndex + 1

		// fmt.Printf("pass CalculatePunctuationMarksContentLength = |%v|\n", content[afterLeftPunctuationMarkContentIndex:])
		length := CalculatePunctuationMarksContentLength(content[afterLeftPunctuationMarkContentIndex:], leftPunctuationMark, rightPunctuationMark, global.GoAnalyzerInvalidScopePunctuationMarkMap)

		// fmt.Printf("after CalculatePunctuationMarksContentLength, length = %v\n", length)

		rightPunctuationMarkIndex := leftPunctuationMarkIndex + length + 1
		if rightPunctuationMarkIndex >= len(content) {
			// fmt.Printf("rightPunctuationMarkIndex %v >= len(content) %v\n", rightPunctuationMarkIndex, len(content))
			break
		}

		// fmt.Printf("relative rightPunctuationMarkIndex = %v, rightPunctuationMark = %v\n", rightPunctuationMarkIndex, string(rune(rightPunctuationMark)))
		// fmt.Printf("pass content = |%v|\n", content[leftPunctuationMarkIndex+1:rightPunctuationMarkIndex])

		subPunctuationContent := RecursiveTraitMultiPunctuationMarksContent(content[leftPunctuationMarkIndex+1:rightPunctuationMarkIndex], &utility.PunctuationMarkInfo{
			PunctuationMark: leftPunctuationMark,
			Index:           leftPunctuationMarkInfo.Index + 1 + passLeftLength + leftPunctuationMarkIndex,
		}, &utility.PunctuationMarkInfo{
			PunctuationMark: rightPunctuationMark,
			Index:           leftPunctuationMarkInfo.Index + 1 + passLeftLength + rightPunctuationMarkIndex,
		}, scopeLeftPunctuationMarkList, maxDeep, deep+1)
		if subPunctuationContent != nil {
			punctuationContent.SubPunctuationContentList = append(punctuationContent.SubPunctuationContentList, subPunctuationContent)
		}

		// fmt.Printf("update content to |%v|\n", content[rightPunctuationMarkIndex+1:])
		content = content[rightPunctuationMarkIndex+1:]
		// fmt.Printf("update pass left from %v to %v\n", passLeftLength, passLeftLength+rightPunctuationMarkIndex+1)
		passLeftLength += rightPunctuationMarkIndex + 1
		// fmt.Println("--------------------------------")
	}

	return punctuationContent
}

// SplitContent 划分内容节点
type SplitContent struct {
	ContentList         []string
	SubSplitContentList []*SplitContent
}

// RecursiveSplitUnderSameDeepPunctuationMarksContent 相同深度的成对标点符号下的内容划分
// @content 待分析的字符串
// @punctuationMarkList 指定成对标点符号
// @splitter 指定分隔符
// @return
func RecursiveSplitUnderSameDeepPunctuationMarksContent(content string, punctuationMarkList []int, splitter string) *SplitContent {
	if punctuationContentNode := TraitMultiPunctuationMarksContent(content, punctuationMarkList, 1); punctuationContentNode != nil {
		return splitUnderSameDeepPunctuationMarksContent(punctuationContentNode, splitter, 0, 0)
	}
	return nil
}

// RecursiveSplitUnderSameDeepPunctuationMarksContentNode 相同深度的成对标点符号下的内容划分
// @punctuationContentNode 成对标点符号的内容根节点，注意：必须是根节点，不能是某个子节点，节点深度必须为 2
// @splitter 指定分隔符
// @return
func RecursiveSplitUnderSameDeepPunctuationMarksContentNode(punctuationContentNode *utility.NewPunctuationContent, splitter string) *SplitContent {
	return splitUnderSameDeepPunctuationMarksContent(punctuationContentNode, splitter, 0, 0)
}

// splitUnderSameDeepPunctuationMarksContent 相同深度的成对标点符号下的内容划分的递归算法
// @punctuationContentNode 成对标点符号的内容根节点，注意：必须是根节点，不能是某个子节点，节点深度 >= 2，分析结果中深度大于 2 的数据不正确
// @splitter 指定分隔符
// @maxDeep 递归最大深度
// @deep 当前深度
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
		for _, subNode := range punctuationContentNode.SubPunctuationContentList {
			// Note: 这里用于判断的依据是子节点相对父节点的左 区间符号 的下标
			// Note: 但是节点的 区间符号 数据中记录的下标是相对于根节点的下标 -> 必须是根节点
			// Note: 所以当节点数只有2时，这个下标可以代表相对父节点（根节点）的下标 -> 节点深度 >= 2
			if subNode.LeftPunctuationMark.Index <= rightIndex && rightIndex <= subNode.RightPunctuationMark.Index {
				inner = true
				offset = subNode.RightPunctuationMark.Index - leftIndex + 1
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

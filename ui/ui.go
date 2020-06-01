package ui

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go-worker/global"
	"github.com/go-worker/utility"
)

// OutputNoteInfo 输出常规信息
func OutputNoteInfo(content ...interface{}) {
	utility.FormatOutput(global.LogMarkNote, "%v", content...)
}

// OutputNoteInfoWithFormat 输出带有格式的常规信息
func OutputNoteInfoWithFormat(format string, content ...interface{}) {
	utility.FormatOutput(global.LogMarkNote, format, content...)
}

// OutputWarnInfo 输出警告信息
func OutputWarnInfo(format string, content ...interface{}) {
	utility.FormatOutput(global.LogMarkWarn, format, content...)
}

// OutputErrorInfo 输出错误信息
func OutputErrorInfo(format string, content ...interface{}) {
	utility.FormatOutput(global.LogMarkError, format, content...)
}

// ParseStyleTemplate 解析样式模板
func ParseStyleTemplate(templateStyleRegexp *regexp.Regexp, content string) string {
	replaceContent := content
	for _, styleTemplateSubmatchList := range templateStyleRegexp.FindAllStringSubmatch(replaceContent, -1) {
		styleContet := ""
		styleChar := styleTemplateSubmatchList[1]
		styleNum := styleTemplateSubmatchList[2]
		if styleChar != "" && styleNum != "" {
			num, parseNumError := strconv.Atoi(styleNum)
			if parseNumError == nil {
				for index := 0; index != num; index++ {
					styleContet = strings.Repeat(styleChar, num)
				}
			} else {
				OutputWarnInfo("%v", parseNumError)
			}
		}
		replaceContent = strings.Replace(replaceContent, styleTemplateSubmatchList[0], styleContet, -1)
	}
	return replaceContent
}

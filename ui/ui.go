package ui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-worker/global"
	"github.com/go-worker/utility2"
)

// OutputNoteInfoWithOutFormat 输出常规信息
func OutputNoteInfoWithOutFormat(content ...interface{}) {
	utility2.FormatOutput(global.LogMarkNote, "%v", content...)
}

// OutputNoteInfo 输出带有格式的常规信息
func OutputNoteInfo(format string, content ...interface{}) {
	utility2.FormatOutput(global.LogMarkNote, format, content...)
}

// OutputWarnInfo 输出警告信息
func OutputWarnInfo(format string, content ...interface{}) {
	utility2.FormatOutput(global.LogMarkWarn, format, content...)
}

// OutputErrorInfo 输出错误信息
func OutputErrorInfo(format string, content ...interface{}) {
	utility2.FormatOutput(global.LogMarkError, format, content...)
}

// ParseStyleTemplate 解析样式模板
func ParseStyleTemplate(templateStyleRegexp *regexp.Regexp, content string) string {
	replaceContent := content
	for _, styleTemplateSubmatchList := range templateStyleRegexp.FindAllStringSubmatch(replaceContent, -1) {
		styleContet := ""
		styleChar := styleTemplateSubmatchList[1]
		styleNum := styleTemplateSubmatchList[2]
		if styleChar != "" && styleNum != "" {
			charASCII, parseCharASCIIError := strconv.Atoi(styleChar)
			num, parseNumError := strconv.Atoi(styleNum)
			if parseCharASCIIError == nil && parseNumError == nil {
				for index := 0; index != num; index++ {
					styleContet = strings.Repeat(string(charASCII), num)
				}
			} else {
				OutputWarnInfo(CommonError6, fmt.Sprintf("charASCII error: %v, num error: %v", parseCharASCIIError, parseNumError))
			}
		}
		replaceContent = strings.Replace(replaceContent, styleTemplateSubmatchList[0], styleContet, -1)
	}
	return replaceContent
}

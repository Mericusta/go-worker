package ui

import (
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

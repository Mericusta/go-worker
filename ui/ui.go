package ui

import "github.com/go-worker/utility"

// OutputErrorInfo 输出错误信息
func OutputErrorInfo(err error) {
	utility.ErrorOutput("%v\n", err)
}

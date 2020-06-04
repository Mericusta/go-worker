package utility3

import (
	"fmt"
	"io/ioutil"

	"github.com/go-worker/global"
	"github.com/go-worker/regexps"
	"github.com/go-worker/ui"
)

// TrimSpaceLine 移除空白行
func TrimSpaceLine(content string) string {
	spaceLineRegexp, hasSpaceLineRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AESpaceLine]
	if !hasSpaceLineRegexp {
		ui.OutputWarnInfo(ui.CommonWarn3, global.AESpaceLine)
		return ""
	}
	replaceContent := spaceLineRegexp.ReplaceAllString(content, "\n")
	return replaceContent
}

// TraverseDirectorySpecificFile 遍历文件夹获取所有绑定类型的文件
func TraverseDirectorySpecificFile(directory, syntax string) []string {
	traverseFileList := make([]string, 0)
	toTraverseDirectoryList := []string{directory}
	fileTypeRegexp, hasFileTypeRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEFileNameType]
	if !hasFileTypeRegexp {
		return traverseFileList
	}
	for len(toTraverseDirectoryList) != 0 {
		toTraverseDirectory := toTraverseDirectoryList[0]
		fileList, readDirError := ioutil.ReadDir(toTraverseDirectory)
		if readDirError == nil {
			for _, file := range fileList {
				filePath := fmt.Sprintf("%v/%v", toTraverseDirectory, file.Name())
				if file.IsDir() {
					toTraverseDirectoryList = append(toTraverseDirectoryList, filePath)
				} else if fileTypeRegexp.ReplaceAllString(filePath, "$TYPE") == syntax {
					traverseFileList = append(traverseFileList, filePath)
				} else {
				}
			}
		} else {
			ui.OutputWarnInfo(ui.CommonError9, toTraverseDirectory, readDirError)
		}
		toTraverseDirectoryList = toTraverseDirectoryList[1:]
	}
	return traverseFileList
}

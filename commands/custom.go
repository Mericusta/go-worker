package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-worker/ui"
	"github.com/go-worker/utility"
	"github.com/go-worker/utility2"
)

type Custom struct {
	*CommandStruct
	Params *customParam
}

var CustomExecutor map[int]func([]string)

func init() {
	CustomExecutor = map[int]func([]string){
		1: RecursivelyCountFileSizeInDirectory,
		2: ConcurrentScanDirectory,
	}
}

func (command *Custom) Execute() error {
	// 解析指令的选项和参数
	parseCommandParamsError := command.parseCommandParams()
	if parseCommandParamsError != nil {
		return parseCommandParamsError
	}
	fmt.Println("params = %+v", command.Params)

	executor, hasExecutor := CustomExecutor[command.Params.optionValue]
	if !hasExecutor || executor == nil {
		return fmt.Errorf(ui.CMDCustomExecutorNotExist, command.Params.optionValue)
	}
	executor(command.Params.paramList)

	return nil
}

type customParam struct {
	option      string
	optionValue int
	paramList   []string
}

func (command *Custom) parseCommandParams() error {
	optionValueList := strings.Split(command.CommandStruct.InputString, " ")
	optionValue, parseOptionValueError := strconv.Atoi(optionValueList[2])
	if parseOptionValueError != nil {
		return parseOptionValueError
	}
	command.Params = &customParam{
		option:      optionValueList[1],
		optionValue: optionValue,
		paramList:   optionValueList[3:],
	}
	return nil
}

// ----------------------------------------------------------------

// RecursivelyCountFileSizeInDirectory 递归统计目录下指定文件类型的大小
func RecursivelyCountFileSizeInDirectory(paramList []string) {
	if len(paramList) < 2 {
		ui.OutputErrorInfo(ui.CMDCustomExecutorHasNotEnoughParam, 1)
		return
	}
	directory := paramList[0]
	fileType := paramList[1]
	fileSizeListMap := make(map[int64][]string)
	utility2.TestOutput("directory = %v", directory)
	directoryStat, getStatError := os.Stat(directory)
	if getStatError != nil {
		ui.OutputErrorInfo(ui.CommonError7, directory, getStatError)
		return
	}
	if !directoryStat.IsDir() {
		ui.OutputWarnInfo(ui.CommonError8, directory)
		return
	}
	utility.TraverseDirectorySpecificFileWithFunction(directory, fileType, func(filePath string, info os.FileInfo) {
		if _, hasSizeList := fileSizeListMap[info.Size()]; !hasSizeList {
			fileSizeListMap[info.Size()] = make([]string, 0)
		}
		fileSizeListMap[info.Size()] = append(fileSizeListMap[info.Size()], filePath)
	})
	for fileSize, fileList := range fileSizeListMap {
		utility2.TestOutput("fileSize = %v", fileSize)
		utility2.TestOutput("fileList = %v", fileList)
	}
	return
}

// ----------------------------------------------------------------

// ConcurrentScanDirectory 并发扫描目录
func ConcurrentScanDirectory(paramList []string) {
	if len(paramList) < 2 {
		ui.OutputErrorInfo(ui.CMDCustomExecutorHasNotEnoughParam, 2)
		return
	}
	directory := paramList[0]
	fileType := paramList[1]
	ignoreDirectoryList := paramList[2:]
	scanChannel := make(chan []string)
	go func() {
		scanChannel <- []string{directory}
	}()

	scanMap := make(map[int64][]string)
	goRoutineNum := 1
	maxGoRoutineNum := goRoutineNum
	for subFileList := range scanChannel {
		for _, subFile := range subFileList {
			fileInfo, getStatError := os.Stat(subFile)
			if getStatError != nil {
				ui.OutputErrorInfo(ui.CommonError7, directory, getStatError)
				continue
			}
			subFilePath, getAbsError := filepath.Abs(subFile)
			if getAbsError != nil {
				ui.OutputErrorInfo(ui.CommonError11, directory, getAbsError)
				continue
			}
			if fileInfo.IsDir() {
				goRoutineNum++
				maxGoRoutineNum++
				go func(directory string, ch chan []string) {
					directoryFileList, readDirError := ioutil.ReadDir(subFilePath)
					if readDirError != nil {
						ui.OutputErrorInfo(ui.CommonError12, subFilePath, readDirError)
						scanChannel <- []string{}
					} else {
						targetList := make([]string, 0)
						for _, directoryFile := range directoryFileList {
							if (directoryFile.IsDir() && func() bool {
								for _, ignoreDirectory := range ignoreDirectoryList {
									if directoryFile.Name() == ignoreDirectory {
										return false
									}
								}
								return true
							}()) || filepath.Ext(directoryFile.Name()) == fileType {
								targetList = append(targetList, filepath.Join(subFilePath, directoryFile.Name()))
							}
						}
						scanChannel <- targetList
					}
				}(subFilePath, scanChannel)
			} else {
				if _, hasSize := scanMap[fileInfo.Size()]; !hasSize {
					scanMap[fileInfo.Size()] = make([]string, 0)
				}
				scanMap[fileInfo.Size()] = append(scanMap[fileInfo.Size()], subFilePath)
			}
		}
		goRoutineNum--
		if goRoutineNum == 0 {
			ui.OutputNoteInfo("end with %v go routine", maxGoRoutineNum)
			break
		}
	}
	for fileSize, filePathList := range scanMap {
		ui.OutputNoteInfo("fileSize = %v", fileSize)
		for _, filePath := range filePathList {
			ui.OutputNoteInfo("filePath = %v", filePath)
		}
	}
}

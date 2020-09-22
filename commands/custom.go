package commands

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
		3: MorrisTraverseBinaryTree,
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
	executeBeginTime := time.Now()
	executor(command.Params.paramList)
	executeEndTime := time.Now()
	ui.OutputNoteInfo("execute custom command %v done, using: %v ns", command.Params.optionValue, executeEndTime.Sub(executeBeginTime).Nanoseconds())

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

// Command Example: custom execute 1 . md

// RecursivelyCountFileSizeInDirectory 递归统计目录下指定文件类型的大小
func RecursivelyCountFileSizeInDirectory(paramList []string) {
	if len(paramList) < 2 {
		ui.OutputErrorInfo(ui.CMDCustomExecutorHasNotEnoughParam, 1)
		return
	}
	directory := paramList[0]
	fileType := paramList[1]
	fileSizeListMap := make(map[int64][]string)
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
	for fileSize, filePathList := range fileSizeListMap {
		ui.OutputNoteInfo("fileSize = %v", fileSize)
		for _, filePath := range filePathList {
			ui.OutputNoteInfo("filePath = %v", filePath)
		}
	}
	return
}

// ----------------------------------------------------------------

// Command Example: custom execute 2 . md .git

// ConcurrentScanDirectory 并发扫描目录下指定文件类型的大小
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

// ----------------------------------------------------------------

// Command Example: custom execute 3

// BTNode 二叉树结点
type BTNode struct {
	Value int
	Left  *BTNode
	Right *BTNode
}

// MorrisTraverseBinaryTree Morris 遍历二叉树
func MorrisTraverseBinaryTree(paramList []string) {
	rootNode := randomGenerateBinaryTree()
	currentNode := rootNode
	var mostRightNode *BTNode
	outputTemplate := "Morris Traverse Binary Tree:"
	for currentNode != nil {
		ui.OutputNoteInfo("%v current Node is %v", outputTemplate, currentNode.Value)
		if currentNode.Left == nil {
			ui.OutputNoteInfo("%v current node left child is nil, move current node to right child", outputTemplate)
			currentNode = currentNode.Right
		} else {
			ui.OutputNoteInfo("%v current node left child is not nil", outputTemplate)
			mostRightNode = currentNode.Left.Right
			for mostRightNode != nil && mostRightNode.Right != nil {
				utility2.TestOutput("change most right node from %+v to %+v", mostRightNode, mostRightNode.Right)
				mostRightNode = mostRightNode.Right
			}

			ui.OutputNoteInfo("%v most right node is %v", outputTemplate, mostRightNode.Value)
			if mostRightNode.Right == nil {
				ui.OutputNoteInfo("%v most right node right child is nil, point to current node, move current node to left child", outputTemplate)
				mostRightNode.Right = currentNode
				currentNode = currentNode.Left
			} else if mostRightNode.Right == currentNode {
				ui.OutputNoteInfo("%v most right node right child is not nil, move current node to right child", outputTemplate)
				mostRightNode.Right = nil
				currentNode = currentNode.Right
			}
		}
	}
}

// randomGenerateBinaryTree 随机生成二叉树
func randomGenerateBinaryTree() *BTNode {
	rand.Seed(time.Now().UnixNano())
	// nodeCount := rand.Intn(31) + 1
	nodeCount := 7
	nodeList := make([]*BTNode, 0)
	for index := 0; index != nodeCount; index++ {
		nodeList = append(nodeList, &BTNode{
			Value: index + 1,
		})
	}
	rootNode := nodeList[0]
	nodeList = nodeList[1:]
	toAppendChildrenNodeList := []*BTNode{rootNode}
	for len(toAppendChildrenNodeList) != 0 {
		toAppendChildrenNode := toAppendChildrenNodeList[0]
		toAppendChildrenNodeList = toAppendChildrenNodeList[1:]
		if len(nodeList) > 0 {
			toAppendChildrenNode.Left = nodeList[0]
			toAppendChildrenNodeList = append(toAppendChildrenNodeList, nodeList[0])
			nodeList = nodeList[1:]
		} else {
			break
		}
		if len(nodeList) > 0 {
			toAppendChildrenNode.Right = nodeList[0]
			toAppendChildrenNodeList = append(toAppendChildrenNodeList, nodeList[0])
			nodeList = nodeList[1:]
		} else {
			break
		}
	}
	return rootNode
}

// ----------------------------------------------------------------

// Command Example: custom execute 4

// test objects:
// charge_config.shop_data
// shop_config.shop_id
// shoplist.shop_id

// example rules:
// table.field -> table.field [table.field]
// table.field ->

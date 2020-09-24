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
		4: RandomGenerateOnmyojiEquipments,
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

// Command Example: custom execute 4 100 1

type ce4AttributeType int

const (
	// MaxLevel max times of update attribute
	MaxLevel int = 5
	// MaxSubAttributeNum max num of sub attribute
	MaxSubAttributeNum int = 4
	// AttributeR Attronite Rate
	AttributeR ce4AttributeType = 10
	// ATKR ATK Rencentage
	ATKR ce4AttributeType = 1
	// DEFR DEF Rencentage
	DEFR ce4AttributeType = 2
	// LIFER LIFE Rencentage
	LIFER ce4AttributeType = 3
	// SV Speed Value
	SV ce4AttributeType = 4
	// ATKV ATK Value
	ATKV ce4AttributeType = 11
	// DEFV DEF Value
	DEFV ce4AttributeType = 12
	// LIFEV LIFE Value
	LIFEV ce4AttributeType = 13
	// CR Critical Rate
	CR ce4AttributeType = 21
	// CD Critical Demage
	CD ce4AttributeType = 22
	// AttributeEffect Attribute Effect
	AttributeEffect ce4AttributeType = 50
	// EH Effect Hit
	EH ce4AttributeType = 41
	// ER Effect Resistance
	ER ce4AttributeType = 42
)

var attributeTypeMap map[int]ce4AttributeType

func initAttributeType() {
	attributeTypeList := []ce4AttributeType{
		ATKR, DEFR, LIFER, SV,
		ATKV, DEFV, LIFEV,
		CR, CD,
		EH, ER,
	}

	attributeTypeMap = map[int]ce4AttributeType{}
	for index, attributeType := range attributeTypeList {
		attributeTypeMap[index] = attributeType
	}
}

type ce4Attribute struct {
	Type  ce4AttributeType
	Value int
}

type ce4Equipment struct {
	MainAttribute *ce4Attribute
	Attribute1    *ce4Attribute
	Attribute2    *ce4Attribute
	Attribute3    *ce4Attribute
	Attribute4    *ce4Attribute
}

// position type isMain min max
// 1 ATKV 1 486 486
// 2 ATKR

// main
// - ATKV 486 486
// - DEFV 104 104
// - LIFEV 2052 2052
// - xxxR 55 55

// sub
// - xxxR = initR + eachR * level
// - xxxV = initV + eachV * level

// - xxxInitR const
// - xxxEachR const
// - level random

// equipment level const
// - attribute0 level + ... + attributeN level = equipment level

// type            init   each
// AttributeR
// ATKR            2,3    2,3
// DEFR            2,3    2,3
// LIFER           2,3    2,3
// SV              2,3    2,3
// --------------------------------
// ATKV            22,27  22,24
// DEFV            4,5    4,5
// LIFEV           91,114 95,105
// CR              2,3    3
// CD              3,4    4
// --------------------------------
// AttributeEffect
// EH              3,4    3,4
// ER              3,4    3,4

var attributeConfigMap map[ce4AttributeType]map[string]map[string]int

func initAttributeConfigMap() {
	attributeConfigMap = map[ce4AttributeType]map[string]map[string]int{
		AttributeR: map[string]map[string]int{
			"init": map[string]int{
				"min": 2,
				"max": 3,
			},
			"each": map[string]int{
				"min": 2,
				"max": 3,
			},
		},
		AttributeEffect: map[string]map[string]int{
			"init": map[string]int{
				"min": 3,
				"max": 4,
			},
			"each": map[string]int{
				"min": 3,
				"max": 4,
			},
		},
		ATKV: map[string]map[string]int{
			"init": map[string]int{
				"min": 22,
				"max": 27,
			},
			"each": map[string]int{
				"min": 22,
				"max": 24,
			},
		},
		DEFV: map[string]map[string]int{
			"init": map[string]int{
				"min": 4,
				"max": 5,
			},
			"each": map[string]int{
				"min": 4,
				"max": 5,
			},
		},
		LIFEV: map[string]map[string]int{
			"init": map[string]int{
				"min": 91,
				"max": 114,
			},
			"each": map[string]int{
				"min": 95,
				"max": 105,
			},
		},
		CR: map[string]map[string]int{
			"init": map[string]int{
				"min": 2,
				"max": 3,
			},
			"each": map[string]int{
				"min": 3,
				"max": 3,
			},
		},
		CD: map[string]map[string]int{
			"init": map[string]int{
				"min": 3,
				"max": 4,
			},
			"each": map[string]int{
				"min": 4,
				"max": 4,
			},
		},
	}
}

// RandomGenerateOnmyojiEquipments 模拟 Onmyoji 御魂生成机制随机生成御魂
func RandomGenerateOnmyojiEquipments(paramList []string) {
	if len(paramList) < 1 {
		ui.OutputErrorInfo(ui.CMDCustomExecutorHasNotEnoughParam, 4)
		return
	}
	var generateNumber int
	var atoiGenerateNumberError error
	var generatePosition int
	var atoiGeneratePositionError error
	generateNumber, atoiGenerateNumberError = strconv.Atoi(paramList[0])
	if len(paramList) > 1 {
		generatePosition, atoiGeneratePositionError = strconv.Atoi(paramList[1])
	}
	if atoiGenerateNumberError != nil || atoiGeneratePositionError != nil {
		ui.OutputErrorInfo(ui.CMDCustomExecutorParseParamError, 4)
		return
	}

	rand.Seed(time.Now().Unix())

	var logConstant1 string = "No.%v equipment"

	for index := 0; index != generateNumber; index++ {
		ui.OutputNoteInfo(fmt.Sprintf("generate %v", logConstant1), index)
		position := generatePosition
		if generatePosition == 0 {
			position = rand.Intn(6) + 1
		}
		ui.OutputNoteInfo(fmt.Sprintf("%v position: %v", logConstant1, position), index)

		initAttributeNum := rand.Intn(3) + 2
		ui.OutputNoteInfo(fmt.Sprintf("%v init attribute num: %v", logConstant1, initAttributeNum), index)

		leftUpdateTimes := MaxLevel - (4 - initAttributeNum)
		ui.OutputNoteInfo(fmt.Sprintf("%v left update times: %v", logConstant1, leftUpdateTimes), index)

		ui.OutputNoteInfo(ui.CommonNote2)
	}
}

func randomGenerateSubAttributes() map[ce4AttributeType]*ce4Attribute {
	ce4AttributeMap := make(map[ce4AttributeType]*ce4Attribute)

	for index := 0; index != MaxSubAttributeNum; index++ {
		for {
			randomValue := rand.Intn(len(attributeTypeMap))
			if attributeType, hasAttributeType := attributeTypeMap[randomValue]; hasAttributeType {
				if _, hasAttribute := ce4AttributeMap[attributeType]; hasAttribute {
					continue
				}

			} else {
				ui.OutputWarnInfo(ui.CMDCustomExecutorOutOfRangeError, 4)
				continue
			}
		}
	}

	return ce4AttributeMap
}

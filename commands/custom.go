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
		5: ProofOfArrayOrdered,
		6: GoCommandToolTemplater,
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
// Command Expression:
// - custom : command const content
// - execute: command const content
// - 1      : specify executor
// - .      : specify directory to count
// - md     : specify file type to count

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
// Command Expression:
// - custom : command const content
// - execute: command const content
// - 2      : specify executor
// - .      : specify directory to scan
// - md     : specify file type to scan
// - .git   : specify directory to ignore

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
// Command Expression:
// - custom : command const content
// - execute: command const content
// - 3      : specify executor

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
// Command Expression:
// - custom : command const content
// - execute: command const content
// - 4      : specify executor
// - 100    : specify the number of items to generate
// - 1      : specify item's position

type ce4AttributeType int

func (attributeType ce4AttributeType) String() string {
	switch attributeType {
	case ATKR:
		return "ATKR"
	case DEFR:
		return "DEFR"
	case LIFER:
		return "LIFER"
	case SV:
		return "SV"
	case ATKV:
		return "ATKV"
	case DEFV:
		return "DEFV"
	case LIFEV:
		return "LIFEV"
	case CR:
		return "CR"
	case CD:
		return "CD"
	case EH:
		return "EH"
	case ER:
		return "ER"
	case AttributeRGroup:
		return "AttributeRGroup"
	case AttributeEGroup:
		return "AttributeEGroup"
	}
	return "unknown"
}

// 0  0  0  0  0     0    0    0  0     0    0
// ER EH CD CR LIFEV DEFV ATKV SV LIFER DEFR ATKR
// AttributeRGroup = 00000001111
// AttributeEGroup = 11000000000
// Attribute & AttributeRGroup != 0 -> Attribute has Attribute R
// Attribute & AttributeEGroup != 0 -> Attribute has Attribute Effect

type ce4EquipmentSuitType int

func (equipmentSuitType ce4EquipmentSuitType) String() string {
	switch equipmentSuitType {
	case ATK4Group:
		return "ATK4Group"
	case DEF4Group:
		return "DEF4Group"
	case LIFE4Group:
		return "LIFE4Group"
	case CR4Group:
		return "CR4Group"
	case EH4Group:
		return "EH4Group"
	case ER4Group:
		return "ER4Group"
	}
	return "unknown"
}

const (
	// ATKR ATK Rencentage
	ATKR ce4AttributeType = 1 << iota
	// DEFR DEF Rencentage
	DEFR
	// LIFER LIFE Rencentage
	LIFER
	// SV Speed Value
	SV
	// ATKV ATK Value
	ATKV
	// DEFV DEF Value
	DEFV
	// LIFEV LIFE Value
	LIFEV
	// CR Critical Rate
	CR
	// CD Critical Demage
	CD
	// EH Effect Hit
	EH
	// ER Effect Resistance
	ER

	// AttributeRGroup Attronite Rate Group
	AttributeRGroup ce4AttributeType = ATKR | DEFR | LIFER | SV
	// AttributeEGroup Attribute Effect Group
	AttributeEGroup ce4AttributeType = EH | ER

	// MaxLevel max times of update attribute
	MaxLevel int = 5
	// MaxSubAttributeNum max num of sub attribute
	MaxSubAttributeNum int = 4

	// ATK4Group ATK 4 item suit type Group
	ATK4Group ce4EquipmentSuitType = 10
	// DEF4Group DEF 4 item suit type Group
	DEF4Group ce4EquipmentSuitType = 20
	// LIFE4Group LIFE 4 item suit type Group
	LIFE4Group ce4EquipmentSuitType = 30
	// CR4Group CR 4 item suit type Group
	CR4Group ce4EquipmentSuitType = 40
	// EH4Group EH 4 item suit type Group
	EH4Group ce4EquipmentSuitType = 50
	// ER4Group ER 4 item suit type Group
	ER4Group ce4EquipmentSuitType = 60
)

var attributeTypeMap map[int]ce4AttributeType

func initAttributeTypeMap() {
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
	SuitType        ce4EquipmentSuitType
	MainAttribute   *ce4Attribute
	SubAttributeMap map[ce4AttributeType]*ce4Attribute
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
// AttributeRGroup
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
// AttributeEGroup
// EH              3,4    3,4
// ER              3,4    3,4

var attributeConfigMap map[ce4AttributeType]map[string]map[string]int
var attributeMainValueMap map[ce4AttributeType]int

func initAttributeConfigMap() {
	attributeConfigMap = map[ce4AttributeType]map[string]map[string]int{
		AttributeRGroup: {
			"init": {
				"min": 2,
				"max": 3,
			},
			"each": {
				"min": 2,
				"max": 3,
			},
		},
		AttributeEGroup: {
			"init": {
				"min": 3,
				"max": 4,
			},
			"each": {
				"min": 3,
				"max": 4,
			},
		},
		ATKV: {
			"init": {
				"min": 22,
				"max": 27,
			},
			"each": {
				"min": 22,
				"max": 24,
			},
		},
		DEFV: {
			"init": {
				"min": 4,
				"max": 5,
			},
			"each": {
				"min": 4,
				"max": 5,
			},
		},
		LIFEV: {
			"init": {
				"min": 91,
				"max": 114,
			},
			"each": {
				"min": 95,
				"max": 105,
			},
		},
		CR: {
			"init": {
				"min": 2,
				"max": 3,
			},
			"each": {
				"min": 3,
				"max": 3,
			},
		},
		CD: {
			"init": {
				"min": 3,
				"max": 4,
			},
			"each": {
				"min": 4,
				"max": 4,
			},
		},
	}
	attributeMainValueMap = map[ce4AttributeType]int{
		ATKV:            486,
		DEFV:            104,
		LIFEV:           2052,
		AttributeRGroup: 55,
		AttributeEGroup: 55,
	}
}

func checkAttributeConfig() (bool, int, int) {
	for _, attributeConfig := range attributeConfigMap {
		if attributeConfig["init"]["max"] < attributeConfig["init"]["min"] {
			return false, attributeConfig["init"]["max"], attributeConfig["init"]["min"]
		} else if attributeConfig["each"]["max"] < attributeConfig["each"]["min"] {
			return false, attributeConfig["each"]["max"], attributeConfig["each"]["min"]
		}
	}
	return true, 0, 0
}

func getAttributeConfig(attributeType ce4AttributeType) map[string]map[string]int {
	return attributeConfigMap[getAttributeGroup(attributeType)]
}

func getAttributeMainValue(attributeType ce4AttributeType) int {
	return attributeMainValueMap[getAttributeGroup(attributeType)]
}

func getAttributeGroup(attributeType ce4AttributeType) ce4AttributeType {
	if attributeType&AttributeRGroup != 0 {
		return AttributeRGroup
	} else if attributeType&AttributeEGroup != 0 {
		return AttributeEGroup
	} else {
		return attributeType
	}
}

var logConstant1 string = "No.%v equipment"

func logWarp(index int, preFormat, postFormat string, value ...interface{}) {
	format := logConstant1
	if len(preFormat) != 0 {
		format = fmt.Sprintf("%v %v", preFormat, format)
	}
	if len(postFormat) != 0 {
		format = fmt.Sprintf("%v %v", format, postFormat)
	}

	outputNoteInfoValue := make([]interface{}, 0)
	outputNoteInfoValue = append(outputNoteInfoValue, index)
	outputNoteInfoValue = append(outputNoteInfoValue, value...)

	ui.OutputNoteInfo(format, outputNoteInfoValue...)
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

	initAttributeTypeMap()
	initAttributeConfigMap()
	if noError, errorMaxValue, errorMinValue := checkAttributeConfig(); !noError {
		ui.OutputErrorInfo(ui.CMDCustomExecutorConfigDataError, 4, fmt.Sprintf("config max value %v should greater than or equal to min value %v", errorMaxValue, errorMinValue))
		return
	}

	rand.Seed(time.Now().Unix())

	for index := 0; index != generateNumber; index++ {
		logWarp(index, "generate", "")
		position := generatePosition
		if generatePosition == 0 {
			position = rand.Intn(6) + 1
		}
		logWarp(index, "", "position: %v", position)

		initAttributeNum := rand.Intn(3) + 2
		logWarp(index, "", "init attribute num: %v", initAttributeNum)

		leftUpdateTimes := MaxLevel - (4 - initAttributeNum)
		logWarp(index, "", "left update times: %v", leftUpdateTimes)

		attributeMap, attributeTypeList := randomGenerateSubAttributes()
		logWarp(index, "", "attribute list: %v", attributeTypeList)
		for _, attribute := range attributeMap {
			logWarp(index, "", "init attribute %s, value %v", attribute.Type, attribute.Value)
		}

		randomUpdateSubAttributes(leftUpdateTimes, attributeMap, attributeTypeList)
		logWarp(index, "", "after update:")
		for _, attribute := range attributeMap {
			logWarp(index, "", "attribute %s, value %v", attribute.Type, attribute.Value)
		}

		ce4Equipment := makeUpEquipment(position, attributeMap)
		logWarp(index, "", "after generate: \nsuit type: %s\nmain attribute: %s:%v\nsub attributes:\n%v", ce4Equipment.SuitType, ce4Equipment.MainAttribute.Type, ce4Equipment.MainAttribute.Value, func() string {
			content := "%v\n%v\n%v\n%v"
			list := make([]interface{}, 0, 4)
			for _, attribute := range ce4Equipment.SubAttributeMap {
				list = append(list, fmt.Sprintf("- %s:%v", attribute.Type, attribute.Value))
			}
			content = fmt.Sprintf(content, list...)
			return content
		}())

		ui.OutputNoteInfo(ui.CommonNote2)
	}
}

func makeUpEquipment(position int, attributeMap map[ce4AttributeType]*ce4Attribute) *ce4Equipment {
	mainAttribute := &ce4Attribute{}
	randomValue := rand.Intn(len(attributeTypeMap))
	attributeType, hasAttributeType := attributeTypeMap[randomValue]
	if !hasAttributeType {
		ui.OutputWarnInfo(ui.CMDCustomExecutorOutOfRangeError, 4, fmt.Sprintf("attributeTypeMap does not have randomValue %v", randomValue))
		return nil
	}
	mainAttribute.Type = attributeType
	mainAttribute.Value = getAttributeMainValue(attributeType)

	equipment := &ce4Equipment{
		SuitType:        ce4EquipmentSuitType((rand.Intn(6) + 1) * 10),
		MainAttribute:   mainAttribute,
		SubAttributeMap: attributeMap,
	}
	return equipment
}

func randomGenerateSubAttributes() (map[ce4AttributeType]*ce4Attribute, []ce4AttributeType) {
	ce4AttributeMap := make(map[ce4AttributeType]*ce4Attribute, 4)
	ce4AttributeList := make([]ce4AttributeType, 0, MaxSubAttributeNum)

	for index := 0; index != MaxSubAttributeNum; index++ {
		attribute := &ce4Attribute{}
		for {
			randomValue := rand.Intn(len(attributeTypeMap))
			if attributeType, hasAttributeType := attributeTypeMap[randomValue]; hasAttributeType {
				if _, hasAttribute := ce4AttributeMap[attributeType]; hasAttribute {
					continue
				}
				attribute.Type = attributeType
				break
			} else {
				ui.OutputWarnInfo(ui.CMDCustomExecutorOutOfRangeError, 4, fmt.Sprintf("attributeTypeMap does not have randomValue %v", randomValue))
				continue
			}
		}
		attributeConfig := getAttributeConfig(attribute.Type)
		if attributeConfig == nil {
			ui.OutputWarnInfo(ui.CMDCustomExecutorOutOfRangeError, 4, fmt.Sprintf("attributeConfigMap does not have attribute type %v", attribute.Type))
			continue
		}
		attribute.Value = attributeConfig["init"]["min"] + rand.Intn(attributeConfig["init"]["max"]-attributeConfig["init"]["min"])
		ce4AttributeMap[attribute.Type] = attribute
		ce4AttributeList = append(ce4AttributeList, attribute.Type)
	}

	return ce4AttributeMap, ce4AttributeList
}

func randomUpdateSubAttributes(leftUpdateTimes int, ce4AttributeMap map[ce4AttributeType]*ce4Attribute, attributeTypeList []ce4AttributeType) {
	for index := 0; index != leftUpdateTimes; index++ {
		attributeTypeToUpdate := attributeTypeList[rand.Intn(len(attributeTypeList))]
		attributeConfig := getAttributeConfig(attributeTypeToUpdate)
		if attributeConfig == nil {
			ui.OutputWarnInfo(ui.CMDCustomExecutorOutOfRangeError, 4, fmt.Sprintf("attributeConfigMap does not have attribute type %v", attributeTypeToUpdate))
			continue
		}
		attribute, hasAttribute := ce4AttributeMap[attributeTypeToUpdate]
		if !hasAttribute {
			ui.OutputErrorInfo(ui.CMDCustomExecutorOutOfRangeError, 4, fmt.Sprintf("ce4AttributeMap does not have attribute type %v", attributeTypeToUpdate))
			continue
		}
		if attributeConfig["each"]["max"] != attributeConfig["each"]["min"] {
			attribute.Value = attribute.Value + attributeConfig["each"]["min"] + rand.Intn(attributeConfig["each"]["max"]-attributeConfig["each"]["min"])
		} else {
			attribute.Value = attribute.Value + attributeConfig["each"]["min"]
		}
	}
}

// x1 + x2 + x3 + x4 + x5 + x6 = t1
// y1 + y2 + y3 + y4 + y5 + y6 = t2
// z1 + z2 + z3 + z4 + z5 + z6 = t3
//
// xn % en = yn % en = zn % en

// t = 90

// cr()

// ----------------------------------------------------------------

<<<<<<< HEAD
// Command Example: custom execute 5 resources/templater_example.go
// Command Expression:
// - custom                        : command const content
// - execute                       : command const content
// - 5                             : specify executor
// - resources/templater_example.go: specify a file to analyze

// GoCommandToolTemplater go 语言命令行工具：模板代码生成器
func GoCommandToolTemplater(paramList []string) {
	if len(paramList) < 1 {
		ui.OutputErrorInfo(ui.CMDCustomExecutorHasNotEnoughParam, 5)
		return
	}
	filename := paramList[0]

	contentByte, readFileError := ioutil.ReadFile(filename)
	if readFileError != nil {
		ui.OutputErrorInfo(ui.CommonError5, filename, readFileError)
		return
	}

	removeGoFileCommentLine(contentByte)
	// contentByteWithoutComment := removeGoFileCommentLine(contentByte)
	// utility2.TestOutput("after clear all comment line:\n|%v|", string(removeGoFileCommentLine(contentByte)))

	// goFileAnalysis := &GoFileAnalysis{FunctionMap: make(map[string]*GoFunctionAnalysis), functionList: make([]string, 0)}
	// analyzeGoFunctionDefinition(goFileAnalysis, contentByte)

	// analyzeGoFunctionBody(goFileAnalysis, contentByte)

	// for functionName, functionAnalysis := range goFileAnalysis.FunctionMap {
	// 	utility2.TestOutput("function name: %v", functionName)
	// 	utility2.TestOutput("function class: %v", functionAnalysis.Class)
	// 	utility2.TestOutput("function params map: %v", functionAnalysis.ParamsMap)
	// 	utility2.TestOutput("function return map: %v", functionAnalysis.ReturnMap)

	// 	// for functionName, callMap := range functionAnalysis.MemberCallMap {
	// 	// 	utility2.TestOutput("function name: %v", functionName)
	// 	// 	for callFunctionName, callFunctionTimes := range callMap {
	// 	// 		utility2.TestOutput("call function: %v, times: %v", callFunctionName, callFunctionTimes)
	// 	// 	}
	// 	// }
	// }

=======
// Command Example: custom execute 5 UP 1,2,3,...
// Command Expression:
// - custom : command const content
// - execute: command const content
// - 5      : specify executor
// - UP     : specify prove type: UP, DOWN, EQUAL
// - 1,2,3  : input list to prove

const (
	// MonotonicityUP 单调递增
	MonotonicityUP = iota << 0
	// MonotonicityDOWN 单调递减
	MonotonicityDOWN
	// MonotonicityEQUAL 相等
	MonotonicityEQUAL
	// GREATER 大于等于
	GREATER = MonotonicityUP | MonotonicityEQUAL
	// LOWER 小于等于
	LOWER = MonotonicityDOWN | MonotonicityEQUAL
)

// ProofOfArrayOrdered 数组单调性证明
func ProofOfArrayOrdered(paramList []string) {
	if len(paramList) < 2 {
		ui.OutputErrorInfo(ui.CMDCustomExecutorHasNotEnoughParam, 1)
		return
	}
	proveTypeString := paramList[0]
	toProveStringSlice := strings.Split(paramList[1], ",")

	proveType := MonotonicityEQUAL
	if proveTypeString == "UP" {
		proveType = MonotonicityUP
	} else if proveTypeString == "DOWN" {
		proveType = MonotonicityDOWN
	} else if proveTypeString == "GREATER" {
		proveType = GREATER
	} else if proveTypeString == "LOWER" {
		proveType = LOWER
	}

	toProveSlice := make([]int, 0, len(toProveStringSlice))
	for _, alpha := range toProveStringSlice {
		integer, atioError := strconv.Atoi(alpha)
		if atioError != nil {
			ui.OutputErrorInfo(ui.CMDCustomExecutorParseParamError, atioError)
			return
		}
		toProveSlice = append(toProveSlice, integer)
	}

	outputNoteFormat := "array %v Monotonicity %v"

	last := toProveSlice[0]
	if proveType == MonotonicityUP {
		for _, value := range toProveSlice[1:] {
			if (value - last) > 0 {
				last = value
				continue
			} else {
				ui.OutputNoteInfo(outputNoteFormat, "is not", proveTypeString)
				return
			}
		}
	} else if proveType == MonotonicityDOWN {
		for _, value := range toProveSlice[1:] {
			if (value - last) < 0 {
				last = value
				continue
			} else {
				ui.OutputNoteInfo(outputNoteFormat, "is not", proveTypeString)
				return
			}
		}
	} else if proveType == GREATER {
		for _, value := range toProveSlice[1:] {
			if (value - last) >= 0 {
				last = value
				continue
			} else {
				ui.OutputNoteInfo(outputNoteFormat, "is not", proveTypeString)
				return
			}
		}
	} else if proveType == LOWER {
		for _, value := range toProveSlice[1:] {
			if (value - last) <= 0 {
				last = value
				continue
			} else {
				ui.OutputNoteInfo(outputNoteFormat, "is not", proveTypeString)
				return
			}
		}
	} else {
		for _, value := range toProveSlice[1:] {
			if (value - last) == 0 {
				last = value
				continue
			} else {
				ui.OutputNoteInfo(outputNoteFormat, "is not", proveTypeString)
				return
			}
		}
	}

	ui.OutputNoteInfo(outputNoteFormat, "is", proveTypeString)

	return
>>>>>>> de950eeb2f060ac77c06f3230eb153efcdbd2184
}

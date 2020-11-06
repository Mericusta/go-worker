package commands

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-worker/global"
	"github.com/go-worker/regexps"
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
		7: GoFileSplitter,
		8: GoGrammaTree,
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
}

// ----------------------------------------------------------------

// custom execute 6 resources/template_example.template.go

// Command Example: custom execute 6 resources/template_example.template.go
// Command Expression:
// - custom                                : command const content
// - execute                               : command const content
// - 6                                     : specify executor
// - resources/template_example.template.go: specify a file to analyze

var TemplateType string = "template.TypeName"
var TemplateFileKey string = ".template."

type GoTemplateFunctionAnalysis struct {
	Analysis *GoFunctionAnalysis
	// DeductionAnalysis *GoFunctionAnalysis
	// ToDeductionTemplateParamIndexList  []int
	ToDeductionTemplateReturnIndexList []int
	ParamDeductionGroupMap             map[int]map[int]goValueType
	// ParamDeductionAnalysisList         []*GoTemplateFunctionAnalysis
	ReturnDeductionGroupMap map[int]map[int]goValueType
	DeductionFunctionMap    map[int]string // deduction group : deduction function
}

// GoCommandToolTemplater go 语言命令行工具：模板代码生成器
func GoCommandToolTemplater(paramList []string) {
	checkTypeAtomicExpression()

	if len(paramList) < 1 {
		ui.OutputErrorInfo(ui.CMDCustomExecutorHasNotEnoughParam, 5)
		return
	}
	fileRelativePath := paramList[0]

	abs, _ := filepath.Abs(".")
	fileABS := path.Join(abs, fileRelativePath)
	fileRootPath := path.Dir(fileABS)
	goAnalysis, analyzeError := analyzeGo(fileRootPath, []string{fileABS})
	if analyzeError != nil {
		ui.OutputWarnInfo(ui.CMDAnalyzeOccursError, analyzeError)
		return
	}

	utility2.TestOutput(ui.CommonNote2)
	utility2.TestOutput("goAnalysis = %+v", goAnalysis)

	// for fileNo, goFileAnalysis := range goAnalysis.FileNoAnalysisMap {
	// 	utility2.TestOutput("fileNo = %v, goFileAnalysis.FilePath = %v", fileNo, goFileAnalysis.FilePath)

	// 	templateFunctionAnalysisMap := make(map[string]*GoTemplateFunctionAnalysis)

	// 	// get template functions
	// 	// utility2.TestOutput("get template functions from file analysis")
	// 	for functionName, functionAnalysis := range goFileAnalysis.FunctionMap {
	// 		// utility2.TestOutput("function name: %v", functionName)

	// 		toDeductionTemplateParamIndexList := make([]int, 0)
	// 		toDeductionTemplateReturnIndexList := make([]int, 0)

	// 		for paramIndex, paramValue := range functionAnalysis.ParamsMap {
	// 			if paramValue.Type == TemplateType {
	// 				// utility2.TestOutput("in param list, template value %v", paramValue.Name)
	// 				toDeductionTemplateParamIndexList = append(toDeductionTemplateParamIndexList, paramIndex)
	// 			}
	// 		}

	// 		for returnIndex, returnValue := range functionAnalysis.ReturnMap {
	// 			if returnValue.Type == TemplateType {
	// 				// utility2.TestOutput("in return list, template value %v", returnValue.Name)
	// 				toDeductionTemplateReturnIndexList = append(toDeductionTemplateReturnIndexList, returnIndex)
	// 			}
	// 		}

	// 		// utility2.TestOutput("to deduction template param index list: %v", toDeductionTemplateParamIndexList)
	// 		// utility2.TestOutput("to deduction template return index list: %v", toDeductionTemplateReturnIndexList)

	// 		if len(toDeductionTemplateParamIndexList) != 0 || len(toDeductionTemplateReturnIndexList) != 0 {
	// 			// utility2.TestOutput("function %v is template function, need deduction", functionName)
	// 			templateFunctionAnalysisMap[functionName] = &GoTemplateFunctionAnalysis{
	// 				Analysis: functionAnalysis,
	// 				// DeductionAnalysis: &GoFunctionAnalysis{
	// 				// 	Class:          functionAnalysis.Class,
	// 				// 	ClassValue:     functionAnalysis.ClassValue,
	// 				// 	ClassValueType: functionAnalysis.ClassValueType,
	// 				// },
	// 				// ToDeductionTemplateParamIndexList:  toDeductionTemplateParamIndexList,
	// 				ToDeductionTemplateReturnIndexList: toDeductionTemplateReturnIndexList,
	// 				ParamDeductionGroupMap:             make(map[int]map[int]goValueType, 0),
	// 				// ParamDeductionAnalysisList:         make([]*GoTemplateFunctionAnalysis, 0),
	// 				ReturnDeductionGroupMap: make(map[int]map[int]goValueType, 0),
	// 				DeductionFunctionMap:    make(map[int]string),
	// 			}
	// 		}
	// 	}

	// 	// utility2.TestOutput("deduction function param type")
	// 	for _, functionAnalysis := range goFileAnalysis.FunctionMap {
	// 		// inner package call function
	// 		paramDeduction(functionAnalysis.InnerPackageCallMap, templateFunctionAnalysisMap, func(callFunction string, callParam []string) {
	// 			// utility2.TestOutput(ui.CommonNote2)
	// 			// utility2.TestOutput("%v call inner package template function %v, param %v", functionName, callFunction, callParam)
	// 		}, func(param string, paramType goValueType) {
	// 			// utility2.TestOutput("param deduction: param %v to type %v", param, paramType)
	// 		})

	// 		// outer package call function
	// 		for _, callFunctionMap := range functionAnalysis.OuterPackageCallMap {
	// 			paramDeduction(callFunctionMap, templateFunctionAnalysisMap, func(callFunction string, callParam []string) {
	// 				// utility2.TestOutput(ui.CommonNote2)
	// 				// utility2.TestOutput("%v call package %v template function %v, param %v", functionName, callFrom, callFunction, callParam)
	// 			}, func(param string, paramType goValueType) {
	// 				// utility2.TestOutput("param deduction: param %v to type %v", param, paramType)
	// 			})
	// 		}

	// 		// member call function
	// 		for _, callFunctionMap := range functionAnalysis.MemberCallMap {
	// 			paramDeduction(callFunctionMap, templateFunctionAnalysisMap, func(callFunction string, callParam []string) {
	// 				// utility2.TestOutput(ui.CommonNote2)
	// 				// utility2.TestOutput("%v call member %v template function %v, param %v", functionName, callFrom, callFunction, callParam)
	// 			}, func(param string, paramType goValueType) {
	// 				// utility2.TestOutput("param deduction: param %v to type %v", param, paramType)
	// 			})
	// 		}
	// 	}

	// 	// utility2.TestOutput("deduction function return type")
	// 	for _, templateFunctionAnalysis := range templateFunctionAnalysisMap {
	// 		// utility2.TestOutput(ui.CommonNote2)
	// 		templateFunctionAnalysis.ReturnDeductionGroupMap = returnDeduction(templateFunctionAnalysis)
	// 	}

	// 	// utility2.TestOutput("generate new function definition with deduction type")
	// 	for functionName, templateFunctionAnalysis := range templateFunctionAnalysisMap {
	// 		// utility2.TestOutput("function %v deduction", functionName)
	// 		for deductionGroup, paramIndexDeductionMap := range templateFunctionAnalysis.ParamDeductionGroupMap {
	// 			// utility2.TestOutput("deductionGroup = %v", deductionGroup)
	// 			deductionFunction := goFunctionDefintion

	// 			// replace class
	// 			rpClassContent := goTemplateRPFunctionClass
	// 			if len(templateFunctionAnalysis.Analysis.ClassValue) != 0 && len(templateFunctionAnalysis.Analysis.ClassValueType) != 0 {
	// 				rpClassContent = fmt.Sprintf("%v ", rpClassContent)
	// 				rpClassContent = strings.Replace(rpClassContent, goTemplateRPFunctionClassValueKey, templateFunctionAnalysis.Analysis.ClassValue, 1)
	// 				rpClassContent = strings.Replace(rpClassContent, goTemplateRPFunctionClassValueTypeKey, templateFunctionAnalysis.Analysis.ClassValueType, 1)
	// 				deductionFunction = strings.Replace(deductionFunction, goTemplateRPFunctionClassKey, rpClassContent, 1)
	// 			} else {
	// 				deductionFunction = strings.Replace(deductionFunction, goTemplateRPFunctionClassKey, "", 1)
	// 			}

	// 			// replace function name
	// 			deductionFunction = strings.Replace(deductionFunction, goTemplateRPFunctionNameKey, fmt.Sprintf("%vType%v", functionName, deductionGroup), 1)

	// 			// replace param list
	// 			// utility2.TestOutput("paramIndexDeductionMap = %v", paramIndexDeductionMap)
	// 			rpParamListContent := ""
	// 			for index, deductionType := range paramIndexDeductionMap {
	// 				paramValueTypeContent := goTemplateRPFunctionParamList
	// 				paramValueTypeContent = strings.Replace(paramValueTypeContent, goTemplateRPFunctionParamValueKey, templateFunctionAnalysis.Analysis.ParamsMap[index].Name, 1)
	// 				paramValueTypeContent = strings.Replace(paramValueTypeContent, goTemplateRPFunctionParamTypeKey, typeStringMap[deductionType], 1)
	// 				if rpParamListContent == "" {
	// 					rpParamListContent = paramValueTypeContent
	// 				} else {
	// 					rpParamListContent = fmt.Sprintf("%v, %v", rpParamListContent, paramValueTypeContent)
	// 				}
	// 			}
	// 			deductionFunction = strings.Replace(deductionFunction, goTemplateRPFunctionParamListKey, rpParamListContent, 1)

	// 			// replace return list
	// 			returnDeductionMap := templateFunctionAnalysis.ReturnDeductionGroupMap[deductionGroup]
	// 			// utility2.TestOutput("returnDeductionMap = %v", returnDeductionMap)
	// 			rpReturnListContent := ""
	// 			for deductionIndex, deductionType := range returnDeductionMap {
	// 				returnValueName := templateFunctionAnalysis.Analysis.ReturnMap[deductionIndex].Name
	// 				returnValueTypeContent := goTemplateRPFunctionReturnList
	// 				returnValueTypeContent = strings.Replace(returnValueTypeContent, goTemplateRPFunctionReturnValueKey, returnValueName, 1)
	// 				returnValueTypeContent = strings.Replace(returnValueTypeContent, goTemplateRPFunctionReturnTypeKey, typeStringMap[deductionType], 1)
	// 				if rpReturnListContent == "" {
	// 					rpReturnListContent = returnValueTypeContent
	// 				} else {
	// 					rpReturnListContent = fmt.Sprintf("%v, %v", rpReturnListContent, returnValueTypeContent)
	// 				}
	// 			}
	// 			if len(returnDeductionMap) > 1 {
	// 				rpReturnListContent = fmt.Sprintf(" (%v)", rpReturnListContent)
	// 			}
	// 			deductionFunction = strings.Replace(deductionFunction, goTemplateRPFunctionReturnListKey, rpReturnListContent, 1)

	// 			// utility2.TestOutput("deductionFunction = %v", deductionFunction)
	// 			templateFunctionAnalysis.DeductionFunctionMap[deductionGroup] = deductionFunction
	// 		}
	// 		// utility2.TestOutput(ui.CommonNote2)
	// 	}

	// 	// utility2.TestOutput("output new function definition to file")
	// 	newFileName := strings.ReplaceAll(filename, TemplateFileKey, ".")
	// 	newFile, createFileError := utility.CreateFile(newFileName)
	// 	if createFileError != nil {
	// 		ui.OutputErrorInfo(ui.CommonError4, createFileError)
	// 		return
	// 	}
	// 	defer newFile.Close()
	// 	fileContent := ""
	// 	for _, templateFunctionAnalysis := range templateFunctionAnalysisMap {
	// 		for _, deductionFunction := range templateFunctionAnalysis.DeductionFunctionMap {
	// 			function := strings.Replace(deductionFunction, goTemplateRPFunctionBodyKey, string(templateFunctionAnalysis.Analysis.BodyContent), 1)
	// 			// utility2.TestOutput("%v function = \n%v\n", deductionGroup, utility3.TrimSpaceLine(function))
	// 			// utility2.TestOutput(ui.CommonNote2)
	// 			fileContent = fmt.Sprintf("%v\n%v\n", fileContent, utility3.TrimSpaceLine(function))
	// 		}
	// 	}
	// 	newFile.WriteString(fileContent)
	// }

}

func paramDeduction(callFunctionAnalysisMap map[string][]*GoFunctionCallAnalysis, templateFunctionAnalysisMap map[string]*GoTemplateFunctionAnalysis, testLog1 func(string, []string), testLog2 func(string, goValueType)) {
	for callFunction, callAnalysisList := range callFunctionAnalysisMap {
		if templateFunctionAnalysis, isTemplateFunction := templateFunctionAnalysisMap[callFunction]; isTemplateFunction {
			for deductionGroup, callAnalysis := range callAnalysisList {
				testLog1(callFunction, callAnalysis.ParamList)
				paramDeductionMap := make(map[int]goValueType)
				for index, param := range callAnalysis.ParamList {
					paramType := goValueTypeDeduction(param)
					testLog2(param, paramType)
					paramDeductionMap[index] = paramType
				}
				templateFunctionAnalysis.ParamDeductionGroupMap[deductionGroup] = paramDeductionMap
			}
		}
	}
}

// after param deduction
// for each param deduction result
// it will deduction more than one gourp
func returnDeduction(afterParamDeductionTemplateFunctionAnalysis *GoTemplateFunctionAnalysis) map[int]map[int]goValueType {
	returnDeductionGroup := make(map[int]map[int]goValueType)
	for deductionGroup := range afterParamDeductionTemplateFunctionAnalysis.ParamDeductionGroupMap {
		returnDeductionMap := make(map[int]goValueType, 0)
		for _, deductionIndex := range afterParamDeductionTemplateFunctionAnalysis.ToDeductionTemplateReturnIndexList {
			// TODO:
			returnDeductionMap[deductionIndex] = tUnknown
			// returnValueName := afterParamDeductionTemplateFunctionAnalysis.Analysis.ReturnMap[deductionIndex].Name
			// utility2.TestOutput("deduction function %v No.%v return value(name: %v) type to %v", afterParamDeductionTemplateFunctionAnalysis.Analysis.Name, deductionIndex, returnValueName, tUnknown)
		}
		returnDeductionGroup[deductionGroup] = returnDeductionMap
	}
	return returnDeductionGroup
}

var goFunctionDefintion string = "func RP_CLASSRP_FUNCTION_NAME(RP_PARAM_LIST)RP_RETURN_LIST RP_FUNCTION_BODY"

var goTemplateRPFunctionClassKey string = "RP_CLASS"
var goTemplateRPFunctionClass string = "(CLASS_VALUE CLASS_VALUE_TYPE)"

var goTemplateRPFunctionClassValueKey string = "CLASS_VALUE"
var goTemplateRPFunctionClassValueTypeKey string = "CLASS_VALUE_TYPE"

var goTemplateRPFunctionNameKey string = "RP_FUNCTION_NAME"

var goTemplateRPFunctionParamListKey string = "RP_PARAM_LIST"
var goTemplateRPFunctionParamList string = "PARAM_VALUE PARAM_TYPE"
var goTemplateRPFunctionParamValueKey = "PARAM_VALUE"
var goTemplateRPFunctionParamTypeKey string = "PARAM_TYPE"

var goTemplateRPFunctionReturnListKey string = "RP_RETURN_LIST"
var goTemplateRPFunctionReturnList string = "RETURN_VALUE RETURN_TYPE"
var goTemplateRPFunctionReturnValueKey = "RETURN_VALUE"
var goTemplateRPFunctionReturnTypeKey string = "RETURN_TYPE"

var goTemplateRPFunctionBodyKey = "RP_FUNCTION_BODY"

type goValueType int

const (
	tUnknown    goValueType = iota // 0
	tInt                           // 1
	tInt8                          // 2
	tInt16                         // 3
	tInt32                         // 4
	tInt64                         // 5
	tUint                          // 6
	tUint8                         // 7
	tUint16                        // 8
	tUint32                        // 9
	tUint64                        // 10
	tFloat32                       // 11
	tFloat64                       // 12
	tComplex64                     // 13
	tComplex128                    // 14
	tBool                          // 15
	tString                        // 16
	tUintptr                       // 17
	fFunc                          // 18
)

var atomicExpressionEnumGoValueTypeMap map[global.AtomicExpressionEnum]goValueType
var goTypeConvertRegexp *regexp.Regexp
var valueTypeStringMap map[string]goValueType
var typeStringMap map[goValueType]string

func checkTypeAtomicExpression() {
	atomicExpressionEnumGoValueTypeMap = map[global.AtomicExpressionEnum]goValueType{
		global.AEInteger: tInt,
		global.AEFloat:   tFloat64,
		global.AEComplex: tComplex128,
	}
	for toCheckTypeAtomicExpressionEnum := range atomicExpressionEnumGoValueTypeMap {
		if _, hasRegexp := regexps.AtomicExpressionEnumRegexpMap[toCheckTypeAtomicExpressionEnum]; !hasRegexp {
			ui.OutputWarnInfo(ui.CommonError16, toCheckTypeAtomicExpressionEnum)
		}
	}

	if goTypeConvertRegexp = regexps.GetRegexpByTemplateEnum(global.GoTypeConvertTemplate); goTypeConvertRegexp == nil {
		ui.OutputWarnInfo(ui.CommonError18, global.GoTypeConvertTemplate)
	}

	valueTypeStringMap = map[string]goValueType{
		"int":        tInt,
		"int8":       tInt8,
		"int16":      tInt16,
		"int32":      tInt32,
		"int64":      tInt64,
		"uint":       tUint,
		"uint8":      tUint8,
		"uint16":     tUint16,
		"uint32":     tUint32,
		"uint64":     tUint64,
		"float32":    tFloat32,
		"float64":    tFloat64,
		"complex64":  tComplex64,
		"complex128": tComplex128,
		"bool":       tBool,
		"string":     tString,
		"uintptr":    tUintptr,
	}

	typeStringMap = make(map[goValueType]string)
	typeStringMap[tUnknown] = TemplateType
	for typeString, typeEnum := range valueTypeStringMap {
		typeStringMap[typeEnum] = typeString
	}
}

func goValueTypeDeduction(valueString string) goValueType {
	for toCheckTypeAtomicExpressionEnum, goValueType := range atomicExpressionEnumGoValueTypeMap {
		if regexps.AtomicExpressionEnumRegexpMap[toCheckTypeAtomicExpressionEnum].MatchString(valueString) {
			return goValueType
		}
	}

	if goTypeConvertRegexp.MatchString(valueString) {
		identifier := goTypeConvertRegexp.ReplaceAllString(valueString, "$IDENTIFIER")
		if goValueType, isGoType := valueTypeStringMap[identifier]; isGoType {
			return goValueType
		}
	}

	return tUnknown
}

// ----------------------------------------------------------------

// custom execute 7 resources/template_example.template.go true

// Command Example: custom execute 7 resources/template_example.template.go
// Command Expression:
// - custom                                : command const content
// - execute                               : command const content
// - 7                                     : specify executor
// - resources/template_example.template.go: specify a file to split

type goFileLineState int

func (s goFileLineState) String() string {
	switch s {
	case lineStateNone:
		return "None"
	case lineStateSpace:
		return "Space"
	case lineStateComment:
		return "Comment"
	case lineStatePackageScope:
		return "Package Scope"
	case lineStateMultiLineImportScope:
		return "Multi-Line Import Scope"
	case lineStateSingleLineImportScope:
		return "Single-Line Import Scope"
	case lineStatePackageVariableScope:
		return "Package Variable Scope"
	case lineStateInterfaceScope:
		return "Interface Scope"
	case lineStateStructScope:
		return "Struct Scope"
	case lineStateFunctionScope:
		return "Function Scope"
	case lineStateMemberFunctionScope:
		return "Member Function Scope"
	case lineStateTypeRenameScope:
		return "Type Rename Scope"
	case lineStateMultiLineConstScope:
		return "Multi-Line Const Scope"
	case lineStateSingleLineConstScope:
		return "Single-Line Const Scope"
	}
	return ""
}

const (
	lineStateNone    goFileLineState = 1 << iota // 0000 0000
	lineStateSpace                               // 0000 0001
	lineStateComment                             // 0000 0010
	lineStateInScope                             // 0000 0100
	lineStateTODO2
	lineStateTODO3
	lineStateTODO4
	lineStateTODO5
	lineStateTODO6
	lineStatePackageScope          // 0000 0000 0000 0000 0000 0001 0000 0000
	lineStateMultiLineImportScope  // 0000 0000 0000 0000 0000 0010 0000 0000
	lineStateSingleLineImportScope // 0000 0000 0000 0000 0000 0100 0000 0000
	lineStatePackageVariableScope  // 0000 0000 0000 0000 0000 1000 0000 0000
	lineStateInterfaceScope        // 0000 0000 0000 0000 0001 0000 0000 0000
	lineStateStructScope           // 0000 0000 0000 0000 0010 0000 0000 0000
	lineStateFunctionScope         // 0000 0000 0000 0000 0100 0000 0000 0000
	lineStateMemberFunctionScope   // 0000 0000 0000 0000 1000 0000 0000 0000
	lineStateTypeRenameScope       // 0000 0000 0000 0001 0000 0000 0000 0000
	lineStateMultiLineConstScope   // 0000 0000 0000 0010 0000 0000 0000 0000
	lineStateSingleLineConstScope  // 0000 0000 0000 0100 0000 0000 0000 0000
)

// split file content to different scopes
// scope has some attribute: struct scope
// scope has its type:
// - package
// - import
// - package variable/constant
// - struct/interface
// - function

type goFileScopeType int

const (
	scopePackage goFileScopeType = iota + 1
	scopeMultiLineImport
	scopeSignleLineImport
	scopePackageVariable
	scopeInterface
	scopeStruct
	scopeFunction
	scopeMemberFunction
	scopeTypeRename
	scopeMultiLineConst
	scopeSingleLineConst
)

type scope struct {
	LineStart int
	LineEnd   int
	ScopeType goFileScopeType
	Content   string
}

func (s scope) isOneLineScope() bool {
	return s.LineStart == s.LineEnd
}

// GoPackageScope go 包切分结果
type GoPackageScope struct {
	Package                  *scope
	MultiLineImport          *scope
	SingleLineImport         []*scope
	PackageVariable          []*scope
	InterfaceDefinition      map[string]*scope
	StructDefinition         map[string]*scope
	FunctionDefinition       map[string]*scope
	MemberFunctionDefinition map[string]map[string]*scope
	TypeRename               map[string]map[string]*scope
	MultiLineConst           []*scope
	SingleLineConst          map[string]*scope
}

var goSplitterPackageSubMatchNameIndexMap map[string]int
var goSplitterSingleLineImportSubMatchNameIndexMap map[string]int
var goSplitterMultiLineImportContentSubMatchNameIndexMap map[string]int
var goSplitterPackageVariableSubMatchNameIndexMap map[string]int
var goSplitterInterfaceSubMatchNameIndexMap map[string]int
var goSplitterStructSubMatchNameIndexMap map[string]int
var goSplitterFunctionSubMatchNameIndexMap map[string]int
var goSplitterMemberFunctionSubMatchNameIndexMap map[string]int
var goSplitterTypeRenameSubMatchNameIndexMap map[string]int
var goSplitterMultiLineConstSubMatchNameIndexMap map[string]int
var goSplitterSingleLineConstSubMatchNameIndexMap map[string]int

// GoFileSplitter go 文件切分示例
func GoFileSplitter(paramList []string) {
	if len(paramList) < 1 {
		ui.OutputErrorInfo(ui.CMDCustomExecutorHasNotEnoughParam, 5)
		return
	}
	filename := paramList[0]
	outputString := paramList[1]
	var output bool = false
	if strings.ToLower(outputString) == "true" {
		output = true
	}

	if !checkGoSplitterRegexp() {
		return
	}

	SplitGoFile(filename, output)
}

// SplitGoFile go 文件切分
func SplitGoFile(filename string, output bool) *GoPackageScope {
	gps := &GoPackageScope{
		SingleLineImport:         make([]*scope, 0),
		PackageVariable:          make([]*scope, 0),
		InterfaceDefinition:      make(map[string]*scope),
		StructDefinition:         make(map[string]*scope),
		FunctionDefinition:       make(map[string]*scope),
		MemberFunctionDefinition: make(map[string]map[string]*scope),
		TypeRename:               make(map[string]map[string]*scope),
		MultiLineConst:           make([]*scope, 0),
		SingleLineConst:          make(map[string]*scope),
	}

	var lineState goFileLineState = lineStateNone
	var keyInterface interface{}

	if output {
		utility2.TestOutput("split file content to line text one by one")
	}
	lineIndex := 0
	utility.ReadFileLineOneByOne(filename, func(line string) bool {
		lineIndex++

		if len(line) == 0 {
			return true
		}

		// if lineIndex > 95 {
		// 	return false
		// }

		if output {
			utility2.TestOutput("No.%v line content = |%v|", lineIndex, line)
		}

		if lineState == lineStateNone {
			lineState = getLineState(line)
		}
		if output {
			utility2.TestOutput("line state is: %v", lineState.String())
			utility2.TestOutput("key is: %v", keyInterface)
		}

		switch lineState {
		case lineStatePackageScope:
			lineState = packageScope(line, lineIndex, gps, lineStatePackageScope, lineStateNone)
		case lineStateMultiLineImportScope:
			lineState = multiLineImportScope(line, lineIndex, gps, lineStateMultiLineImportScope, lineStateNone)
		case lineStateSingleLineImportScope:
			lineState = signleLineImportScope(line, lineIndex, gps, lineStateSingleLineImportScope, lineStateNone)
		case lineStatePackageVariableScope:
			lineState = packageVariableScope(line, lineIndex, gps, lineStatePackageVariableScope, lineStateNone)
		case lineStateInterfaceScope:
			lineState, keyInterface = interfaceScope(line, lineIndex, gps, keyInterface, lineStateInterfaceScope, lineStateNone)
		case lineStateStructScope:
			lineState, keyInterface = structScope(line, lineIndex, gps, keyInterface, lineStateStructScope, lineStateNone)
		case lineStateFunctionScope:
			lineState, keyInterface = functionScope(line, lineIndex, gps, keyInterface, lineStateFunctionScope, lineStateNone)
		case lineStateMemberFunctionScope:
			lineState, keyInterface = memberFunctionScope(line, lineIndex, gps, keyInterface, lineStateMemberFunctionScope, lineStateNone)
		case lineStateTypeRenameScope:
			lineState = typeRenameScope(line, lineIndex, gps, lineStateTypeRenameScope, lineStateNone)
		case lineStateMultiLineConstScope:
			lineState, keyInterface = multiLineConstScope(line, lineIndex, gps, keyInterface, lineStateMultiLineConstScope, lineStateNone)
		case lineStateSingleLineConstScope:
			lineState = singleLineConstScope(line, lineIndex, gps, lineStateSingleLineConstScope, lineStateNone)
		case lineStateNone:
		default:
			ui.OutputNoteInfo("unknown line state %v", lineState)
		}

		if output {
			utility2.TestOutput(ui.CommonNote2)
		}

		return true
	})

	if gps.Package != nil {
		ui.OutputNoteInfo(ui.CommonNote2)
		ui.OutputNoteInfo("package scope = |%v|", gps.Package.Content)
	}
	if gps.MultiLineImport != nil {
		ui.OutputNoteInfo(ui.CommonNote2)
		ui.OutputNoteInfo("multi-line import scope = |%v|", gps.MultiLineImport.Content)
	}
	if gps.SingleLineImport != nil {
		ui.OutputNoteInfo(ui.CommonNote2)
		ui.OutputNoteInfo("single-line import scope:")
		for _, scopeData := range gps.SingleLineImport {
			ui.OutputNoteInfo("|%v|", scopeData.Content)
		}
	}
	if gps.PackageVariable != nil {
		ui.OutputNoteInfo(ui.CommonNote2)
		ui.OutputNoteInfo("package variable scope:")
		for _, scopeData := range gps.PackageVariable {
			ui.OutputNoteInfo("|%v|", scopeData.Content)
		}
	}
	if gps.InterfaceDefinition != nil {
		ui.OutputNoteInfo(ui.CommonNote2)
		ui.OutputNoteInfo("interface scope:")
		for _, scopeData := range gps.InterfaceDefinition {
			ui.OutputNoteInfo("|%v|", scopeData.Content)
		}
	}
	if gps.StructDefinition != nil {
		ui.OutputNoteInfo(ui.CommonNote2)
		ui.OutputNoteInfo("struct scope:")
		for _, scopeData := range gps.StructDefinition {
			ui.OutputNoteInfo("|%v|", scopeData.Content)
		}
	}
	if gps.FunctionDefinition != nil {
		ui.OutputNoteInfo(ui.CommonNote2)
		ui.OutputNoteInfo("function scope:")
		for _, scopeData := range gps.FunctionDefinition {
			ui.OutputNoteInfo("|%v|", scopeData.Content)
		}
	}
	if gps.MemberFunctionDefinition != nil {
		ui.OutputNoteInfo(ui.CommonNote2)
		ui.OutputNoteInfo("member function scope:")
		for _, functionMap := range gps.MemberFunctionDefinition {
			for _, scopeData := range functionMap {
				ui.OutputNoteInfo("|%v|", scopeData.Content)
			}
		}
	}
	if gps.TypeRename != nil {
		ui.OutputNoteInfo(ui.CommonNote2)
		ui.OutputNoteInfo("type rename scope:")
		for _, renameMap := range gps.TypeRename {
			for _, scopeData := range renameMap {
				ui.OutputNoteInfo("|%v|", scopeData.Content)
			}
		}
	}
	if gps.MultiLineConst != nil {
		ui.OutputNoteInfo(ui.CommonNote2)
		ui.OutputNoteInfo("multi-line const scope:")
		for _, scopeData := range gps.MultiLineConst {
			ui.OutputNoteInfo("|%v|", scopeData.Content)
		}
	}
	if gps.SingleLineConst != nil {
		ui.OutputNoteInfo(ui.CommonNote2)
		ui.OutputNoteInfo("const value scope:")
		for _, scopeData := range gps.SingleLineConst {
			ui.OutputNoteInfo("|%v|", scopeData.Content)
		}
	}

	return gps
}

func getLineState(line string) goFileLineState {
	utility2.TestOutput("get line state by |%v|", line)
	if regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopePackageTemplate).MatchString(line) {
		return lineStatePackageScope
	} else if regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeMultiLineImportStartTemplate).MatchString(line) {
		return lineStateMultiLineImportScope
	} else if regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeSingleLineImportTemplate).MatchString(line) {
		return lineStateSingleLineImportScope
	} else if regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopePackageVariableTemplate).MatchString(line) {
		return lineStatePackageVariableScope
	} else if regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeInterfaceTemplate).MatchString(line) {
		return lineStateInterfaceScope
	} else if regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeStructTemplate).MatchString(line) {
		return lineStateStructScope
	} else if isMatchGoScopeFunction(line) {
		return lineStateFunctionScope
	} else if regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeMemberFunctionTemplate).MatchString(line) {
		return lineStateMemberFunctionScope
	} else if regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeTypeRenameTemplate).MatchString(line) {
		return lineStateTypeRenameScope
	} else if regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeMultiLineConstStartTemplate).MatchString(line) {
		return lineStateMultiLineConstScope
	} else if regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeSingleLineConstTemplate).MatchString(line) {
		return lineStateSingleLineConstScope
	}
	return lineStateNone
}

func isMatchGoScopeFunction(content string) bool {
	replacedContent, _ := utility.ReplaceToUniqueString(content, global.GoKeywordEmptyInterface)
	return regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeFunctionTemplate).MatchString(replacedContent)
}

func packageScope(line string, lineIndex int, gps *GoPackageScope, continueState, endState goFileLineState) goFileLineState {
	if gps.Package == nil {
		gps.Package = &scope{
			LineStart: lineIndex,
			LineEnd:   lineIndex,
			ScopeType: scopePackage,
			Content:   line,
		}
	}
	return endState
}

func multiLineImportScope(line string, lineIndex int, gps *GoPackageScope, continueState, endState goFileLineState) goFileLineState {
	// scope begin
	if gps.MultiLineImport == nil {
		gps.MultiLineImport = &scope{
			LineStart: lineIndex,
			LineEnd:   lineIndex,
			ScopeType: scopeMultiLineImport,
			Content:   line,
		}
		return continueState
	}

	// scope content
	gps.MultiLineImport.Content = fmt.Sprintf("%v\n%v", gps.MultiLineImport.Content, line)

	// scope end
	if regexps.AtomicExpressionEnumRegexpMap[global.AEGoFileSplitterScopeEnd].MatchString(line) {
		gps.MultiLineImport.LineEnd = lineIndex
		return endState
	}

	return continueState
}

func signleLineImportScope(line string, lineIndex int, gps *GoPackageScope, continueState, endState goFileLineState) goFileLineState {
	gps.SingleLineImport = append(gps.SingleLineImport, &scope{
		LineStart: lineIndex,
		LineEnd:   lineIndex,
		ScopeType: scopeSignleLineImport,
		Content:   line,
	})
	return endState
}

func packageVariableScope(line string, lineIndex int, gps *GoPackageScope, continueState, endState goFileLineState) goFileLineState {
	gps.PackageVariable = append(gps.PackageVariable, &scope{
		LineStart: lineIndex,
		LineEnd:   lineIndex,
		ScopeType: scopePackageVariable,
		Content:   line,
	})
	return endState
}

func interfaceScope(line string, lineIndex int, gps *GoPackageScope, keyInterface interface{}, continueState, endState goFileLineState) (goFileLineState, interface{}) {
	var key string
	if keyInterface != nil {
		var ok bool
		key, ok = keyInterface.(string)
		if !ok {
			ui.OutputErrorInfo(ui.CommonError19, "keyInterface", "string")
			return endState, nil
		}
	}

	utility2.TestOutput("key = %v", key)

	// scope begin
	if len(key) == 0 {
		var interfaceKey string
		var scopeEnd string
		utility2.TestOutput("goSplitterInterfaceSubMatchNameIndexMap = %v", goSplitterInterfaceSubMatchNameIndexMap)
		for _, subMatchList := range regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeInterfaceTemplate).FindAllStringSubmatch(line, -1) {
			if index, hasIndex := goSplitterInterfaceSubMatchNameIndexMap["NAME"]; hasIndex {
				interfaceKey = strings.TrimSpace(subMatchList[index])
			}
			if index, hasIndex := goSplitterInterfaceSubMatchNameIndexMap["SCOPE_END"]; hasIndex {
				scopeEnd = strings.TrimSpace(subMatchList[index])
			}
		}

		gps.InterfaceDefinition[interfaceKey] = &scope{
			LineStart: lineIndex,
			LineEnd:   lineIndex,
			ScopeType: scopeInterface,
			Content:   line,
		}
		// one line interface
		if len(scopeEnd) != 0 {
			utility2.TestOutput("interface %v is one line interface", interfaceKey)
			return endState, nil
		}
		utility2.TestOutput("interface line = %v", line)
		return continueState, interfaceKey
	}

	// scope content
	gps.InterfaceDefinition[key].Content = fmt.Sprintf("%v\n%v", gps.InterfaceDefinition[key].Content, line)

	// scope end
	if regexps.AtomicExpressionEnumRegexpMap[global.AEGoFileSplitterScopeEnd].MatchString(line) {
		gps.InterfaceDefinition[key].LineEnd = lineIndex
		utility2.TestOutput("interface scope is end at line %v", lineIndex)
		return endState, nil
	}

	return continueState, key
}

func structScope(line string, lineIndex int, gps *GoPackageScope, keyInterface interface{}, continueState, endState goFileLineState) (goFileLineState, interface{}) {
	var key string
	if keyInterface != nil {
		var ok bool
		key, ok = keyInterface.(string)
		if !ok {
			ui.OutputErrorInfo(ui.CommonError19, "keyInterface", "string")
			return endState, nil
		}
	}

	// scope begin
	if len(key) == 0 {
		var structKey string
		var scopeEnd string
		for _, subMatchList := range regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeStructTemplate).FindAllStringSubmatch(line, -1) {
			if structNameIndex, hasIndex := goSplitterStructSubMatchNameIndexMap["NAME"]; hasIndex {
				structKey = strings.TrimSpace(subMatchList[structNameIndex])
			}
			if scopeEndIndex, hasIndex := goSplitterStructSubMatchNameIndexMap["SCOPE_END"]; hasIndex {
				scopeEnd = strings.TrimSpace(subMatchList[scopeEndIndex])
			}
		}
		gps.StructDefinition[structKey] = &scope{
			LineStart: lineIndex,
			LineEnd:   lineIndex,
			ScopeType: scopeStruct,
			Content:   line,
		}
		// one line struct
		if len(scopeEnd) != 0 {
			return endState, nil
		}
		return continueState, structKey
	}

	// scope content
	gps.StructDefinition[key].Content = fmt.Sprintf("%v\n%v", gps.StructDefinition[key].Content, line)

	// scope end
	if regexps.AtomicExpressionEnumRegexpMap[global.AEGoFileSplitterScopeEnd].MatchString(line) {
		gps.StructDefinition[key].LineEnd = lineIndex
		return endState, nil
	}

	return continueState, key
}

func functionScope(line string, lineIndex int, gps *GoPackageScope, keyInterface interface{}, continueState, endState goFileLineState) (goFileLineState, interface{}) {
	var key string
	if keyInterface != nil {
		var ok bool
		key, ok = keyInterface.(string)
		if !ok {
			ui.OutputErrorInfo(ui.CommonError19, "keyInterface", "string")
			return endState, nil
		}
	}

	replacedLine, _ := utility.ReplaceToUniqueString(line, global.GoKeywordEmptyInterface)
	// utility2.TestOutput("replace %v to %v", global.GoKeywordEmptyInterface, replacedString)
	// utility2.TestOutput("replaced line = %v", replacedLine)

	// scope begin
	if len(key) == 0 {
		var functionKey string
		// var body string
		// var definition string
		var scopeEnd string
		for _, subMatchList := range regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeFunctionTemplate).FindAllStringSubmatch(replacedLine, -1) {
			if functionNameIndex, hasIndex := goSplitterFunctionSubMatchNameIndexMap["NAME"]; hasIndex {
				functionKey = strings.TrimSpace(subMatchList[functionNameIndex])
			}
			// if contentIndex, hasIndex := goSplitterFunctionSubMatchNameIndexMap["BODY"]; hasIndex {
			// 	body = strings.TrimSpace(subMatchList[contentIndex])
			// }
			// if index, hasIndex := goSplitterFunctionSubMatchNameIndexMap["DEFINITION"]; hasIndex {
			// 	definition = strings.TrimSpace(subMatchList[index])
			// }
			if scopeEndIndex, hasIndex := goSplitterFunctionSubMatchNameIndexMap["SCOPE_END"]; hasIndex {
				scopeEnd = strings.TrimSpace(subMatchList[scopeEndIndex])
			}
		}
		// utility2.TestOutput("functionKey = %v", functionKey)
		// utility2.TestOutput("body = %v", body)
		// utility2.TestOutput("definition = %v", definition)
		// utility2.TestOutput("scopeEnd = %v", scopeEnd)

		gps.FunctionDefinition[functionKey] = &scope{
			LineStart: lineIndex,
			LineEnd:   lineIndex,
			ScopeType: scopeFunction,
			Content:   line,
		}
		// one line function
		if len(scopeEnd) != 0 {
			return endState, nil
		}

		return continueState, functionKey
	}

	// scope content
	gps.FunctionDefinition[key].Content = fmt.Sprintf("%v\n%v", gps.FunctionDefinition[key].Content, replacedLine)

	// scope end
	if regexps.AtomicExpressionEnumRegexpMap[global.AEGoFileSplitterScopeEnd].MatchString(replacedLine) {
		gps.FunctionDefinition[key].LineEnd = lineIndex
		return endState, nil
	}

	return continueState, key
}

type memberFunctionKey struct {
	Class string
	Name  string
}

func memberFunctionScope(line string, lineIndex int, gps *GoPackageScope, keyInterface interface{}, continueState, endState goFileLineState) (goFileLineState, interface{}) {
	var key *memberFunctionKey
	if keyInterface != nil {
		var ok bool
		key, ok = keyInterface.(*memberFunctionKey)
		if !ok {
			ui.OutputErrorInfo(ui.CommonError19, "keyInterface", "*memberFunctionKey")
			return endState, nil
		}
	}

	// scope begin
	if key == nil {
		var functionStructString string
		var functionKeyName string
		// var content string
		var scopeEnd string
		for _, subMatchList := range regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeMemberFunctionTemplate).FindAllStringSubmatch(line, -1) {
			if functionStructIndex, hasIndex := goSplitterMemberFunctionSubMatchNameIndexMap["MEMBER"]; hasIndex {
				functionStructString = strings.TrimSpace(subMatchList[functionStructIndex])
			}
			if functionNameIndex, hasIndex := goSplitterMemberFunctionSubMatchNameIndexMap["NAME"]; hasIndex {
				functionKeyName = strings.TrimSpace(subMatchList[functionNameIndex])
			}
			// if contentIndex, hasIndex := goSplitterMemberFunctionSubMatchNameIndexMap["CONTENT"]; hasIndex {
			// 	content = strings.TrimSpace(subMatchList[contentIndex])
			// }
			if scopeEndIndex, hasIndex := goSplitterMemberFunctionSubMatchNameIndexMap["SCOPE_END"]; hasIndex {
				scopeEnd = strings.TrimSpace(subMatchList[scopeEndIndex])
			}
		}
		// utility2.TestOutput("functionStructString = %v", functionStructString)
		// utility2.TestOutput("functionKeyName = %v", functionKeyName)
		// utility2.TestOutput("content = %v", content)
		// utility2.TestOutput("scopeEnd = %v", scopeEnd)

		// var functionClassValue string
		var functionClassValueType string
		var functionKeyClass string
		if len(functionStructString) != 0 {
			memberString := regexps.AtomicExpressionEnumRegexpMap[global.AEBracketsContent].ReplaceAllString(functionStructString, "$CONTENT")
			memberStringList := strings.Split(strings.TrimSpace(memberString), " ")
			if len(memberStringList) == 2 {
				// functionClassValue = memberStringList[0]
				functionClassValueType = memberStringList[1]
				functionKeyClass = utility.TraitStructName(functionClassValueType)
				// utility2.TestOutput("functionClassValue = %v", functionClassValue)
				// utility2.TestOutput("functionClassValueType = %v", functionClassValueType)
				// utility2.TestOutput("functionKeyClass = %v", functionKeyClass)
			} else {
				ui.OutputWarnInfo(ui.CMDAnalyzeGoFunctionDefinitionSyntaxError)
			}
		}

		if _, hasClass := gps.MemberFunctionDefinition[functionKeyClass]; !hasClass {
			gps.MemberFunctionDefinition[functionKeyClass] = make(map[string]*scope)
		}
		gps.MemberFunctionDefinition[functionKeyClass][functionKeyName] = &scope{
			LineStart: lineIndex,
			LineEnd:   lineIndex,
			ScopeType: scopeMemberFunction,
			Content:   line,
		}
		// empty struct
		if len(scopeEnd) != 0 {
			// utility2.TestOutput("%v.%v is one line function, %v", functionKeyClass, functionKeyName, scopeEnd)
			return endState, nil
		}

		return continueState, &memberFunctionKey{Class: functionKeyClass, Name: functionKeyName}
	}

	// scope content
	gps.MemberFunctionDefinition[key.Class][key.Name].Content = fmt.Sprintf("%v\n%v", gps.MemberFunctionDefinition[key.Class][key.Name].Content, line)

	// scope end
	if regexps.AtomicExpressionEnumRegexpMap[global.AEGoFileSplitterScopeEnd].MatchString(line) {
		gps.MemberFunctionDefinition[key.Class][key.Name].LineEnd = lineIndex
		// utility2.TestOutput("%v is function %v.%v scope end line", lineIndex, key.Class, key.Name)
		return endState, nil
	}

	return continueState, key
}

func typeRenameScope(line string, lineIndex int, gps *GoPackageScope, continueState, endState goFileLineState) goFileLineState {
	var typeName string
	var renameFrom string
	// var renameType string
	for _, subMatchList := range regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeTypeRenameTemplate).FindAllStringSubmatch(line, -1) {
		if typeNameIndex, hasIndex := goSplitterTypeRenameSubMatchNameIndexMap["NAME"]; hasIndex {
			typeName = strings.TrimSpace(subMatchList[typeNameIndex])
		}
		if renameFromIndex, hasIndex := goSplitterTypeRenameSubMatchNameIndexMap["FROM"]; hasIndex {
			renameFrom = strings.TrimSpace(subMatchList[renameFromIndex])
		}
		// if renameTypeIndex, hasIndex := goSplitterTypeRenameSubMatchNameIndexMap["TYPE"]; hasIndex {
		// 	renameType = strings.TrimSpace(subMatchList[renameTypeIndex])
		// }
	}
	// utility2.TestOutput("typeName = %v", typeName)
	// utility2.TestOutput("renameFrom = %v", renameFrom)
	// utility2.TestOutput("renameType = %v", renameType)
	if _, hasFrom := gps.TypeRename[renameFrom]; !hasFrom {
		gps.TypeRename[renameFrom] = make(map[string]*scope)
	}
	gps.TypeRename[renameFrom][typeName] = &scope{
		LineStart: lineIndex,
		LineEnd:   lineIndex,
		ScopeType: scopeTypeRename,
		Content:   line,
	}
	return endState
}

func multiLineConstScope(line string, lineIndex int, gps *GoPackageScope, keyInterface interface{}, continueState, endState goFileLineState) (goFileLineState, interface{}) {
	var key int
	if keyInterface != nil {
		var ok bool
		key, ok = keyInterface.(int)
		if !ok {
			ui.OutputErrorInfo(ui.CommonError19, "keyInterface", "*multiLineConstKey")
			utility2.TestOutput("key = %+v", key)
			return endState, nil
		}
	}

	// scope begin
	if keyInterface == nil {
		gps.MultiLineConst = append(gps.MultiLineConst, &scope{
			LineStart: lineIndex,
			LineEnd:   lineIndex,
			ScopeType: scopeMultiLineConst,
			Content:   line,
		})
		return continueState, len(gps.MultiLineConst) - 1
	}

	if key >= len(gps.MultiLineConst) {
		ui.OutputErrorInfo(ui.CMDCustomExecutorRunError, 7, fmt.Sprintf("multi-line scope syntax error at line: %v", lineIndex))
		return lineStateNone, nil
	}

	// scope content
	gps.MultiLineConst[key].Content = fmt.Sprintf("%v\n%v", gps.MultiLineConst[key].Content, line)

	// scope end
	if regexps.AtomicExpressionEnumRegexpMap[global.AEGoFileSplitterScopeEnd].MatchString(line) {
		gps.MultiLineConst[key].LineEnd = lineIndex
		return endState, nil
	}

	return continueState, len(gps.MultiLineConst) - 1
}

func singleLineConstScope(line string, lineIndex int, gps *GoPackageScope, continueState, endState goFileLineState) goFileLineState {
	var constName string
	// var constType string
	// var constValue string
	for _, subMatchList := range regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeSingleLineConstTemplate).FindAllStringSubmatch(line, -1) {
		if constNameIndex, hasIndex := goSplitterSingleLineConstSubMatchNameIndexMap["NAME"]; hasIndex {
			constName = strings.TrimSpace(subMatchList[constNameIndex])
		}
		// if constTypeIndex, hasIndex := goSplitterSingleLineConstSubMatchNameIndexMap["TYPE"]; hasIndex {
		// 	constType = strings.TrimSpace(subMatchList[constTypeIndex])
		// }
		// if constValueIndex, hasIndex := goSplitterSingleLineConstSubMatchNameIndexMap["VALUE"]; hasIndex {
		// 	constValue = strings.TrimSpace(subMatchList[constValueIndex])
		// }
	}
	// utility2.TestOutput("constName = %v", constName)
	// utility2.TestOutput("constType = %v", constType)
	// utility2.TestOutput("constValue = %v", constValue)
	gps.SingleLineConst[constName] = &scope{
		LineStart: lineIndex,
		LineEnd:   lineIndex,
		ScopeType: scopeSingleLineConst,
		Content:   line,
	}
	return endState
}

func getImportPackageAliasPathFromLine(line string, findRegexp *regexp.Regexp, subMatchNameIndexMap map[string]int) (string, string) {
	var importPackageAlias string
	var importPackagePath string
	for _, importSubmatchList := range findRegexp.FindAllStringSubmatch(line, -1) {
		if aliasIndex, hasAliasIndex := subMatchNameIndexMap["ALIAS"]; hasAliasIndex {
			importPackageAlias = strings.TrimSpace(importSubmatchList[aliasIndex])
		}
		if aliasIndex, hasAliasIndex := subMatchNameIndexMap["CONTENT"]; hasAliasIndex {
			importPackagePath = strings.TrimSpace(importSubmatchList[aliasIndex])
		}
		if len(importPackageAlias) == 0 {
			packagePathList := strings.Split(strings.Trim(importPackagePath, "\""), "/")
			importPackageAlias = packagePathList[len(packagePathList)-1]
		}
	}
	return importPackageAlias, importPackagePath
}

// checkGoSplitterRegexp 检查 go 文件切割器的所有模板/原子表达式
func checkGoSplitterRegexp() bool {
	ok := true
	if goFileSplitterScopePackageRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopePackageTemplate); goFileSplitterScopePackageRegexp != nil {
		goSplitterPackageSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileSplitterScopePackageRegexp.SubexpNames() {
			goSplitterPackageSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopePackageTemplate)
		ok = false
	}

	if goFileSplitterScopeMultiLineImportStartRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeMultiLineImportStartTemplate); goFileSplitterScopeMultiLineImportStartRegexp == nil {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopeMultiLineImportStartTemplate)
		ok = false
	}

	if goFileSplitterScopeMultiLineImportContentRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeMultiLineImportContentTemplate); goFileSplitterScopeMultiLineImportContentRegexp != nil {
		goSplitterMultiLineImportContentSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileSplitterScopeMultiLineImportContentRegexp.SubexpNames() {
			goSplitterMultiLineImportContentSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopeMultiLineImportContentTemplate)
		ok = false
	}

	if goFileSplitterScopeSingleLineImportRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeSingleLineImportTemplate); goFileSplitterScopeSingleLineImportRegexp != nil {
		goSplitterSingleLineImportSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileSplitterScopeSingleLineImportRegexp.SubexpNames() {
			goSplitterSingleLineImportSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopeSingleLineImportTemplate)
		ok = false
	}

	if goFileSplitterScopeEnd, has := regexps.AtomicExpressionEnumRegexpMap[global.AEGoFileSplitterScopeEnd]; !has || goFileSplitterScopeEnd == nil {
		ui.OutputErrorInfo(ui.CommonError16, global.AEGoFileSplitterScopeEnd)
		ok = false
	}

	if goFileSplitterScopePackageVariableRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopePackageVariableTemplate); goFileSplitterScopePackageVariableRegexp != nil {
		goSplitterPackageVariableSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileSplitterScopePackageVariableRegexp.SubexpNames() {
			goSplitterPackageVariableSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopePackageVariableTemplate)
		ok = false
	}

	if goFileSplitterScopeInterfaceRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeInterfaceTemplate); goFileSplitterScopeInterfaceRegexp != nil {
		goSplitterInterfaceSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileSplitterScopeInterfaceRegexp.SubexpNames() {
			goSplitterInterfaceSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopeInterfaceTemplate)
		ok = false
	}

	if goFileSplitterScopeStructRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeStructTemplate); goFileSplitterScopeStructRegexp != nil {
		goSplitterStructSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileSplitterScopeStructRegexp.SubexpNames() {
			goSplitterStructSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopeStructTemplate)
		ok = false
	}

	if goFileSplitterScopeFunctionRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeFunctionTemplate); goFileSplitterScopeFunctionRegexp != nil {
		goSplitterFunctionSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileSplitterScopeFunctionRegexp.SubexpNames() {
			goSplitterFunctionSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopeFunctionTemplate)
		ok = false
	}

	if goFileSplitterScopeMemberFunctionRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeMemberFunctionTemplate); goFileSplitterScopeMemberFunctionRegexp != nil {
		goSplitterMemberFunctionSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileSplitterScopeMemberFunctionRegexp.SubexpNames() {
			goSplitterMemberFunctionSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopeMemberFunctionTemplate)
		ok = false
	}

	if goFileSplitterScopeTypeRenameRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeTypeRenameTemplate); goFileSplitterScopeTypeRenameRegexp != nil {
		goSplitterTypeRenameSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileSplitterScopeTypeRenameRegexp.SubexpNames() {
			goSplitterTypeRenameSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopeTypeRenameTemplate)
		ok = false
	}

	if goFileSplitterScopeMultiLineConstStartRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeMultiLineConstStartTemplate); goFileSplitterScopeMultiLineConstStartRegexp == nil {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopeMultiLineConstStartTemplate)
		ok = false
	}

	if goFileSplitterScopeMultiLineConstContentRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeMultiLineConstContentTemplate); goFileSplitterScopeMultiLineConstContentRegexp != nil {
		goSplitterMultiLineConstSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileSplitterScopeMultiLineConstContentRegexp.SubexpNames() {
			goSplitterMultiLineConstSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopeMultiLineConstContentTemplate)
		ok = false
	}

	if goFileSplitterScopeSingleLineConstRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeSingleLineConstTemplate); goFileSplitterScopeSingleLineConstRegexp != nil {
		goSplitterSingleLineConstSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileSplitterScopeSingleLineConstRegexp.SubexpNames() {
			goSplitterSingleLineConstSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopeSingleLineConstTemplate)
		ok = false
	}

	if bracketsContentRegexp, hasBracketsContentRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEBracketsContent]; !hasBracketsContentRegexp || bracketsContentRegexp == nil {
		ui.OutputErrorInfo(ui.CommonError16, global.AEBracketsContent)
		ok = false
	}

	ok = true
	return ok
}

// ----------------------------------------------------------------

// custom execute 8 resources/template_example.template.go true

// Command Example: custom execute 8 resources/template_example.template.go
// Command Expression:
// - custom                                : command const content
// - execute                               : command const content
// - 8                                     : specify executor
// - resources/template_example.template.go: specify a file to split

var (
	keywordVar      string = "var"
	keywordFunc     string = "func"
	keywordReturn   string = "return"
	keywordStruct   string = "struct"
	kewordInterface string = "interface"
)

type grammaState int

const (
	nonSense                   grammaState = 1 << iota // 0000 0000 0000 0001
	variableDefinition                                 // 0000 0000 0000 0010
	functionDefinition                                 // 0000 0000 0000 0100
	structDefinition                                   // 0000 0000 0000 1000
	interfaceDefintion                                 // 0000 0000 0001 0000
	todo00100000                                       // 0000 0000 0010 0000
	todo01000000                                       // 0000 0000 0100 0000
	todo10000000                                       // 0000 0000 1000 0000
	identifier                                         // 0000 0001 0000 0000
	scopeLeftBracket                                   // 0000 0010 0000 0000
	scopeRightBracket                                  // 0000 0100 0000 0000
	scopeLeftCurlyBraces                               // 0000 1000 0000 0000
	scopeRightRightCurlyBraces                         // 0001 0000 0000 0000
	scopeLeftQuote                                     // 0010 0000 0000 0000
	scopeRightQuote                                    // 0100 0000 0000 0000
	allSense                                           // 1000 0000 0000 0000
	splitterComma                                      // 0001 0000 0000 0000 0000
)

// var[keyword] i[identifier] int[type] =[assign] 1[value] -> gramma tree
// single-line variable definition:
// - keyword: keyword
// - identifier:
// - type:
// - assign: =
// - value:
// variableDefinition
// - identifier 1 -> name
// - identifier 2 -> type

// func[keyword] f[identifier]([])
//

var lengthKeywordStateMap map[int]map[string]grammaState

type grammaNode struct {
	gramma grammaState
	next   *grammaNode
}

type grammaRule struct {
	gramma        grammaState
	boundaryState map[rune]map[grammaState]int
	nextStateMap  map[grammaState]int
}

func (rule *grammaRule) MatchState(state grammaState) bool {
	if _, hasAllSense := rule.nextStateMap[allSense]; hasAllSense {
		return true
	}
	_, hasNextState := rule.nextStateMap[state]
	return hasNextState
}

var keywordFuncGrammaRule *grammaRule

func initKeywordFuncGrammaRule() {
	keywordFuncGrammaRule = &grammaRule{
		gramma: functionDefinition,
		nextStateMap: map[grammaState]int{
			identifier: 1,
		},
	}
}

var identifierGrammaRule *grammaRule

func initIdentifierGrammaRule() {
	identifierGrammaRule = &grammaRule{
		gramma: identifier,
		nextStateMap: map[grammaState]int{
			allSense: -1,
		},
	}
}

var scopeLeftBracketGrammaRule *grammaRule

func initScopeLeftBracketGrammaRule() {
	scopeLeftBracketGrammaRule = &grammaRule{
		gramma: scopeLeftBracket,
		nextStateMap: map[grammaState]int{
			allSense: -1,
		},
	}
}

type rule struct {
}

var functionDefinitionGrammaRule interface{}

func initFunctionDefinitionGrammaRule() {
	functionDefinitionGrammaRule = &grammaRule{
		gramma: functionDefinition, // means keyword func

	}
}

func GoGrammaTree(paramList []string) {
	if len(paramList) < 1 {
		ui.OutputErrorInfo(ui.CMDCustomExecutorHasNotEnoughParam, 5)
		return
	}
	// filename := paramList[0]
	// outputString := paramList[1]
	// var output bool = false
	// if strings.ToLower(outputString) == "true" {
	// 	output = true
	// }
	initGrammaKeywordLengthMap()

	testContent1 := `
func TemplateOperatorREM(t1 int, t2 interface{}, t3 struct{ t int }) struct{ v int } {
	return struct{ v int }{v: 1}
}
`
	// testContent1 := "var i int = 1"
	utility2.TestOutput("Test Content1: |%v|", testContent1)

	generateGrammaTree(testContent1)
}

func initGrammaKeywordLengthMap() {
	keywordStateMap := map[string]grammaState{
		keywordVar:  variableDefinition,
		keywordFunc: functionDefinition,
		// keywordReturn: structDefinition,
		keywordStruct:   structDefinition,
		kewordInterface: interfaceDefintion,
	}
	lengthKeywordStateMap = make(map[int]map[string]grammaState)
	for keyword, state := range keywordStateMap {
		if _, hasLength := lengthKeywordStateMap[len(keyword)]; !hasLength {
			lengthKeywordStateMap[len(keyword)] = make(map[string]grammaState)
		}
		lengthKeywordStateMap[len(keyword)][keyword] = state
	}
}

func generateGrammaTree(content string) *grammaNode {
	var rootNode *grammaNode
	// var rootNodeRule *grammaRule
	// var currentNode *grammaNode
	// var currentNodeRule *grammaRule

	// var grammerScope int

	var currentState grammaState = nonSense
	var currentRule *grammaRule
	var builder strings.Builder
	for _, r := range content {
		// var newNode *grammaNode

		// 遇到空白字符或边界字符尝试解释已读取的字符串
		if isSpaceRune(r) || isBoundaryRune(r) {
			// 根据内容解析状态
			state := getStateByContent(builder.String())

			// 根据状态获取规则
			rule := getRuleByState(currentState)

			if currentState == nonSense && currentRule == nil {
				currentState = state
				currentRule = rule
				continue
			}

			if currentRule.MatchState(state) {

			}
		}

		_, err := builder.WriteRune(r)
		if err != nil {
			ui.OutputErrorInfo("generate gramma tree write rune %v error: %v", r, err)
			return nil
		}

		// // 根据内容获取状态
		// state = getStateByContent(len, builder.String())

		// if rootNode == nil && currentNode == nil && state == anySense {
		// 	continue
		// }

		// newNode = getGrammaNodeByState(state)
		// // 获取到了
		// if rootNode == nil && currentNode == nil {
		// 	rootNode = newNode
		// 	rootNodeRule = getNodeRule(rootNode.gramma)

		// 	currentNode = rootNode
		// 	currentNodeRule = rootNodeRule

		// 	// 清空 builder
		// 	builder.Reset()
		// } else if rootNode != nil && currentNode != nil {
		// 	if checkIfNewNodeMatchCurrentNodeRules(currentNode, newNode) {
		// 		currentNode.next = newNode
		// 		currentNode = newNode
		// 	}
		// }
	}

	return rootNode
}

func isSpaceRune(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func isBoundaryRune(r rune) bool {
	return r == '(' || r == ')' || r == '{' || r == '}'
}

func isSplitter(r rune) bool {
	return r == ','
}

func getStateByContent(content string) grammaState {
	if keywordStateMap, hasLength := lengthKeywordStateMap[len(content)]; hasLength {
		if state, hasKeyword := keywordStateMap[content]; hasKeyword {
			return state
		}
	}
	return identifier
}

func getRuleByState(state grammaState) *grammaRule {
	return nil
}

func explainPrevious(content string) *grammaNode {
	if len(content) == 0 {
		return nil
	}
	return nil
}

func getGrammaNodeByState(state grammaState) *grammaNode {
	switch state {
	case variableDefinition, functionDefinition:
		return &grammaNode{gramma: state}
	}
	return nil
}

// func generateVariableDefinitionNode() *grammaNode {
// 	return &grammaNode{gramma: variableDefinition}
// }

// func generateFunctionDefinitionNode() *grammaNode

// func checkIfNewNodeMatchCurrentNodeRules(currentNode, newNode *grammaNode) bool {
// 	currentNodeRule := getNodeRule(currentNode.gramma)
// 	// if currentNodeRule
// 	return false
// }

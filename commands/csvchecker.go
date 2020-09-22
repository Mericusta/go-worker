package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PersonalTool/code/go/go_analyzer/utility"
	"github.com/go-worker/global"
	"github.com/go-worker/ui"
)

type csvRule struct {
	description      string
	isList           bool
	listSplitter     string
	checkFileHeadMap map[string]string
}

type checkObject struct {
	file  string
	index string
	line  int
}

func (toCheckObject *checkObject) AnalyzeCSVCheckingTree() ([]string, *checkingTreeNode) {
	filePathList := make([]string, 0)
	rootNode := &checkingTreeNode{
		filename: toCheckObject.file,
		children: make(map[string]map[string]*checkingTreeNode),
	}

	toAnalyzeTreeNodeList := []*checkingTreeNode{rootNode}
	for len(toAnalyzeTreeNodeList) != 0 {
		toAnalyzeTreeNode := toAnalyzeTreeNodeList[0]
		if headRuleMap, hasToAnalyzeFile := rulesMap[toAnalyzeTreeNode.filename]; hasToAnalyzeFile {
			for head, rule := range headRuleMap {
				toAnalyzeTreeNode.children[head] = make(map[string]*checkingTreeNode)
				for filename := range rule.checkFileHeadMap {
					childNode := &checkingTreeNode{
						filename: filename,
						children: make(map[string]map[string]*checkingTreeNode),
					}
					toAnalyzeTreeNode.children[head][filename] = childNode
					toAnalyzeTreeNodeList = append(toAnalyzeTreeNodeList, childNode)
				}
			}
		}
		filePathList = append(filePathList, toAnalyzeTreeNode.filename)
		toAnalyzeTreeNodeList = toAnalyzeTreeNodeList[1:]
	}
	return filePathList, rootNode
}

type checkingTreeNode struct {
	filename string
	children map[string]map[string]*checkingTreeNode // head - target_file - target_node
}

func (toCheckingTreeNode *checkingTreeNode) CheckingRun(fileStructMap map[string]*GTCSVStruct, value string) {
	// toCheckingTreeNodeCSVStruct, hasCSVStruct := fileStructMap[toCheckingTreeNode.filename]
	// if !hasCSVStruct {
	// 	ui.OutputWarnInfo(ui.CMDCSVCheckerStructNotFoundError, toCheckingTreeNode.filename)
	// 	return
	// }

}

func check(toCheckObject *checkObject, csvStruct *CSVStruct, checkHead string, checker func(string, string) bool) bool {
	if checker == nil {
		// ERROR
		return false
	}

	fileScanner := bufio.NewScanner(csvStruct.file)
	line := 0
	checkIndex := -1
	for index, head := range csvStruct.indexHeadMap {
		if head == checkHead {
			checkIndex = index
			break
		}
	}
	if checkIndex == -1 {
		// ERROR
		return false
	}
	for ; fileScanner.Scan(); line++ {
		columnContentList := strings.Split(fileScanner.Text(), "\t")
		if len(columnContentList) <= checkIndex {
			// ERROR
			return false
		}

		// TODO:
		return checker(toCheckObject.index, columnContentList[checkIndex])
		// if columnContentList[checkIndex] == value {
		// 	return true
		// }
	}
	if err := fileScanner.Err(); err != nil {
		ui.OutputErrorInfo(ui.CommonError13, csvStruct.path, err)
	}

	return false
}

// source
var checkerMap map[string]map[string]map[string]map[string]func(string) bool

var rulesMap map[string]map[string]*csvRule

func init() {
	// checkerMap = make(map[string]map[string]map[string]map[string]func(string) bool)
	// checkerMap["charge_config.csv"] = make(map[string]map[string]map[string]func(string) bool)
	// checkerMap["charge_config.csv"]["shop_data"] = make(map[string]map[string]func(string) bool)
	// checkerMap["charge_config.csv"]["shop_data"][""]

	// checkerMap = make(map[string]map[string]map[string]func(string) bool)
	// checkerMap["charge_config.csv"]["shop_data"]["shop_data"]

	rulesMap = make(map[string]map[string]*csvRule)
	rulesMap["charge_config.csv"] = make(map[string]*csvRule)
	rulesMap["charge_config.csv"]["shop_data"] = &csvRule{
		description:  "shop_config.id,shoplist.shop_id",
		isList:       true,
		listSplitter: ";",
		checkFileHeadMap: map[string]string{
			"shop_config.csv": "shop_id",
			"shoplist.csv":    "shop_id",
		},
	}
}

type CSVChecker struct {
	*CommandStruct
}

func (command *CSVChecker) Execute() error {
	ui.OutputNoteInfo(ui.CMDCSVCheckerWelcome)

	// step 0: specify global data: region, zone
	// splitter := "\t"              // GT csv config
	// region := 1                   // checkObject
	// zone := -1                    // checkObject
	checkIndex := "ma83.3yuanbox" // checkObject

	toCheckObject := &checkObject{
		file:  "charge_config.csv",
		index: checkIndex,
	}

	// step 1: analyze checking tree
	filePathList, checkingTreeRootNode := toCheckObject.AnalyzeCSVCheckingTree()
	utility.TestOutput("filePathList = %v", filePathList)
	utility.TestOutput(ui.CommonNote2)

	checkingTreeNodeList := []*checkingTreeNode{checkingTreeRootNode}
	for len(checkingTreeNodeList) != 0 {
		checkingTreeNode := checkingTreeNodeList[0]
		utility.TestOutput("node: %v", checkingTreeNode.filename)
		if len(checkingTreeNode.children) != 0 {
			for head, targetFileMap := range checkingTreeNode.children {
				for filename, childNode := range targetFileMap {
					utility.TestOutput("node: %v, head: %v, target file: %v", checkingTreeNode.filename, head, filename)
					checkingTreeNodeList = append(checkingTreeNodeList, childNode)
				}
			}
		}
		checkingTreeNodeList = checkingTreeNodeList[1:]
	}
	utility.TestOutput(ui.CommonNote2)

	// step 2: analyze all csv file struct in checking tree
	// TODO: interface, return base class
	fileStructMap := AnalyzeGTCSVFileConcurrently(filePathList)

	for filename, fileStruct := range fileStructMap {
		fmt.Printf("filename = %v, fileStruct = %+v\n", filename, fileStruct)
		utility.TestOutput(ui.CommonNote2)
		if fileStruct.file != nil {
			fileStruct.file.Close()
		}
	}
	utility.TestOutput(ui.CommonNote2)

	// step 3: check
	// toCheckObject.check()

	return nil
}

func (command *CSVChecker) parseCommandParams() error {
	return nil
}

// load csv file

type csvStructInterface interface {
	analyzeCSVStruct() error
}

type CSVStruct struct {
	path          string
	file          *os.File
	cols          int
	startIndex    int
	splitter      string
	indexHeadMap  map[int]string
	lineOperation func(lineIndex int, lineContent string)
	fileOperation func(...interface{})
	csvRulesMap   map[int]*csvRule
}

func (csv *CSVStruct) analyzeCSVStruct() error {
	csvFile, openFileError := os.OpenFile(csv.path, os.O_RDONLY, 0644)
	if openFileError != nil {
		return openFileError
	}

	csv.file = csvFile

	fileScanner := bufio.NewScanner(csvFile)
	line := 0
	for ; fileScanner.Scan(); line++ {
		// custom logic
		if csv.lineOperation != nil {
			csv.lineOperation(line, fileScanner.Text())
		}
		if line == csv.startIndex {
			columnContentList := strings.Split(fileScanner.Text(), csv.splitter)
			csv.cols = len(columnContentList)
			break
		}
	}
	if err := fileScanner.Err(); err != nil {
		ui.OutputErrorInfo(ui.CommonError13, csv.path, err)
	}

	// custom logic for GT struct
	if csv.fileOperation != nil {
		csv.fileOperation()
	}
	return nil
}

type GTCSVStruct struct {
	CSVStruct
	indexExplainMap map[int]string
	indexTypeMap    map[int]string
}

func (gtCSVStruct *GTCSVStruct) analyzeCSVStruct() error {
	gtCSVStruct.CSVStruct.lineOperation = func(lineIndex int, lineContent string) {
		if lineIndex == 0 {
			for index, columnContent := range strings.Split(lineContent, gtCSVStruct.splitter) {
				gtCSVStruct.indexExplainMap[index] = columnContent
			}
		} else if lineIndex == 1 {
			for index, columnContent := range strings.Split(lineContent, gtCSVStruct.splitter) {
				gtCSVStruct.indexTypeMap[index] = columnContent
			}
		} else if lineIndex == 2 {
			for index, columnContent := range strings.Split(lineContent, gtCSVStruct.splitter) {
				gtCSVStruct.indexHeadMap[index] = columnContent
			}
		}
	}
	return gtCSVStruct.CSVStruct.analyzeCSVStruct()
}

// AnalyzeGTCSVFileConcurrently 并发加载并分析 GT csv 文件结构
func AnalyzeGTCSVFileConcurrently(filenameList []string) map[string]*GTCSVStruct {
	fileStructMap := make(map[string]*GTCSVStruct)
	analyzeCSVStructWaitGroup := sync.WaitGroup{}
	analyzeCSVStructWaitGroup.Add(len(filenameList))
	for _, filename := range filenameList {
		gtCSVStruct := &GTCSVStruct{
			CSVStruct: CSVStruct{
				path:         filepath.Join(global.ResourceDirectory, filename),
				startIndex:   3,
				splitter:     "\t",
				indexHeadMap: make(map[int]string),
			},
			indexExplainMap: make(map[int]string),
			indexTypeMap:    make(map[int]string),
		}
		fileStructMap[filename] = gtCSVStruct
		go func(path string) {
			analyzeCSVStructError := fileStructMap[path].analyzeCSVStruct()
			if analyzeCSVStructError != nil {
				ui.OutputErrorInfo("analyze GT csv struct error: %v", analyzeCSVStructError)
			}
			analyzeCSVStructWaitGroup.Done()
		}(filename)
	}
	analyzeCSVStructWaitGroup.Wait()
	return fileStructMap
}

// LoadCSVFileConcurrently 并发加载并分析 csv 文件
func LoadCSVFileConcurrently(filePathList []string) map[string]*os.File {
	pathFileMap := make(map[string]*os.File)
	pathFileMapMutex := new(sync.Mutex)
	pathErrorMap := make(map[string]error)
	pathErrorMapMutex := new(sync.Mutex)
	loadingWaitGroup := sync.WaitGroup{}
	loadingWaitGroup.Add(len(filePathList))
	for _, filePath := range filePathList {
		go func(path string) {
			csvFile, openFileError := os.OpenFile(path, os.O_RDONLY, 0644)
			if openFileError != nil {
				pathErrorMapMutex.Lock()
				pathErrorMap[path] = openFileError
				pathErrorMapMutex.Unlock()
			} else {
				pathFileMapMutex.Lock()
				pathFileMap[path] = csvFile
				pathFileMapMutex.Unlock()
			}
			loadingWaitGroup.Done()
		}(filepath.Join(global.ResourceDirectory, filePath))
	}
	loadingWaitGroup.Wait()
	for path, loadingError := range pathErrorMap {
		ui.OutputErrorInfo(ui.CommonError5, path, loadingError)
	}
	return pathFileMap
}

// object

// checking rule
//
// example: charge_config.csv
// #地区ID	#区服ID用于内部区分	#充值ID	商店相关数据
// string	string	string	string
// region	zone_id	charge_id	shop_data
// 2,3	-1	ma83.day.18yuanbox1	145,200144
//
// charge_config.shop_data => shop_config.id,shoplist.shop_id[;shop_config.id,shoplist.shop_id]
//
// rule1: charge_config.shop_data -> shop_config, if exist
// rule2: charge_config.shop_data -> shoplist, if exist
// rule3: shoplist.group -> shop_config.group_config, if exist

// shop_config.

func analyzeRule() {
	// originRule := "shop_config.id,shoplist.shop_id"
	// TK_CSV_FILE.TK_CSV_FIELD
}

// 1 load charge_config.csv
// 2 load rule: shop_config.id,shoplist.shop_id[;shop_config.id,shoplist.shop_id] -> load all rules about charge_config.csv
// 3 load shop_config.csv, shoplist.csv and analyze csv struct -> load all csv struct in rules
// 4 load rule: ...

// func analyzeChargeConfigShopDataValue(columnContent string) bool {
// 	// columnContent
// 	columnContent = "145,200144"

// 	for _, eachShopData := range strings.Split(columnContent, ";") {
// 		shopDataList := strings.Split(eachShopData, ",")
// 		if len(shopDataList) != 2 {
// 			ui.OutputWarnInfo(ui.CMDCSVCheckerFieldValueFormatError, "charge_config", "shop_data")
// 			continue
// 		}
// 		shopConfigID := shopDataList[0]
// 		shoplistShopID := shopDataList[1]

// 	}

// 	return true
// }

// func chargeConfigShopDataRule1(shopIDList []string) map[string]bool {
// 	checkResult := make(map[string]bool)
// 	for _, shopID := range shopIDList {
// 		checkResult[shopID] = false

// 		// get shop_config file object
// 		//
// 	}
// 	return checkResult
// }

// ------------

// 0 specify global data: region, zone
// 1 load csv file charge_config
// 2 check all columns rule: charge_config.shop_data
// 2-1 charge_config.shop_data -> shop_config and shoplist
// 3 load csv file shop_config and shoplist
// 4

// ---

// map[source file][source head]{
//     decription -> non-programmer define
//     isList -> non-programmer define
//     listSplitter -> non-programmer define

//     map[target file][target head]{ -> program analyze
//         object1{
//             file
//             head
//             checker -> programmer define
//         },
//         object2{
//             file
//             head
//             checker -> programmer define
//         },
//         ...
//     },
// }

// ---

// type 1: know value, and the value point to some csv structs, check if value exists in a those config

// type 2: know value, and the value referenced by another config's value, check if

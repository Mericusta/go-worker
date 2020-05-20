package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/go-worker/ui"
	"github.com/go-worker/utility"
)

const workerOperationIndex = 1
const workerProjectIndex = 2
const workerProjectParamsIndex = 3

var workerOperationList []string
var workerProjectExecutorMap map[string]func([]string) error

var goTypeStringList []string

type workerParams struct {
	Operation     string
	Project       string
	ProjectParams []string
}

func init() {
	workerOperationList = []string{
		"bind",
		"run",
		"create",
		"help",
		"quit",
	}
	workerProjectExecutorMap = map[string]func([]string) error{
		"0001": workerProject0001Executor,
	}
}

func parseWorkerParams(inputString string) (*workerParams, error) {
	var params *workerParams
	inputStringList := strings.Split(inputString, " ")
	if len(inputStringList) > workerOperationIndex {
		params = new(workerParams)
		params.Operation = inputStringList[workerOperationIndex]
		params.Project = inputStringList[workerProjectIndex]
		params.ProjectParams = append(params.ProjectParams, inputStringList[workerProjectParamsIndex:]...)
	} else {
		return nil, fmt.Errorf(ui.FSMUnknownCommand)
	}
	return params, nil
}

// WorkerExecutor worker 命令执行
func WorkerExecutor(inputString string) error {
	utility.TestOutput("WorkerExecutor, inputString = %v", inputString)
	params, parseParamsError := parseWorkerParams(inputString)
	if parseParamsError != nil {
		return parseParamsError
	}
	if params != nil {
		utility.TestOutput("WorkerExecutor, operation = %v", params.Operation)
		utility.TestOutput("WorkerExecutor, project = %v", params.Project)
		switch params.Operation {
		case "bind":
			{

			}
		case "run":
			{
				projectExecutor, hasProjectExecutor := workerProjectExecutorMap[params.Project]
				if hasProjectExecutor {
					projectExecuteError := projectExecutor(params.ProjectParams)
					if projectExecuteError != nil {
						return projectExecuteError
					}
				} else {
					return fmt.Errorf("worker run project %v but no executor found", params.Project)
				}
			}
		case "help":
			{
				fmt.Print("TEST: worker help")
			}
		case "quit":
			{
				fmt.Print("TEST: worker quit")
			}
		}
	}
	return nil
}

type goStructMemberInfo struct {
	structName    string
	structType    string
	structCsv     string
	structComment string
}

// project 0001 is to maintain csv struct in JoyNova Server
// it provides append/remove/update for go struct member
const project0001StructNameIndex = 0
const project0001OperationIndex = 1
const project0001OperationParamsIndex = 2

func workerProject0001Executor(projectParams []string) error {
	if project0001OperationParamsIndex > len(projectParams) {
		return fmt.Errorf(ui.FSMUnknownCommand)
	}

	structName := projectParams[project0001StructNameIndex]
	operation := projectParams[project0001OperationIndex]
	var operationParamsList []string
	if len(projectParams) > project0001OperationParamsIndex {
		operationParamsList = projectParams[project0001OperationParamsIndex:]
	}

	goFilePath, goFileContent, openGoFileError := openGoFile()
	if openGoFileError != nil {
		return openGoFileError
	}

	goStructExpression := fmt.Sprintf(`(?s)type %v struct \{.*?\}`, structName)
	goStructRegexp := regexp.MustCompile(goStructExpression)

	goStructString := goStructRegexp.FindString(string(goFileContent))
	goStructStringInterval := goStructRegexp.FindStringIndex(goFileContent)
	utility.TestOutput("len(goStructString) = %v, len(goStructStringInterval) = %v", len(goStructString), len(goStructStringInterval))

	if len(goStructString) == 0 || len(goStructStringInterval) == 0 {
		return fmt.Errorf("not find struct %v content in file", structName)
	}

	utility.TestOutput("goStructString = %v", goStructString)
	utility.TestOutput("goStructStringInterval = %v", goStructStringInterval)
	utility.TestOutput("len = %v = %v", len(goStructString), goStructStringInterval[1]-goStructStringInterval[0])

	goStructLineList := strings.Split(goStructString, "\n")

	utility.TestOutput("operation = %v", operation)
	var newGoStructString string
	switch operation {
	case "append":
		// 增加
		var appendGoStructMemberInfoError error
		newGoStructString, appendGoStructMemberInfoError = appendGoStructMemberInfo(goStructLineList)
		if appendGoStructMemberInfoError != nil {
			return appendGoStructMemberInfoError
		}
	case "remove":
		// 删除
		var removeGoStructMemberInfoError error
		newGoStructString, removeGoStructMemberInfoError = removeGoStructMemberInfo(goStructLineList, operationParamsList)
		if removeGoStructMemberInfoError != nil {
			return removeGoStructMemberInfoError
		}
	case "update":
		// 修改
		var updateGoStructMemberInfoError error
		newGoStructString, updateGoStructMemberInfoError = updateGoStructMemberInfo(goStructLineList, operationParamsList)
		if updateGoStructMemberInfoError != nil {
			return updateGoStructMemberInfoError
		}
	case "list":
		// 查看
	case "help":
		// 帮助
	default:
		return fmt.Errorf(ui.FSMUnknownCommand)
	}

	// 写入
	newFileContent := goFileContent[0:goStructStringInterval[0]] + newGoStructString + goFileContent[goStructStringInterval[1]:]
	reconverFileContent := strings.ReplaceAll(newFileContent, "\\n", "\n")
	ioutil.WriteFile(goFilePath, []byte(reconverFileContent), 0644)

	return nil
}

func openGoFile() (string, string, error) {
	fmt.Printf(ui.FSMUnknownCommand)

	input := bufio.NewScanner(os.Stdin)
	if hasInput := input.Scan(); !hasInput {
		return "", "", fmt.Errorf("file name is required")
	}

	projectDir, getProjectDirError := os.Getwd()
	if getProjectDirError != nil {
		return "", "", getProjectDirError
	}
	utility.TestOutput("project dir = %v", projectDir)
	file, openFileError := os.OpenFile(fmt.Sprintf("%v/resources/%v", projectDir, input.Text()), os.O_RDWR, 0644)
	defer file.Close()
	if openFileError != nil {
		return "", "", openFileError
	}
	fileContentByteList, readAllError := ioutil.ReadAll(file)
	if readAllError != nil {
		return "", "", readAllError
	}
	return fmt.Sprintf("%v/resources/struct.go", projectDir), string(fileContentByteList), nil
}

func appendGoStructMemberInfo(goStructLineList []string) (string, error) {
	newGoStructMemberInfo := new(goStructMemberInfo)
	fmt.Println(ui.FSMUnknownCommand)
	input := bufio.NewScanner(os.Stdin)
	if hasInput := input.Scan(); hasInput {
		inputStringList := strings.Split(input.Text(), "|")
		if len(inputStringList) < 3 {
			return "", fmt.Errorf("not enough value")
		}
		utility.TestOutput("inputStringList = %v", inputStringList)
		newGoStructMemberInfo.structName = inputStringList[0]
		newGoStructMemberInfo.structType = inputStringList[1]
		newGoStructMemberInfo.structCsv = inputStringList[2]
		if len(inputStringList) > 4 {
			newGoStructMemberInfo.structComment = strings.Join(inputStringList[3:], ",")
		} else {
			newGoStructMemberInfo.structComment = inputStringList[3]
		}
	} else {
		return "", fmt.Errorf("input nothing")
	}

	newGoStructMemberInfoString := fmt.Sprintf("%v %v `csv:\"%v\"` //%v", newGoStructMemberInfo.structName, newGoStructMemberInfo.structType, newGoStructMemberInfo.structCsv, newGoStructMemberInfo.structComment)
	firstSlice := goStructLineList[:len(goStructLineList)-1]
	utility.TestOutput("len(firstSlice) = %v", len(firstSlice))
	lastSlice := append([]string{}, goStructLineList[len(goStructLineList)-1:]...)
	utility.TestOutput("lastSlice = %v", lastSlice)
	newSlice := append(firstSlice, newGoStructMemberInfoString)
	utility.TestOutput("first slice append newGoStructMemberInfoString, len(newSlice) = %v", len(newSlice))
	utility.TestOutput("lastSlice = %v", lastSlice)
	newSlice = append(newSlice, lastSlice...)
	utility.TestOutput("newSlice append lastSlice..., len(newSlice) = %v", len(newSlice))
	return strings.Join(newSlice, "\n"), nil
}

func removeGoStructMemberInfo(goStructLineList []string, memberToRemove []string) (string, error) {
	if len(memberToRemove) == 0 {
		return "", fmt.Errorf(ui.FSMUnknownCommand)
	}
	for _, operationParam := range memberToRemove {
		for lineIndex, goStructLine := range goStructLineList {
			keywordIndex := strings.Index(strings.TrimSpace(goStructLine), operationParam)
			if keywordIndex != -1 && keywordIndex == 0 {
				goStructLineList = append(goStructLineList[:lineIndex], goStructLineList[lineIndex+1:]...)
				break
			}
		}
	}
	return strings.Join(goStructLineList, "\n"), nil
}

const project0001UpdateMemberIndex = 0
const project0001UpdateMemberKey = 1
const project0001UpdateMemberValue = 2

func updateGoStructMemberInfo(goStructLineList []string, paramsList []string) (string, error) {
	goStructMemberInfoMap := analysisGoStructData(goStructLineList)
	if len(paramsList) > project0001UpdateMemberValue {
		found := false
		for lineIndex, goStructMemberInfo := range goStructMemberInfoMap {
			if goStructMemberInfo.structName == paramsList[project0001UpdateMemberIndex] {
				switch paramsList[project0001UpdateMemberKey] {
				case "name":
					goStructMemberInfo.structName = paramsList[project0001UpdateMemberValue]
				case "type":
					goStructMemberInfo.structType = paramsList[project0001UpdateMemberValue]
				case "csv":
					goStructMemberInfo.structCsv = paramsList[project0001UpdateMemberValue]
				case "comment":
					goStructMemberInfo.structComment = paramsList[project0001UpdateMemberValue]
				default:
					return "", fmt.Errorf(ui.FSMUnknownCommand)
				}
				found = true
				newGoStructMemberInfoString := fmt.Sprintf("%v %v `csv:\"%v\"` //%v", goStructMemberInfo.structName, goStructMemberInfo.structType, goStructMemberInfo.structCsv, goStructMemberInfo.structComment)
				goStructLineList[lineIndex] = newGoStructMemberInfoString
				break
			}
		}
		if !found {
			return "", fmt.Errorf(ui.FSMUnknownCommand)
		}
	} else {
		return "", fmt.Errorf(ui.FSMUnknownCommand)
	}
	return strings.Join(goStructLineList, "\n"), nil
}

func analysisGoStructData(goStructLineList []string) map[int]*goStructMemberInfo {
	// 查看
	goStructMemberInfoMap := make(map[int]*goStructMemberInfo, 0)
	for lineIndex, goStructLine := range goStructLineList {
		if lineIndex == 0 || lineIndex == len(goStructLineList)-1 {
			continue
		}
		goStructLineStringList := strings.Fields(goStructLine)
		utility.TestOutput("goStructLineStringList = %v", goStructLineStringList)
		goStructCsvField := strings.Split(goStructLineStringList[2], "\"")[1]
		goStructCommentIndex := strings.Index(goStructLine, "//")
		goStructCommentString := ""
		if goStructCommentIndex != -1 {
			goStructCommentString = strings.TrimSpace(goStructLine[goStructCommentIndex+2:])
		}
		goStructMemberInfoMap[lineIndex] = &goStructMemberInfo{
			structName:    goStructLineStringList[0],
			structType:    goStructLineStringList[1],
			structCsv:     goStructCsvField,
			structComment: goStructCommentString,
		}
	}

	for lineIndex, goStructMemberInfo := range goStructMemberInfoMap {
		utility.TestOutput("lineIndex = %v, member name = %v, member type = %v, member csv = %v, comment = %v", lineIndex, goStructMemberInfo.structName, goStructMemberInfo.structType, goStructMemberInfo.structCsv, goStructMemberInfo.structComment)
	}

	return goStructMemberInfoMap
}

// TODO: func trait
func traitFormat(goStructLineList []string) {

}

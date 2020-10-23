package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-worker/global"
	"github.com/go-worker/regexps"
	"github.com/go-worker/ui"
	"github.com/go-worker/utility"
	"github.com/go-worker/utility2"
	"github.com/go-worker/utility3"
)

type Analyze struct {
	*CommandStruct
	Params *analyzeParam
}

func (command *Analyze) Execute() error {
	// 解析指令的选项和参数
	parseCommandParamsError := command.parseCommandParams()
	if parseCommandParamsError != nil {
		return parseCommandParamsError
	}

	utility2.TestOutput("command.Params = %v", command.Params)

	// // 获取全局配置
	// projectPath := config.GetCurrentProjectPath()
	// fileType := config.WorkerConfig.ProjectSyntax
	// if fileType == "" {
	// 	ui.OutputWarnInfo(ui.CommonError14)
	// }

	// if command.Params.parentValue != "" {
	// 	projectPath = filepath.Join(projectPath, command.Params.parentValue)
	// }

	// if command.Params.sourceValue == "" {
	// 	return fmt.Errorf(ui.CommonError1)
	// }

	// toAnalyzePath := projectPath
	// if command.Params.sourceValue != "." {
	// 	toAnalyzePath = filepath.Join(toAnalyzePath, command.Params.sourceValue)
	// }

	// toWriteFilePath := ""
	// if command.Params.outputValue != "" {
	// 	toWriteFilePath = filepath.Join(toAnalyzePath, command.Params.outputValue)
	// }

	// toAnalyzeWriteFilePathMap := make(map[string]string, 0)
	// switch command.Params.sourceType {
	// case "file":
	// 	toAnalyzeFilePath := filepath.Join(projectPath, fmt.Sprintf("%v.%v", command.Params.sourceValue, fileType))
	// 	if !utility.IsExist(toAnalyzeFilePath) {
	// 		return fmt.Errorf(ui.CMDAnalyzeFileOrDirectoryNotExist, toAnalyzeFilePath)
	// 	}
	// 	if toWriteFilePath != "" {
	// 		toAnalyzeWriteFilePathMap[toAnalyzeFilePath] = toWriteFilePath
	// 	} else {
	// 		toAnalyzeWriteFilePathMap[toAnalyzeFilePath] = fmt.Sprintf("%v.%v", toAnalyzeFilePath, global.SyntaxMarkdown)
	// 	}
	// case "directory":
	// 	directoryStat, getStatError := os.Stat(toAnalyzePath)
	// 	if getStatError != nil {
	// 		return fmt.Errorf(ui.CommonError7, toAnalyzePath, getStatError)
	// 	}
	// 	if !directoryStat.IsDir() {
	// 		ui.OutputWarnInfo(ui.CommonError8, toAnalyzePath)
	// 		return nil
	// 	}
	// 	for _, toAnalyzeFilePath := range utility.TraverseDirectorySpecificFile(toAnalyzePath, fileType) {
	// 		if !utility.IsExist(toAnalyzeFilePath) {
	// 			ui.OutputWarnInfo(ui.CMDAnalyzeFileOrDirectoryNotExist, toAnalyzeFilePath)
	// 			continue
	// 		}
	// 		if toWriteFilePath != "" {
	// 			toAnalyzeWriteFilePathMap[toAnalyzeFilePath] = toWriteFilePath
	// 		} else {
	// 			toAnalyzeWriteFilePathMap[toAnalyzeFilePath] = fmt.Sprintf("%v.%v", toAnalyzeFilePath, global.SyntaxMarkdown)
	// 		}
	// 	}
	// default:
	// 	break
	// }

	// var analyzeFunction func(string, map[string]string) (interface{}, error)
	// switch fileType {
	// case global.SyntaxGo:
	// 	analyzeFunction = analyzeGo
	// case global.SyntaxCpp:
	// 	analyzeFunction = nil
	// }

	// _, analyzeError := analyzeFunction(toAnalyzePath, toAnalyzeWriteFilePathMap)
	// if analyzeError != nil {
	// 	return analyzeError
	// }

	return nil
}

type analyzeParam struct {
	// sourceType  string
	// sourceValue string
	// parentValue string
	sourceValue string
	outputValue string
}

// TODO:
func (command *Analyze) parseCommandParams() error {
	// optionValueString := ""
	// if optionValueRegexp, hasOptionValueRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEAnalyzeOptionValue]; hasOptionValueRegexp {
	// 	optionValueString = optionValueRegexp.FindString(command.CommandStruct.InputString)
	// } else {
	// 	ui.OutputWarnInfo(ui.CommonError15, "analyze", "file|directory|package")
	// }
	// if optionValueString == "" {
	// 	return fmt.Errorf(ui.CommonError1)
	// }
	// optionValueList := strings.Split(optionValueString, " ")
	// command.Params = &analyzeParam{
	// 	sourceType:  optionValueList[0],
	// 	sourceValue: optionValueList[1],
	// }
	// parentValue := ""
	// if parentValueRegexp := regexps.GetRegexpByTemplateEnum(global.OptionParentValueTemplate); parentValueRegexp != nil {
	// 	parentValue = parentValueRegexp.FindString(command.CommandStruct.InputString)
	// } else {
	// 	ui.OutputWarnInfo(ui.CommonError15, "analyze", "parent")
	// }
	// if parentValue != "" {
	// 	parentValueList := strings.Split(parentValue, " ")
	// 	command.Params.parentValue = parentValueList[1]
	// }
	command.Params = &analyzeParam{}
	analyzePath := ""
	if pathRegexp, hasPathRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEPath]; hasPathRegexp {
		utility2.TestOutput("command.CommandStruct.InputString = %v", command.CommandStruct.InputString)
		analyzePath = pathRegexp.ReplaceAllString(command.CommandStruct.InputString, "$PATH")
		command.Params.sourceValue = analyzePath
	} else {
		ui.OutputWarnInfo(ui.CommonError16, global.AEPath)
	}
	outputValue := ""
	if outputValueRegexp := regexps.GetRegexpByTemplateEnum(global.OptionOutputValueTemplate); outputValueRegexp != nil {
		outputValue = outputValueRegexp.FindString(command.CommandStruct.InputString)
	} else {
		ui.OutputWarnInfo(ui.CommonError15, "analyze", "output")
	}
	utility2.TestOutput("outputValue = %v", outputValue)
	if outputValue != "" {
		outputValueList := strings.Split(outputValue, " ")
		command.Params.outputValue = outputValueList[0]
	}
	return nil
}

// 分析 GO 项目

// GoAnalysis go 项目分析结果
type GoAnalysis struct {
	// FileNoAnalysisMap    map[int]*goFileAnalysis
	// PackageNoAnalysisMap map[int]*GoPackageAnalysis
	// MainPackageAnalysis *GoPackageAnalysis
	// PackageFunctionMap   map[string]string
	PackageAnalysisMap map[string]*GoPackageAnalysis
}

// GoPackageAnalysis go 包分析结果
type GoPackageAnalysis struct {
	// No                          int
	// FileNoList                  []int
	// ImportPackageAnalysisNoList []int
	FunctionAnalysisMap map[string]*GoFunctionAnalysis
}

// goFileAnalysis go 文件分析结果
type goFileAnalysis struct {
	// No             int
	FilePath       string
	PackageName    string
	PackagePath    string
	ImportAliasMap map[string]string
	FunctionMap    map[string]*GoFunctionAnalysis
	functionList   []string
}

// GoFunctionAnalysis go 函数分析结果
type GoFunctionAnalysis struct {
	Class               string // 所属类
	ClassValue          string // 类对象命名
	ClassValueType      string // 类对象类型
	Name                string
	ParamsMap           map[string]*GoFunctionVariable
	ReturnMap           map[string]*GoFunctionVariable
	InnerPackageCallMap map[string]map[int]*GoFunctionCallAnalysis            // call function from inner package
	OuterPackageCallMap map[string]map[string]map[int]*GoFunctionCallAnalysis // call function from outer package
	MemberCallMap       map[string]map[string]map[int]*GoFunctionCallAnalysis // call member function
	CallerMap           map[string]map[string]map[int]*GoFunctionCallAnalysis // package path : function : call index : call content
	BodyContent         []byte
	VariableMap         map[string]*GoFunctionVariable
}

// GoFunctionCallAnalysis go 函数调用分析结果
type GoFunctionCallAnalysis struct {
	From      string
	Content   string
	ParamList []string
}

type GoFunctionVariable struct {
	Index int
	Name  string
	Type  string
}

func analyzeGo(toAnalyzePathList []string) (interface{}, error) {
	const mainPackageNo = 0
	// packageNo := mainPackageNo
	// fileNo := 0
	goAnalysis := &GoAnalysis{
		// FileNoAnalysisMap:    make(map[int]*goFileAnalysis),
		// PackageNoAnalysisMap: make(map[int]*GoPackageAnalysis),
		// PackageFunctionMap:   make(map[string]string),
		PackageAnalysisMap: make(map[string]*GoPackageAnalysis),
	}
	// packagePathAnalysisMap := make(map[string]*GoPackageAnalysis)

	for _, toAnalyzeFilePath := range toAnalyzePathList {
		// var toWriteFile *os.File
		// defer func() {
		// 	if toWriteFile != nil {
		// 		toWriteFile.Close()
		// 	}
		// }()
		// if toWriteFilePath := toAnalyzeWriteFilePathMap[toAnalyzeFilePath]; len(toWriteFilePath) != 0 {
		// 	if utility.IsExist(toWriteFilePath) {
		// 		var openFileError error
		// 		toWriteFile, openFileError = os.OpenFile(toWriteFilePath, os.O_RDWR|os.O_TRUNC, 0644)
		// 		if openFileError != nil {
		// 			return goAnalysis, openFileError
		// 		}
		// 	} else {
		// 		var createFileError error
		// 		toWriteFile, createFileError = utility.CreateFile(toWriteFilePath)
		// 		if createFileError != nil {
		// 			return goAnalysis, createFileError
		// 		}
		// 	}
		// 	if toWriteFile == nil {
		// 		return goAnalysis, fmt.Errorf("analyze to write file is nil")
		// 	}
		// }

		toAnalyzeFile, inputError := os.Open(toAnalyzeFilePath)
		defer toAnalyzeFile.Close()
		if inputError != nil || toAnalyzeFile == nil {
			return goAnalysis, fmt.Errorf(ui.CommonError5, toAnalyzeFilePath, inputError.Error())
		}

		fileAnalysis, analyzeError := analyzeGoFile(toAnalyzeFile)
		if analyzeError != nil {
			ui.OutputWarnInfo(ui.CMDAnalyzeOccursError, analyzeError)
		}
		// fileAnalysis.No = fileNo

		if _, hasPackageAnalysis := goAnalysis.PackageAnalysisMap[fileAnalysis.PackagePath]; !hasPackageAnalysis {
			// goPackageAnalysis = &GoPackageAnalysis{
			// No: func() int {
			// 	if fileAnalysis.PackageName == "main" {
			// 		return mainPackageNo
			// 	}
			// 	packageNo++
			// 	return packageNo
			// }(),
			// FileNoList:                  make([]int, 0),
			// ImportPackageAnalysisNoList: make([]int, 0),
			// FunctionAnalysisMap: make(map[string]*GoFunctionAnalysis),
			// }
			// packagePathAnalysisMap[fileAnalysis.PackagePath] = goPackageAnalysis
			// goAnalysis.PackageNoAnalysisMap[goPackageAnalysis.No] = goPackageAnalysis
			goAnalysis.PackageAnalysisMap[fileAnalysis.PackagePath] = &GoPackageAnalysis{
				FunctionAnalysisMap: make(map[string]*GoFunctionAnalysis),
			}
		}
		// packagePathAnalysisMap[fileAnalysis.PackagePath].FileNoList = append(packagePathAnalysisMap[fileAnalysis.PackagePath].FileNoList, fileAnalysis.No)
		// goAnalysis.FileNoAnalysisMap[fileAnalysis.No] = fileAnalysis
		// fileNo++

		for functionName, functionAnalysis := range fileAnalysis.FunctionMap {
			if _, hasFunction := goAnalysis.PackageAnalysisMap[fileAnalysis.PackagePath].FunctionAnalysisMap[functionName]; hasFunction {
				ui.OutputWarnInfo(ui.CMDAnalyzeGoDuplicateFunctionInPackage, functionName, fileAnalysis.PackagePath)
			}
			goAnalysis.PackageAnalysisMap[fileAnalysis.PackagePath].FunctionAnalysisMap[functionName] = functionAnalysis
		}
	}

	utility2.TestOutput(ui.CommonNote2)
	utility2.TestOutput("analyze go function caller")
	for packagePath, packageAnalysis := range goAnalysis.PackageAnalysisMap {
		utility2.TestOutput("packagePath = %v", packagePath)
		for functionName, functionAnalysis := range packageAnalysis.FunctionAnalysisMap {
			// Outer Package Call
			for callFromPackagePath, callFunctionMap := range functionAnalysis.OuterPackageCallMap {
				callFromPackageAnalysis, hasCallFromPackageAnalysis := goAnalysis.PackageAnalysisMap[callFromPackagePath]
				if !hasCallFromPackageAnalysis {
					ui.OutputWarnInfo(ui.CMDAnalyzeGoPackageAnalysisNotExist, callFromPackagePath)
					continue
				}
				for callFunction, functionCallAnalysisMap := range callFunctionMap {
					callFunctionAnalysis, hasCallFunctionAnalysis := callFromPackageAnalysis.FunctionAnalysisMap[callFunction]
					if !hasCallFunctionAnalysis {
						ui.OutputWarnInfo(ui.CMDAnalyzeGoFunctionAnalysisNotExist, callFunction)
						continue
					}
					for callIndex, functionCallAnalysis := range functionCallAnalysisMap {
						callFunctionAnalysis.CallerMap[packagePath][functionName][callIndex] = functionCallAnalysis
						utility2.TestOutput("package %v function %v call package %v function %v by params %+v", packagePath, functionName, callFromPackagePath, callFunctionAnalysis.Name, functionCallAnalysis)
					}
				}
			}

			// Member Call
			for callFromPackagePath, callFunctionMap := range functionAnalysis.MemberCallMap {
				callFromPackageAnalysis, hasCallFromPackageAnalysis := goAnalysis.PackageAnalysisMap[callFromPackagePath]
				if !hasCallFromPackageAnalysis {
					ui.OutputWarnInfo(ui.CMDAnalyzeGoPackageAnalysisNotExist, callFromPackagePath)
					continue
				}
				for callFunction, functionCallAnalysisMap := range callFunctionMap {
					callFunctionAnalysis, hasCallFunctionAnalysis := callFromPackageAnalysis.FunctionAnalysisMap[callFunction]
					if !hasCallFunctionAnalysis {
						ui.OutputWarnInfo(ui.CMDAnalyzeGoFunctionAnalysisNotExist, callFunction)
						continue
					}
					for callIndex, functionCallAnalysis := range functionCallAnalysisMap {
						callFunctionAnalysis.CallerMap[packagePath][functionName][callIndex] = functionCallAnalysis
						utility2.TestOutput("package %v function %v call package %v function %v by params %+v", packagePath, functionName, callFromPackagePath, callFunctionAnalysis.Name, functionCallAnalysis)
					}
				}
			}
		}
	}
	utility2.TestOutput(ui.CommonNote2)

	// if mainPackageAnalysis, hasMainPackage := goAnalysis.PackageNoAnalysisMap[mainPackageNo]; hasMainPackage {
	// 	goAnalysis.MainPackageAnalysis = mainPackageAnalysis
	// }

	// for _, packageAnalysis := range packagePathAnalysisMap {
	// 	for _, fileNo := range packageAnalysis.FileNoList {
	// 		if fileAnalysis, hasFileAnalysis := goAnalysis.FileNoAnalysisMap[fileNo]; hasFileAnalysis {
	// 			for _, importPackagePath := range fileAnalysis.ImportAliasMap {
	// 				if importPackageAnalysis, hasImportPackageAnalysis := packagePathAnalysisMap[importPackagePath]; hasImportPackageAnalysis {
	// 					found := false
	// 					for _, importPackageNo := range packageAnalysis.ImportPackageAnalysisNoList {
	// 						if importPackageNo == importPackageAnalysis.No {
	// 							found = true
	// 							break
	// 						}
	// 					}
	// 					if found {
	// 						continue
	// 					}
	// 					packageAnalysis.ImportPackageAnalysisNoList = append(packageAnalysis.ImportPackageAnalysisNoList, importPackageAnalysis.No)
	// 				} else {
	// 					ui.OutputWarnInfo(ui.CMDAnalyzeGoPackageAnalysisNotExist, importPackagePath)
	// 				}
	// 			}
	// 		} else {
	// 			ui.OutputWarnInfo(ui.CMDAnalyzeGoFileAnalysisNotExist, fileNo)
	// 		}
	// 	}
	// }

	// for packagePath, packageAnalysis := range packagePathAnalysisMap {
	// 	utility2.TestOutput("No: %v, Package Path: %v, Import: %v", packageAnalysis.No, packagePath, packageAnalysis.ImportPackageAnalysisNoList)
	// 	goAnalysis.PackageNoAnalysisMap[packageAnalysis.No] = packageAnalysis
	// }

	// nTreeNodeChildrenMap := makeUpNTreeNodeChildrenMapByGoPackage(goAnalysis.PackageNoAnalysisMap)
	// utility2.TestOutput("nTreeNodeChildrenMap = %+v", nTreeNodeChildrenMap)

	// if len(nTreeNodeChildrenMap) != 0 {
	// 	mergedNTree := utility.NTreeHierarchicalMergeAlgorithmImproved(nTreeNodeChildrenMap)
	// 	for level, node := range mergedNTree {
	// 		utility2.TestOutput("level = %v, node = %v", level, node)
	// 	}

	// 	// 输出包级有向图
	// 	abs, getAbsError := filepath.Abs(toAnalyzePath)
	// 	if getAbsError != nil {
	// 		return getAbsError
	// 	}
	// 	utility2.TestOutput("%v", filepath.Base(abs))
	// 	var toWriteGoAnalysisFile *os.File
	// 	toWriteGoAnalysisFilePath := fmt.Sprintf("%v.%v", filepath.Base(abs), global.SyntaxMarkdown)
	// 	if utility.IsExist(toWriteGoAnalysisFilePath) {
	// 		var openFileError error
	// 		toWriteGoAnalysisFile, openFileError = os.OpenFile(toWriteGoAnalysisFilePath, os.O_RDWR|os.O_TRUNC, 0644)
	// 		if openFileError != nil {
	// 			return openFileError
	// 		}
	// 	} else {
	// 		var createFileError error
	// 		toWriteGoAnalysisFile, createFileError = utility.CreateFile(toWriteGoAnalysisFilePath)
	// 		if createFileError != nil {
	// 			return createFileError
	// 		}
	// 	}

	// 	goPackageLevelDirectedGraph := outputGoPackageLevelDirectedGraph(mergedNTree)
	// 	_, writeError := toWriteGoAnalysisFile.WriteString(goPackageLevelDirectedGraph)
	// 	if writeError != nil {
	// 		return writeError
	// 	}
	// }

	return goAnalysis, nil
}

func analyzeGoFile(toAnalyzeFile *os.File) (*goFileAnalysis, error) {
	goFileAnalysis := &goFileAnalysis{
		ImportAliasMap: make(map[string]string),
		FunctionMap:    make(map[string]*GoFunctionAnalysis),
		functionList:   make([]string, 0),
	}

	toAnalyzeContent, readToAnalyzeContentError := ioutil.ReadAll(toAnalyzeFile)
	if readToAnalyzeContentError != nil {
		return nil, readToAnalyzeContentError
	}

	toAnalyzeContent = removeGoFileCommentLine(toAnalyzeContent)

	// 文件路径
	filePath, getFileAbsPathError := filepath.Abs(toAnalyzeFile.Name())
	if getFileAbsPathError != nil {
		return nil, getFileAbsPathError
	}
	goFileAnalysis.FilePath = strings.Replace(filePath, "\\", "/", -1)

	// 解析包名
	analyzeGoKeywordPackage(goFileAnalysis, toAnalyzeContent)

	// 解析依赖包
	analyzeGoImportPackage(goFileAnalysis, toAnalyzeContent)

	// 解析函数定义
	analyzeGoFunctionDefinition(goFileAnalysis, toAnalyzeContent)

	// 解析函数体
	analyzeGoFunctionBody(goFileAnalysis, toAnalyzeContent)

	utility2.TestOutput("goFileAnalysis = %+v", goFileAnalysis)

	// 输出解析结果
	// functionListContent := outputAnalyzeGoFileResult(goFileAnalysis)

	// 输出到文件
	// if toWriteFile != nil {
	// 	_, writeError := toWriteFile.WriteString(functionListContent)
	// 	if writeError != nil {
	// 		return nil, writeError
	// 	}
	// }

	return goFileAnalysis, nil
}

func removeGoFileCommentLine(fileContentByte []byte) []byte {
	if goCommentLineRegexp, hasGoCommentLineRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEGoCommentLine]; hasGoCommentLineRegexp {
		result := goCommentLineRegexp.ReplaceAll(fileContentByte, []byte(""))
		return result
	}
	return fileContentByte
}

func analyzeGoKeywordPackage(goFileAnalysis *goFileAnalysis, fileContentByte []byte) {
	if keywordPackageRegexp, hasKeywordPackageRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEGoKeywordPackageValue]; hasKeywordPackageRegexp {
		if packageValueContentByte := keywordPackageRegexp.Find(fileContentByte); len(packageValueContentByte) != 0 {
			goFileAnalysis.PackageName = strings.Split(string(packageValueContentByte), " ")[1]
			if goFileAnalysis.PackageName == "main" {
				goFileAnalysis.PackagePath = strings.Replace(strings.Replace(filepath.Join(filepath.Dir(goFileAnalysis.FilePath), "main"), "\\", "/", -1), global.GoPathSrc, "", -1)
			} else {
				goFileAnalysis.PackagePath = strings.Replace(strings.Replace(filepath.Dir(goFileAnalysis.FilePath), "\\", "/", -1), global.GoPathSrc, "", -1)
			}
		}
	} else {
		ui.OutputWarnInfo(ui.CMDAnalyzeGoKeywordRegexpNotExist, "package")
	}
}

func analyzeGoImportPackage(goFileAnalysis *goFileAnalysis, fileContentByte []byte) {
	if keywordImportValueRegexp := regexps.GetRegexpByTemplateEnum(global.GoKeywordImportValueTemplate); keywordImportValueRegexp != nil {
		if goKeywordImportAliasRegexp := regexps.GetRegexpByTemplateEnum(global.GoKeywordImportAliasTemplate); goKeywordImportAliasRegexp != nil {
			keywordImportValueSubMatchNameIndexMap := make(map[string]int)
			for index, subMatchName := range keywordImportValueRegexp.SubexpNames() {
				keywordImportValueSubMatchNameIndexMap[subMatchName] = index
			}
			goKeywordImportAliasSubMatchNameIndexMap := make(map[string]int)
			for index, subMatchName := range goKeywordImportAliasRegexp.SubexpNames() {
				goKeywordImportAliasSubMatchNameIndexMap[subMatchName] = index
			}
			for _, keywordImportValueSubmatchList := range keywordImportValueRegexp.FindAllSubmatch(fileContentByte, -1) {
				packageName := make([]byte, 0)
				if aliasIndex, hasAliasIndex := keywordImportValueSubMatchNameIndexMap["ALIAS"]; hasAliasIndex {
					packageName = keywordImportValueSubmatchList[aliasIndex]
				}
				if keywordImportValueSubMatchIndex, hasKeywordImportValueSubMatchIndex := keywordImportValueSubMatchNameIndexMap["VALUE"]; hasKeywordImportValueSubMatchIndex {
					keywordImportValueSubmatchValue := keywordImportValueSubmatchList[keywordImportValueSubMatchIndex]
					// 处理多行形式下的包引用
					for _, goKeywordImportAliasSubMatchList := range goKeywordImportAliasRegexp.FindAllSubmatch(keywordImportValueSubmatchValue, -1) {
						if goKeywordImportAliasSubMatchIndex, hasGoKeywordImportAliasSubMatchIndex := goKeywordImportAliasSubMatchNameIndexMap["ALIAS"]; hasGoKeywordImportAliasSubMatchIndex {
							packageName = goKeywordImportAliasSubMatchList[goKeywordImportAliasSubMatchIndex]
						}
						packagePath := make([]byte, 0)
						if goKeywordImportAliasSubMatchIndex, hasGoKeywordImportAliasSubMatchIndex := goKeywordImportAliasSubMatchNameIndexMap["CONTENT"]; hasGoKeywordImportAliasSubMatchIndex && len(goKeywordImportAliasSubMatchList) > goKeywordImportAliasSubMatchIndex {
							packagePath = goKeywordImportAliasSubMatchList[goKeywordImportAliasSubMatchIndex]
						}
						if len(packageName) == 0 {
							packagePathList := strings.Split(strings.Trim(string(packagePath), "\""), "/")
							packageName = []byte(packagePathList[len(packagePathList)-1])
						}
						goFileAnalysis.ImportAliasMap[strings.TrimSpace(string(packageName))] = strings.Trim(string(packagePath), "\"")
						// clear package name
						packageName = make([]byte, 0)
					}
				}
			}
		} else {
			ui.OutputWarnInfo(ui.CommonError16, global.GoKeywordImportAliasTemplate)
		}

	} else {
		ui.OutputWarnInfo(ui.CMDAnalyzeGoKeywordRegexpNotExist, "import")
	}
}

func analyzeGoFunctionDefinition(goFileAnalysis *goFileAnalysis, fileContentByte []byte) {
	bracketsContentRegexp, hasBracketsContentRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEBracketsContent]
	if !hasBracketsContentRegexp {
		ui.OutputWarnInfo(ui.CommonError16, global.AEBracketsContent)
		return
	}
	functionDefinitionRegexp := regexps.GetRegexpByTemplateEnum(global.GoFunctionDefinitionTemplate)
	if functionDefinitionRegexp == nil {
		ui.OutputWarnInfo(ui.CMDAnalyzeGoKeywordRegexpNotExist, "function")
		return
	}

	// 解析函数定义
	functionDefinitionByteList := functionDefinitionRegexp.FindAll(fileContentByte, -1)
	for _, functionDefinitionByte := range functionDefinitionByteList {
		functionAnalysis := &GoFunctionAnalysis{
			ParamsMap:           make(map[string]*GoFunctionVariable),
			ReturnMap:           make(map[string]*GoFunctionVariable),
			InnerPackageCallMap: make(map[string]map[int]*GoFunctionCallAnalysis),
			OuterPackageCallMap: make(map[string]map[string]map[int]*GoFunctionCallAnalysis),
			MemberCallMap:       make(map[string]map[string]map[int]*GoFunctionCallAnalysis),
			CallerMap:           make(map[string]map[string]map[int]*GoFunctionCallAnalysis),
			VariableMap:         make(map[string]*GoFunctionVariable),
		}
		// 解析所属类
		memberStringWithPunctuationMark := functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$MEMBER")
		if memberStringWithPunctuationMark != "" {
			memberString := bracketsContentRegexp.ReplaceAllString(memberStringWithPunctuationMark, "$CONTENT")
			memberStringList := strings.Split(strings.TrimSpace(memberString), " ")
			if len(memberStringList) != 2 {
				ui.OutputWarnInfo(ui.CMDAnalyzeGoFunctionDefinitionSyntaxError)
				continue
			}
			functionAnalysis.ClassValue = memberStringList[0]
			functionAnalysis.ClassValueType = memberStringList[1]
			functionAnalysis.Class = utility.TraitStructName(memberStringList[1])
		}
		// 解析函数名
		functionAnalysis.Name = functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$NAME")
		// 解析参数表
		paramsStringWithPunctuationMark := functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$PARAM")
		if paramsStringWithPunctuationMark != "" {
			paramsString := bracketsContentRegexp.ReplaceAllString(paramsStringWithPunctuationMark, "$CONTENT")
			typeString := ""
			unknownTypeVariableList := make([]*GoFunctionVariable, 0)
			for index, paramString := range strings.Split(paramsString, ",") {
				paramStringList := strings.Split(strings.TrimSpace(paramString), " ")
				if len(paramStringList) == 1 {
					unknownTypeVariableList = append(unknownTypeVariableList, &GoFunctionVariable{
						Index: index,
						Name:  paramStringList[0],
					})
				} else {
					functionAnalysis.ParamsMap[paramStringList[0]] = &GoFunctionVariable{
						Index: index,
						Name:  paramStringList[0],
						Type:  paramStringList[1],
					}
					if typeString == "" {
						for _, unknownTypeVariable := range unknownTypeVariableList {
							unknownTypeVariable.Type = paramStringList[1]
							functionAnalysis.ParamsMap[unknownTypeVariable.Name] = unknownTypeVariable
						}
					}
				}
			}
		}
		// 解析返回列表
		returnStringWithPunctuationMark := functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$RETURN")
		if returnStringWithPunctuationMark != "" {
			returnContent := bracketsContentRegexp.ReplaceAllString(returnStringWithPunctuationMark, "$CONTENT")
			for index, returnString := range strings.Split(returnContent, ",") {
				returnStringList := strings.Split(strings.TrimSpace(returnString), " ")
				if len(returnStringList) == 1 {
					functionAnalysis.ReturnMap[fmt.Sprintf("%v", len(functionAnalysis.ReturnMap))] = &GoFunctionVariable{
						Index: index,
						Name:  "",
						Type:  returnStringList[0],
					}
				} else if len(returnStringList) == 2 {
					functionAnalysis.ReturnMap[returnStringList[0]] = &GoFunctionVariable{
						Index: index,
						Name:  returnStringList[0],
						Type:  returnStringList[1],
					}
				}
			}
		}
		goFileAnalysis.FunctionMap[functionAnalysis.Name] = functionAnalysis
		goFileAnalysis.functionList = append(goFileAnalysis.functionList, functionAnalysis.Name)
	}
}

func analyzeGoFunctionBody(goFileAnalysis *goFileAnalysis, fileContentByte []byte) {
	functionDefinitionRegexp := regexps.GetRegexpByTemplateEnum(global.GoFunctionDefinitionTemplate)
	if functionDefinitionRegexp == nil {
		ui.OutputWarnInfo(ui.CMDAnalyzeGoKeywordRegexpNotExist, global.GoFunctionDefinitionTemplate)
		return
	}

	// variableDeclarationRegexp := regexps.GetRegexpByTemplateEnum(global.GoVariableDeclarationTemplate)
	// if variableDeclarationRegexp == nil {
	// 	ui.OutputWarnInfo(ui.CMDAnalyzeGoKeywordRegexpNotExist, global.GoVariableDeclarationTemplate)
	// 	return
	// }

	variableInitializationRegexp := regexps.GetRegexpByTemplateEnum(global.GoVariableInitializationTemplate)
	if variableInitializationRegexp == nil {
		ui.OutputWarnInfo(ui.CMDAnalyzeGoKeywordRegexpNotExist, global.GoVariableInitializationTemplate)
		return
	}

	goFunctionCallRegexp := regexps.GetRegexpByTemplateEnum(global.GoFunctionCallTemplate)
	if goFunctionCallRegexp == nil {
		ui.OutputWarnInfo(ui.CMDAnalyzeGoKeywordRegexpNotExist, global.GoFunctionCallTemplate)
		return
	}

	functionDefinitionStartIndexList := functionDefinitionRegexp.FindAllIndex(fileContentByte, -1)
	for index, functionDefinitionStartIndex := range functionDefinitionStartIndexList {
		if index > len(goFileAnalysis.functionList[index]) {
			ui.OutputWarnInfo(ui.CMDAnalyzeGoFunctionError, index)
			continue
		}
		goFunctionAnalysis := goFileAnalysis.FunctionMap[goFileAnalysis.functionList[index]]
		utility2.TestOutput("index = %v, functionDefinitionStartIndex = %v, goFunctionAnalysis.Name = %v", index, functionDefinitionStartIndex, goFunctionAnalysis.Name)

		// passing content: ...}
		definitionLength := utility2.CalculatePunctuationMarksContentLength(string(fileContentByte[functionDefinitionStartIndex[1]:]), global.PunctuationMarkCurlyBraces)
		if definitionLength == 0 {
			ui.OutputWarnInfo(ui.CMDAnalyzeGoFunctionContentSyntaxError)
			continue
		}

		goFunctionAnalysis.BodyContent = fileContentByte[functionDefinitionStartIndex[1]-1 : functionDefinitionStartIndex[1]+1+definitionLength]

		// body content: {...}
		// utility2.TestOutput(ui.CommonNote2)
		// utility2.TestOutput("goFunctionAnalysis.BodyContent = |%v|", string(goFunctionAnalysis.BodyContent))
		utility2.TestOutput(ui.CommonNote2)

		// variable

		// variableDeclarationIndexList := variableDeclarationRegexp.FindAll(goFunctionAnalysis.BodyContent, -1)
		// for _, variableDeclaration := range variableDeclarationIndexList {
		// 	utility2.TestOutput("variableDeclaration = %v", string(variableDeclaration))
		// 	variableName := variableDeclarationRegexp.ReplaceAllString(string(variableDeclaration), "$NAME")
		// 	variableType := variableDeclarationRegexp.ReplaceAllString(string(variableDeclaration), "$TYPE")
		// 	utility2.TestOutput("variableName = %v", variableName)
		// 	utility2.TestOutput("variableType = %v", variableType)
		// 	goFunctionAnalysis.VariableMap[variableName] = &GoFunctionVariable{
		// 		Index: len(goFunctionAnalysis.VariableMap),
		// 		Name:  variableName,
		// 		Type:  variableType,
		// 	}
		// }

		utility2.TestOutput(ui.CommonNote2)

		variableInitializationIndexList := variableInitializationRegexp.FindAll(goFunctionAnalysis.BodyContent, -1)
		for _, variableInitialization := range variableInitializationIndexList {
			utility2.TestOutput("variableInitialization = %v", string(variableInitialization))
			variableListString := strings.TrimSpace(variableInitializationRegexp.ReplaceAllString(string(variableInitialization), "$LIST"))
			variableInitString := strings.TrimSpace(variableInitializationRegexp.ReplaceAllString(string(variableInitialization), "$INIT"))
			utility2.TestOutput("variableList = %v", variableListString)
			utility2.TestOutput("variableInit = %v", variableInitString)

			// parse variable init by what? -> struct or function call or other variable function name
			variableTypeList := make([]string, 0)
			punctuationMarkLeftBracket := strings.IndexRune(variableInitString, global.PunctuationMarkLeftBracket)
			punctuationMarkLeftCurlyBracesIndex := strings.IndexRune(variableInitString, global.PunctuationMarkLeftCurlyBraces)
			punctuationMarkPointIndex := strings.IndexRune(variableInitString, global.PunctuationMarkPoint)
			if punctuationMarkLeftBracket == -1 && punctuationMarkLeftCurlyBracesIndex == -1 {
				if variableAnalysis, isParamVariable := goFunctionAnalysis.ParamsMap[variableInitString]; isParamVariable {
					utility2.TestOutput("variableInit is function param %v:%v", variableInitString, variableAnalysis.Type)
					variableTypeList = append(variableTypeList, variableAnalysis.Type)
				} else if variableAnalysis, isVariable := goFunctionAnalysis.VariableMap[variableInitString]; isVariable {
					utility2.TestOutput("variableInit is function variable %v:%v", variableInitString, variableAnalysis.Type)
					variableTypeList = append(variableTypeList, variableAnalysis.Type)
				} else {
					if punctuationMarkPointIndex == -1 {
						// TODO:
						utility2.TestOutput("variableInit is a function or variable from inner package", variableInitString)
					} else {
						utility2.TestOutput("variableInit is a function or variable from %v", variableInitString, variableAnalysis.Type)
					}
				}
			} else if punctuationMarkLeftBracket != -1 && punctuationMarkLeftCurlyBracesIndex == -1 {
				utility2.TestOutput("variableInit is function call return %v", variableInitString)
				// if punctuationMarkPointIndex != -1 {

				// }
			} else if punctuationMarkLeftBracket == -1 && punctuationMarkLeftCurlyBracesIndex != -1 {
				utility2.TestOutput("variableInit is struct instance: %v", variableInitString[0:punctuationMarkLeftCurlyBracesIndex])
			} else {
				utility2.TestOutput("variableInit is unsupport syntax")
			}

			for _, variableString := range strings.Split(variableListString, ",") {
				if strings.TrimSpace(variableString) != "_" {
					goFunctionAnalysis.VariableMap[strings.TrimSpace(variableString)] = &GoFunctionVariable{
						Index: len(goFunctionAnalysis.VariableMap),
						Name:  strings.TrimSpace(variableString),
						// Type:  variableType,
					}
				}
			}
		}

		continue

		utility2.TestOutput(ui.CommonNote2)

		// function call

		for callIndex, functionCallByteList := range goFunctionCallRegexp.FindAll(goFunctionAnalysis.BodyContent, -1) {
			utility2.TestOutput("callIndex = %v, string(functionCallByteList) = %v", callIndex, string(functionCallByteList))
			utility2.TestOutput("goFileAnalysis.PackageName = %v", goFileAnalysis.PackageName)

			callFrom := goFunctionCallRegexp.ReplaceAllString(string(functionCallByteList), "$CALL")
			callFunction := goFunctionCallRegexp.ReplaceAllString(string(functionCallByteList), "$NAME")
			callParam := goFunctionCallRegexp.ReplaceAllString(string(functionCallByteList), "$PARAM")
			var callParamList []string

			if len(callParam) != 0 {
				callParamStringList := strings.Split(callParam, ",")
				callParamList = make([]string, 0, len(callParamStringList))
				for _, callParam := range callParamStringList {
					callParamList = append(callParamList, utility2.FixBracketMatchingResult(strings.TrimSpace(callParam)))
				}
			}

			if len(callFrom) != 0 {
				if _, hasPackage := goFileAnalysis.ImportAliasMap[callFrom]; hasPackage {
					callFrom = goFileAnalysis.ImportAliasMap[callFrom]
					if _, hasPackage := goFunctionAnalysis.OuterPackageCallMap[callFrom]; !hasPackage {
						goFunctionAnalysis.OuterPackageCallMap[callFrom] = make(map[string]map[int]*GoFunctionCallAnalysis)
					}
					if _, hasCallFunction := goFunctionAnalysis.OuterPackageCallMap[callFrom][callFunction]; !hasCallFunction {
						goFunctionAnalysis.OuterPackageCallMap[callFrom][callFunction] = make(map[int]*GoFunctionCallAnalysis)
					}
					goFunctionAnalysis.OuterPackageCallMap[callFrom][callFunction][len(goFunctionAnalysis.OuterPackageCallMap[callFrom][callFunction])] = &GoFunctionCallAnalysis{
						From:      callFrom,
						Content:   string(functionCallByteList),
						ParamList: callParamList,
					}
					utility2.TestOutput("call from outer package %v, call function %v, call param %v", callFrom, callFunction, callParamList)
				} else {
					if _, hasMember := goFunctionAnalysis.MemberCallMap[callFrom]; !hasMember {
						goFunctionAnalysis.MemberCallMap[callFrom] = make(map[string]map[int]*GoFunctionCallAnalysis)
					}
					if _, hasCallFunction := goFunctionAnalysis.MemberCallMap[callFrom][callFunction]; !hasCallFunction {
						goFunctionAnalysis.MemberCallMap[callFrom][callFunction] = make(map[int]*GoFunctionCallAnalysis)
					}
					goFunctionAnalysis.MemberCallMap[callFrom][callFunction][len(goFunctionAnalysis.MemberCallMap[callFrom][callFunction])] = &GoFunctionCallAnalysis{
						From:      callFrom,
						Content:   string(functionCallByteList),
						ParamList: callParamList,
					}
					utility2.TestOutput("call from member %v, call function %v, call param %v", callFrom, callFunction, callParamList)
				}
			} else {
				if _, hasCallFunction := goFunctionAnalysis.InnerPackageCallMap[callFunction]; !hasCallFunction {
					goFunctionAnalysis.InnerPackageCallMap[callFunction] = make(map[int]*GoFunctionCallAnalysis)
				}
				goFunctionAnalysis.InnerPackageCallMap[callFunction][len(goFunctionAnalysis.InnerPackageCallMap[callFunction])] = &GoFunctionCallAnalysis{
					From:      callFrom,
					Content:   string(functionCallByteList),
					ParamList: callParamList,
				}
				utility2.TestOutput("call from inner package, call function %v, call param %v", callFunction, callParamList)
			}

			// callFunctionAnalysis := goFuncion

			utility2.TestOutput(ui.CommonNote2)
		}
	}
}

func outputAnalyzeGoFileResult(goFileAnalysis *goFileAnalysis) string {
	templateStyleRegexp, hasTemplateStyleRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AETemplateStyle]
	if !hasTemplateStyleRegexp {
		ui.OutputWarnInfo(ui.CommonError16, global.AETemplateStyle)
		return ""
	}

	// file path
	resultContent := strings.Replace(ui.AnalyzeGoFileResultTemplate, global.AnalyzeRPFilePath, goFileAnalysis.FilePath, -1)

	// package content
	packageContent := ui.ParseStyleTemplate(templateStyleRegexp, ui.AnalyzeGoFilePackageContentTemplate)
	packageContent = strings.Replace(packageContent, global.AnalyzeRPPackageName, goFileAnalysis.PackageName, -1)
	packageContent = strings.Replace(packageContent, global.AnalyzeRPPackagePath, goFileAnalysis.PackagePath, -1)
	resultContent = strings.Replace(resultContent, global.AnalyzeRPPackageContent, packageContent, -1)

	// import 内容
	importPackageListContent := ""
	if len(goFileAnalysis.ImportAliasMap) != 0 {
		importPackageListString := ""
		for packageAlias, packagePath := range goFileAnalysis.ImportAliasMap {
			packageAliasPathContent := ui.AnalyzeGoFileImportPackageTemplate

			// style template
			packageAliasPathContent = ui.ParseStyleTemplate(templateStyleRegexp, packageAliasPathContent)

			// package alias
			packageAliasPathContent = strings.Replace(packageAliasPathContent, global.AnalyzeRPPackageAlias, packageAlias, -1)

			// package path
			packageAliasPathContent = strings.Replace(packageAliasPathContent, global.AnalyzeRPPackagePath, packagePath, -1)

			if importPackageListString == "" {
				importPackageListString = packageAliasPathContent
			} else {
				importPackageListString = fmt.Sprintf("%v\n%v", importPackageListString, packageAliasPathContent)
			}
		}
		importPackageListContent = strings.Replace(ui.AnalyzeGoFileImportPackageListTemplate, global.AnalyzeRPImportPackage, importPackageListString, -1)
	}
	resultContent = strings.Replace(resultContent, global.AnalyzeRPImportPackageList, importPackageListContent, -1)

	// function 内容
	functionDefinitionListContent := ""
	if len(goFileAnalysis.functionList) != 0 {
		functionDefinitionListString := ""
		for _, functionName := range goFileAnalysis.functionList {
			functionAnalysis := goFileAnalysis.FunctionMap[functionName]
			functionDefinitionContent := ui.AnalyzeGoFileFunctionDefinitionTemplate

			// sytle template
			functionDefinitionContent = ui.ParseStyleTemplate(templateStyleRegexp, functionDefinitionContent)

			// function name
			functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionName, functionAnalysis.Name, -1)

			// function class
			if functionAnalysis.Class != "" {
				functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionClass, strings.Replace(ui.AnalyzeGoFileFunctionClassTemplate, global.AnalyzeRPFunctionClassName, functionAnalysis.Class, -1), -1)
			} else {
				functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionClass, global.AnalyzeRPEmptyString, -1)
			}

			// function params
			functionParamListContent := ""
			if len(functionAnalysis.ParamsMap) != 0 {
				// functionParamListContent = parseGoFunctionParamOrReturnListContent(templateStyleRegexp, functionAnalysis.ParamsMap, functionParamList)
			}
			functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionParamList, functionParamListContent, -1)

			// function return
			functionReturnListContent := ""
			if len(functionAnalysis.ReturnMap) != 0 {
				// functionReturnListContent = parseGoFunctionParamOrReturnListContent(templateStyleRegexp, functionAnalysis.ReturnMap, functionReturnList)
			}
			functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionReturnList, functionReturnListContent, -1)

			// function body call map
			functionCallMapContent := ""
			if len(functionAnalysis.OuterPackageCallMap) != 0 {
				// functionCallMapContent = parseGoFunctionCallMapContent(templateStyleRegexp, functionAnalysis.OuterPackageCallMap)
			}
			functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionCallMap, functionCallMapContent, -1)

			// add to function definition list content
			if functionDefinitionListString == "" {
				functionDefinitionListString = functionDefinitionContent
			} else {
				functionDefinitionListString = fmt.Sprintf("%v\n%v", functionDefinitionListString, functionDefinitionContent)
			}
		}
		functionDefinitionListContent = strings.Replace(ui.AnalyzeGoFileFunctionDefinitionListTemplate, global.AnalyzeRPFunctionDefinition, functionDefinitionListString, -1)
	}
	resultContent = strings.Replace(resultContent, global.AnalyzeRPFunctionDefinitionList, functionDefinitionListContent, -1)

	// clear space line
	resultContent = utility3.TrimSpaceLine(resultContent)

	return resultContent
}

func outputGoPackageLevelDirectedGraph(mergedNTree map[int]map[int]int) string {
	return ""
}

const (
	functionParamList  = 0
	functionReturnList = 1
)

func parseGoFunctionParamOrReturnListContent(templateStyleRegexp *regexp.Regexp, nameTypeMap map[string]string, parseType int) string {
	nameTypeListTemplate := ui.AnalyzeGoFileFunctionParamListTemplate
	nameTypeTemplate := ui.AnalyzeGoFileFunctionParamNameTypeTemplate
	replaceString := global.AnalyzeRPFunctionParamNameTypeList
	if parseType == functionReturnList {
		nameTypeListTemplate = ui.AnalyzeGoFileFunctionReturnListTemplate
		nameTypeTemplate = ui.AnalyzeGoFileFunctionReturnNameTypeTemplate
		replaceString = global.AnalyzeRPFunctionReturnNameTypeList
	}

	nameTypeListContent := ""
	for nameString, typeString := range nameTypeMap {
		nameTypeContent := ui.ParseStyleTemplate(templateStyleRegexp, nameTypeTemplate)
		nameTypeContent = strings.Replace(nameTypeContent, global.AnalyzeRPName, nameString, -1)
		nameTypeContent = strings.Replace(nameTypeContent, global.AnalyzeRPType, typeString, -1)
		if nameTypeListContent == "" {
			nameTypeListContent = nameTypeContent
		} else {
			nameTypeListContent = fmt.Sprintf("%v\n%v", nameTypeListContent, nameTypeContent)
		}
	}
	return strings.Replace(ui.ParseStyleTemplate(templateStyleRegexp, nameTypeListTemplate), replaceString, nameTypeListContent, -1)
}

func parseGoFunctionCallMapContent(templateStyleRegexp *regexp.Regexp, nameTypeMap map[string]map[string]int) string {
	callMapContent := ""
	for callPackage, functionMap := range nameTypeMap {
		for functionName := range functionMap {
			packageFunctionContent := ui.ParseStyleTemplate(templateStyleRegexp, ui.AnalyzeGoFileFunctionCallPackageMap)
			packageFunctionContent = strings.Replace(packageFunctionContent, global.AnalyzeRPPackage, callPackage, -1)
			packageFunctionContent = strings.Replace(packageFunctionContent, global.AnalyzeRPIdentifier, functionName, -1)
			if callMapContent == "" {
				callMapContent = packageFunctionContent
			} else {
				callMapContent = fmt.Sprintf("%v\n%v", callMapContent, packageFunctionContent)
			}
		}
	}
	return strings.Replace(ui.ParseStyleTemplate(templateStyleRegexp, ui.AnalyzeGoFileFunctionCallMapTemplate), global.AnalyzeRPFunctionCallPackageMap, callMapContent, -1)
}

func makeUpNTreeNodeChildrenMapByGoPackage(goPackageAnalysisMap map[int]*GoPackageAnalysis) map[int][]int {
	nTreeNodeChildrenMap := make(map[int][]int)
	// for packageNo, packageAnalysis := range goPackageAnalysisMap {
	// 	if _, hasNo := nTreeNodeChildrenMap[packageNo]; !hasNo {
	// 		nTreeNodeChildrenMap[packageNo] = make([]int, 0)
	// 	}
	// 	nTreeNodeChildrenMap[packageNo] = append(nTreeNodeChildrenMap[packageNo], packageAnalysis.ImportPackageAnalysisNoList...)
	// }
	return nTreeNodeChildrenMap
}

// 分析 CPP 文件

func analyzeCppFile(toAnalyzeFile, toWriteFile *os.File) error {
	return nil
}

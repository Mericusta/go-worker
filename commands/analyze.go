package commands

import (
	"fmt"
	"os"
	"path"
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

var goAnalyzerPackageSubMatchNameIndexMap map[string]int

// var goAnalyzerSingleLineImportSubMatchNameIndexMap map[string]int
var goAnalyzerImportContentSubMatchNameIndexMap map[string]int
var goAnalyzerPackageVariableSubMatchNameIndexMap map[string]int
var goAnalyzerInterfaceFunctionSubMatchNameIndexMap map[string]int
var goAnalyzerInterfaceSubMatchNameIndexMap map[string]int
var goAnalyzerStructVariableSubMatchNameIndexMap map[string]int
var goAnalyzerStructFunctionSubMatchNameIndexMap map[string]int
var goAnalyzerFunctionSubMatchNameIndexMap map[string]int
var goAnalyzerMemberFunctionSubMatchNameIndexMap map[string]int
var goAnalyzerTypeRenameSubMatchNameIndexMap map[string]int
var goAnalyzerMultiLineConstSubMatchNameIndexMap map[string]int
var goAnalyzerSingleLineConstSubMatchNameIndexMap map[string]int

// GoAnalysis go 项目分析结果
type GoAnalysis struct {
	RootPath           string
	PackageAnalysisMap map[string]*GoPackageAnalysis // package path : package analysis
}

// GoPackageAnalysis go 包分析结果
type GoPackageAnalysis struct {
	Scope                                 map[string]*GoPackageScope                           // file path : go package scope
	PackageName                           string                                               // package name
	PackagePath                           string                                               // package path
	ImportAnalysis                        map[string]map[string]*GoImportAnalysis              // file path : package path : import analysis
	VariableAnalysisMap                   map[string]*GoVariableAnalysis                       // variable name : variable analysis
	ConstAnalysisMap                      map[string]*GoVariableAnalysis                       // const name : variable analysis
	InterfaceAnalysis                     map[string]*GoInterfaceAnalysis                      // interface name : interface analysis
	StructAnalysis                        map[string]*GoStructAnalysis                         // struct name : struct analysis
	FunctionAnalysisMap                   map[string]*GoFunctionAnalysis                       // function name : function analysis
	OtherPackageMemberFunctionAnalysisMap map[string]map[string]map[string]*GoFunctionAnalysis // package path : struct name : function name : function analysis
}

// GoImportAnalysis go 引入包的分析结果
type GoImportAnalysis struct {
	Alias string
	Path  string
}

// GoVariableAnalysis go 变量分析结果
type GoVariableAnalysis struct {
	Name     string
	Type     string
	TypeFrom string
}

// GoInterfaceAnalysis go 接口分析结果
type GoInterfaceAnalysis struct {
	Name     string
	Function map[string]*GoFunctionAnalysis
}

// GoStructAnalysis go 结构体分析结果
type GoStructAnalysis struct {
	Name           string
	Base           map[string]*GoVariableAnalysis // base struct type : base variable analysis
	MemberVariable map[string]*GoVariableAnalysis // valiable name : variable analysis
	MemberFunction map[string]*GoFunctionAnalysis // function name : function analysis
}

// GoFunctionAnalysis go 函数分析结果
type GoFunctionAnalysis struct {
	Class               string                                                // the struct's name of its belong struct, if it is member function
	ClassValue          string                                                // the name of its belong struct's member 'this', if it is member function
	ClassValueType      string                                                // the type of its belong struct's member 'this', if it is member function
	ClassValueTypeFrom  string                                                // the package of its belong struct's member 'this', if it is member function
	Name                string                                                // function name
	ParamsMap           map[string]*GoFunctionVariable                        // param name : function variable
	ReturnMap           map[string]*GoFunctionVariable                        // return name : function variable
	VariableMap         map[string]*GoFunctionVariable                        // variable name : function variable
	CallMap             map[string][]*GoFunctionCallAnalysis                  // call function : function call analysis
	InnerPackageCallMap map[string]map[int]*GoFunctionCallAnalysis            // call function from inner package
	OuterPackageCallMap map[string]map[string]map[int]*GoFunctionCallAnalysis // call function from outer package
	MemberCallMap       map[string]map[string]map[int]*GoFunctionCallAnalysis // call member function
	CallerMap           map[string]map[string]map[int]*GoFunctionCallAnalysis // package path : function : call index : call content
}

// GoFunctionCallAnalysis go 函数调用分析结果
type GoFunctionCallAnalysis struct {
	Content   string
	From      string
	Call      string
	ParamList []string
}

// GoFunctionVariable go 函数内变量
type GoFunctionVariable struct {
	GoVariableAnalysis
	Index int
}

// checkGoAnalyzerRegexp 检查 go 项目分析器的所有模板/原子表达式
func checkGoAnalyzerRegexp() bool {
	ok := true
	if goFileAnalyzerScopePackageRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopePackageTemplate); goFileAnalyzerScopePackageRegexp != nil {
		goAnalyzerPackageSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileAnalyzerScopePackageRegexp.SubexpNames() {
			goAnalyzerPackageSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileAnalyzerScopePackageTemplate)
		ok = false
	}

	// if goFileSplitterScopeMultiLineImportStartRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeMultiLineImportStartTemplate); goFileSplitterScopeMultiLineImportStartRegexp == nil {
	// 	ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopeMultiLineImportStartTemplate)
	// 	ok = false
	// }

	if goFileAnalyzerScopeImportContentRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopeImportContentTemplate); goFileAnalyzerScopeImportContentRegexp != nil {
		goAnalyzerImportContentSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileAnalyzerScopeImportContentRegexp.SubexpNames() {
			goAnalyzerImportContentSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileAnalyzerScopeImportContentTemplate)
		ok = false
	}

	// if goFileSplitterScopeSingleLineImportRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileSplitterScopeSingleLineImportTemplate); goFileSplitterScopeSingleLineImportRegexp != nil {
	// 	goSplitterSingleLineImportSubMatchNameIndexMap = make(map[string]int)
	// 	for index, subMatchName := range goFileSplitterScopeSingleLineImportRegexp.SubexpNames() {
	// 		goSplitterSingleLineImportSubMatchNameIndexMap[subMatchName] = index
	// 	}
	// } else {
	// 	ui.OutputErrorInfo(ui.CommonError18, global.GoFileSplitterScopeSingleLineImportTemplate)
	// 	ok = false
	// }

	// if goFileSplitterScopeEnd, has := regexps.AtomicExpressionEnumRegexpMap[global.AEGoFileSplitterScopeEnd]; !has || goFileSplitterScopeEnd == nil {
	// 	ui.OutputErrorInfo(ui.CommonError16, global.AEGoFileSplitterScopeEnd)
	// 	ok = false
	// }

	if goFileAnalyzerScopePackageVariableRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopePackageVariableTemplate); goFileAnalyzerScopePackageVariableRegexp != nil {
		goAnalyzerPackageVariableSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileAnalyzerScopePackageVariableRegexp.SubexpNames() {
			goAnalyzerPackageVariableSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileAnalyzerScopePackageVariableTemplate)
		ok = false
	}

	if goFileAnalyzerScopeInterfaceFunctionRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopeInterfaceFunctionTemplate); goFileAnalyzerScopeInterfaceFunctionRegexp != nil {
		goAnalyzerInterfaceFunctionSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileAnalyzerScopeInterfaceFunctionRegexp.SubexpNames() {
			goAnalyzerInterfaceFunctionSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileAnalyzerScopeInterfaceFunctionTemplate)
		ok = false
	}

	if goFileAnalyzerScopeInterfaceRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopeInterfaceTemplate); goFileAnalyzerScopeInterfaceRegexp != nil {
		goAnalyzerInterfaceSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileAnalyzerScopeInterfaceRegexp.SubexpNames() {
			goAnalyzerInterfaceSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileAnalyzerScopeInterfaceTemplate)
		ok = false
	}

	if goFileAnalyzerScopeStructVariableRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopeStructVariableTemplate); goFileAnalyzerScopeStructVariableRegexp != nil {
		goAnalyzerStructVariableSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileAnalyzerScopeStructVariableRegexp.SubexpNames() {
			goAnalyzerStructVariableSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileAnalyzerScopeStructVariableTemplate)
		ok = false
	}

	// if goFileAnalyzerScopeStructFunctionRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopeStructFunctionTemplate); goFileAnalyzerScopeStructFunctionRegexp != nil {
	// 	goAnalyzerStructFunctionSubMatchNameIndexMap = make(map[string]int)
	// 	for index, subMatchName := range goFileAnalyzerScopeStructFunctionRegexp.SubexpNames() {
	// 		goAnalyzerStructFunctionSubMatchNameIndexMap[subMatchName] = index
	// 	}
	// } else {
	// 	ui.OutputErrorInfo(ui.CommonError18, global.GoFileAnalyzerScopeStructFunctionTemplate)
	// 	ok = false
	// }

	if goFileAnalyzerScopeFunctionRegexp := regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopeFunctionTemplate); goFileAnalyzerScopeFunctionRegexp != nil {
		goAnalyzerFunctionSubMatchNameIndexMap = make(map[string]int)
		for index, subMatchName := range goFileAnalyzerScopeFunctionRegexp.SubexpNames() {
			goAnalyzerFunctionSubMatchNameIndexMap[subMatchName] = index
		}
	} else {
		ui.OutputErrorInfo(ui.CommonError18, global.GoFileAnalyzerScopeFunctionTemplate)
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

// analyzeGo 分析 go 项目
// @param rootPath 项目根路径
// @param toAnalyzePathList 项目下所有待分析的文件的列表
// @return
func analyzeGo(rootPath string, toAnalyzePathList []string) (*GoAnalysis, error) {
	if !checkGoSplitterRegexp() {
		return nil, nil
	}

	if !checkGoAnalyzerRegexp() {
		return nil, nil
	}

	// const mainPackageNo = 0
	// packageNo := mainPackageNo
	// fileNo := 0
	goAnalysis := &GoAnalysis{
		// FileNoAnalysisMap:    make(map[int]*goFileAnalysis),
		// PackageNoAnalysisMap: make(map[int]*GoPackageAnalysis),
		// PackageFunctionMap:   make(map[string]string),
		RootPath:           rootPath,
		PackageAnalysisMap: make(map[string]*GoPackageAnalysis),
	}
	// packagePathAnalysisMap := make(map[string]*GoPackageAnalysis)

	// 4.1.3.1.6.1
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

		// toAnalyzeFile, inputError := os.Open(toAnalyzeFilePath)
		// defer toAnalyzeFile.Close()
		// if inputError != nil || toAnalyzeFile == nil {
		// 	return goAnalysis, fmt.Errorf(ui.CommonError5, toAnalyzeFilePath, inputError.Error())
		// }

		// fileAnalysis, analyzeError := analyzeGoFile(toAnalyzeFile)
		// if analyzeError != nil {
		// 	ui.OutputWarnInfo(ui.CMDAnalyzeOccursError, analyzeError)
		// }
		// fileAnalysis.No = fileNo

		// 4.1.3.1.6.1.1
		splitFileResult := SplitGoFile(toAnalyzeFilePath, false)

		// 4.1.3.1.6.1.2
		packageName, packagePath, analyzeGoScopePackageError := analyzeGoScopePackage(splitFileResult.Package.Content, rootPath)
		if analyzeGoScopePackageError != nil {
			return nil, analyzeGoScopePackageError
		}

		utility2.TestOutput(ui.CommonNote2)
		utility2.TestOutput("packageName = %v, packagePath = %v", packageName, packagePath)

		if _, hasPackagePath := goAnalysis.PackageAnalysisMap[packagePath]; !hasPackagePath {
			goAnalysis.PackageAnalysisMap[packagePath] = &GoPackageAnalysis{
				PackageName:                           packageName,
				PackagePath:                           packagePath,
				ImportAnalysis:                        make(map[string]map[string]*GoImportAnalysis),
				VariableAnalysisMap:                   make(map[string]*GoVariableAnalysis),
				ConstAnalysisMap:                      make(map[string]*GoVariableAnalysis),
				InterfaceAnalysis:                     make(map[string]*GoInterfaceAnalysis),
				StructAnalysis:                        make(map[string]*GoStructAnalysis),
				FunctionAnalysisMap:                   make(map[string]*GoFunctionAnalysis),
				OtherPackageMemberFunctionAnalysisMap: make(map[string]map[string]map[string]*GoFunctionAnalysis),
				Scope:                                 make(map[string]*GoPackageScope),
			}
		}
		goAnalysis.PackageAnalysisMap[packagePath].Scope[toAnalyzeFilePath] = splitFileResult

		// 4.1.3.1.6.1.3
		if _, hasFile := goAnalysis.PackageAnalysisMap[packagePath].ImportAnalysis[toAnalyzeFilePath]; !hasFile {
			goAnalysis.PackageAnalysisMap[packagePath].ImportAnalysis[toAnalyzeFilePath] = make(map[string]*GoImportAnalysis)
		}

		for _, scopeData := range splitFileResult.SingleLineImport {
			if goImportAnalysis := analyzeGoScopeSingleLineImport(scopeData.Content); goImportAnalysis != nil {
				goAnalysis.PackageAnalysisMap[packagePath].ImportAnalysis[toAnalyzeFilePath][goImportAnalysis.Alias] = goImportAnalysis
			}
		}

		if splitFileResult.MultiLineImport != nil {
			goImportAnalysisList := analyzeGoScopeMultiLineImport(splitFileResult.MultiLineImport.Content)
			for _, goImportAnalysis := range goImportAnalysisList {
				goAnalysis.PackageAnalysisMap[packagePath].ImportAnalysis[toAnalyzeFilePath][goImportAnalysis.Alias] = goImportAnalysis
			}
		}

		utility2.TestOutput(ui.CommonNote2)
		utility2.TestOutput("Output package import analysis:")
		for filePath, importAnlysisMap := range goAnalysis.PackageAnalysisMap[packagePath].ImportAnalysis {
			utility2.TestOutput("filePath = %v", filePath)
			for alias, analysis := range importAnlysisMap {
				utility2.TestOutput("alias = %v, importPath = %v", alias, analysis.Path)
			}
		}
		utility2.TestOutput(ui.CommonNote2)

		// 4.1.3.1.6.1.4
		for _, packageVariableScope := range splitFileResult.PackageVariable {
			if goPackageVariableAnalysis := analyzeGoScopePackageVariable(packageVariableScope.Content); goPackageVariableAnalysis != nil {
				goAnalysis.PackageAnalysisMap[packagePath].VariableAnalysisMap[goPackageVariableAnalysis.Name] = goPackageVariableAnalysis
			}
		}

		utility2.TestOutput(ui.CommonNote2)
		utility2.TestOutput("Output package variable analysis:")
		for name, analysis := range goAnalysis.PackageAnalysisMap[packagePath].VariableAnalysisMap {
			utility2.TestOutput("name = %v, type = %v, type from = %v", name, analysis.Type, analysis.TypeFrom)
		}
		utility2.TestOutput(ui.CommonNote2)

		// 4.1.3.1.6.1.5
		for interfaceName, interfaceScope := range splitFileResult.InterfaceDefinition {
			if goInterfaceAnalysis := analyzeGoScopeInterface(interfaceScope.Content); goInterfaceAnalysis != nil {
				goInterfaceAnalysis.Name = interfaceName
				goAnalysis.PackageAnalysisMap[packagePath].InterfaceAnalysis[interfaceName] = goInterfaceAnalysis
			}
		}

		utility2.TestOutput(ui.CommonNote2)
		utility2.TestOutput("Output package interface analysis:")
		for interfaceName, interfaceAnalysis := range goAnalysis.PackageAnalysisMap[packagePath].InterfaceAnalysis {
			utility2.TestOutput("interface %v function list: %v", interfaceName, len(interfaceAnalysis.Function))
			for functionName, functionAnalysis := range interfaceAnalysis.Function {
				utility2.TestOutput("function: %v", functionName)
				utility2.TestOutput("- param list:")
				for paramName, paramAnalysis := range functionAnalysis.ParamsMap {
					var paramTypeString string
					if len(paramAnalysis.TypeFrom) != 0 {
						paramTypeString = fmt.Sprintf("%v.%v", paramAnalysis.TypeFrom, paramAnalysis.Type)
					} else {
						paramTypeString = paramAnalysis.Type
					}
					utility2.TestOutput("\t- Index: %v, Name: %v, Type: %v", paramAnalysis.Index, paramName, paramTypeString)
				}
				utility2.TestOutput("- return list:")
				for returnIndex, returnAnalysis := range functionAnalysis.ReturnMap {
					var returnTypeString string
					if len(returnAnalysis.TypeFrom) != 0 {
						returnTypeString = fmt.Sprintf("%v.%v", returnAnalysis.TypeFrom, returnAnalysis.Type)
					} else {
						returnTypeString = returnAnalysis.Type
					}
					utility2.TestOutput("\t- Index: %v, Name: %v, Type: %v", returnIndex, returnAnalysis.Name, returnTypeString)
				}
			}
		}
		utility2.TestOutput(ui.CommonNote2)

		// 4.1.3.1.6.1.6
		for structName, structScope := range splitFileResult.StructDefinition {
			if goStructAnalysis := analyzeGoScopeStruct(structScope); goStructAnalysis != nil {
				goStructAnalysis.Name = structName
				goAnalysis.PackageAnalysisMap[packagePath].StructAnalysis[structName] = goStructAnalysis
			}
		}

		utility2.TestOutput(ui.CommonNote2)
		utility2.TestOutput("Output package struct analysis:")
		for structName, structAnalysis := range goAnalysis.PackageAnalysisMap[packagePath].StructAnalysis {
			utility2.TestOutput("struct %v base struct list: %v", structName, len(structAnalysis.Base))
			for baseType, baseVariableAnalysis := range structAnalysis.Base {
				utility2.TestOutput("- base type: %v, from %v", baseType, baseVariableAnalysis.TypeFrom)
			}
			utility2.TestOutput("struct %v variable list: %v", structName, len(structAnalysis.MemberVariable))
			for variableName, variableAnalysis := range structAnalysis.MemberVariable {
				utility2.TestOutput("- variable name: %v, type: %v, from: %v", variableName, variableAnalysis.Type, variableAnalysis.TypeFrom)
			}
		}
		utility2.TestOutput(ui.CommonNote2)

		// 1.1.3.1.6.1.7
		for functionName, functionScope := range splitFileResult.FunctionDefinition {
			if functionAnalysis := analyzeGoScopeFunction(functionScope); functionAnalysis != nil {
				goAnalysis.PackageAnalysisMap[packagePath].FunctionAnalysisMap[functionName] = functionAnalysis
			}
		}

		utility2.TestOutput(ui.CommonNote2)
		utility2.TestOutput("Output package function analysis:")
		for functionName, functionAnalysis := range goAnalysis.PackageAnalysisMap[packagePath].FunctionAnalysisMap {
			utility2.TestOutput("function %v", functionName)
			for paramName, functionVariableAnalysis := range functionAnalysis.ParamsMap {
				utility2.TestOutput("- param index: %v, name: %v, type: %v, from: %v", functionVariableAnalysis.Index, paramName, functionVariableAnalysis.Type, functionVariableAnalysis.TypeFrom)
			}
			for _, functionVariableAnalysis := range functionAnalysis.ReturnMap {
				utility2.TestOutput("- return index: %v, name: %v, type: %v, from: %v", functionVariableAnalysis.Index, functionVariableAnalysis.Name, functionVariableAnalysis.Type, functionVariableAnalysis.TypeFrom)
			}
			for call, callAnalysisList := range functionAnalysis.CallMap {
				utility2.TestOutput("- call: %v, len = %v", call, len(callAnalysisList))
				for _, callAnalysis := range callAnalysisList {
					var paramListString string
					for _, param := range callAnalysis.ParamList {
						if len(paramListString) == 0 {
							paramListString = fmt.Sprintf("%v", param)
						} else {
							paramListString = fmt.Sprintf("%v, %v", paramListString, param)
						}
					}
					if len(callAnalysis.From) != 0 {
						utility2.TestOutput("\t- analysis = |%v.%v(%v)|", callAnalysis.From, callAnalysis.Call, paramListString)
					} else {
						utility2.TestOutput("\t- analysis = |%v(%v)|", callAnalysis.Call, paramListString)
					}
					utility2.TestOutput("\t- content  = |%v|", callAnalysis.Content)
				}
			}
		}
		utility2.TestOutput(ui.CommonNote2)
	}

	// 4.1.3.1.6.2

	// utility2.TestOutput(ui.CommonNote2)
	// utility2.TestOutput("analyze go function caller")
	// for packagePath, packageAnalysis := range goAnalysis.PackageAnalysisMap {
	// 	utility2.TestOutput("packagePath = %v", packagePath)
	// 	for functionName, functionAnalysis := range packageAnalysis.FunctionAnalysisMap {
	// 		// Outer Package Call
	// 		for callFromPackagePath, callFunctionMap := range functionAnalysis.OuterPackageCallMap {
	// 			callFromPackageAnalysis, hasCallFromPackageAnalysis := goAnalysis.PackageAnalysisMap[callFromPackagePath]
	// 			if !hasCallFromPackageAnalysis {
	// 				ui.OutputWarnInfo(ui.CMDAnalyzeGoPackageAnalysisNotExist, callFromPackagePath)
	// 				continue
	// 			}
	// 			for callFunction, functionCallAnalysisMap := range callFunctionMap {
	// 				callFunctionAnalysis, hasCallFunctionAnalysis := callFromPackageAnalysis.FunctionAnalysisMap[callFunction]
	// 				if !hasCallFunctionAnalysis {
	// 					ui.OutputWarnInfo(ui.CMDAnalyzeGoFunctionAnalysisNotExist, callFunction)
	// 					continue
	// 				}
	// 				for callIndex, functionCallAnalysis := range functionCallAnalysisMap {
	// 					callFunctionAnalysis.CallerMap[packagePath][functionName][callIndex] = functionCallAnalysis
	// 					utility2.TestOutput("package %v function %v call package %v function %v by params %+v", packagePath, functionName, callFromPackagePath, callFunctionAnalysis.Name, functionCallAnalysis)
	// 				}
	// 			}
	// 		}

	// 		// Member Call
	// 		for callFromPackagePath, callFunctionMap := range functionAnalysis.MemberCallMap {
	// 			callFromPackageAnalysis, hasCallFromPackageAnalysis := goAnalysis.PackageAnalysisMap[callFromPackagePath]
	// 			if !hasCallFromPackageAnalysis {
	// 				ui.OutputWarnInfo(ui.CMDAnalyzeGoPackageAnalysisNotExist, callFromPackagePath)
	// 				continue
	// 			}
	// 			for callFunction, functionCallAnalysisMap := range callFunctionMap {
	// 				callFunctionAnalysis, hasCallFunctionAnalysis := callFromPackageAnalysis.FunctionAnalysisMap[callFunction]
	// 				if !hasCallFunctionAnalysis {
	// 					ui.OutputWarnInfo(ui.CMDAnalyzeGoFunctionAnalysisNotExist, callFunction)
	// 					continue
	// 				}
	// 				for callIndex, functionCallAnalysis := range functionCallAnalysisMap {
	// 					callFunctionAnalysis.CallerMap[packagePath][functionName][callIndex] = functionCallAnalysis
	// 					utility2.TestOutput("package %v function %v call package %v function %v by params %+v", packagePath, functionName, callFromPackagePath, callFunctionAnalysis.Name, functionCallAnalysis)
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// utility2.TestOutput(ui.CommonNote2)

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

func removeGoFileCommentLine(fileContentByte []byte) []byte {
	if goCommentLineRegexp, hasGoCommentLineRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEGoCommentLine]; hasGoCommentLineRegexp {
		result := goCommentLineRegexp.ReplaceAll(fileContentByte, []byte(""))
		return result
	}
	return fileContentByte
}

// analyzeGoScopePackage
// @param content 待分析 package 域的内容
// @param fileABS 待分析文件的绝对路径
// @return
func analyzeGoScopePackage(content, fileABS string) (packageName string, packagePath string, analyzeError error) {
	for _, subMatchList := range regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopePackageTemplate).FindAllStringSubmatch(content, -1) {
		if index, hasIndex := goAnalyzerPackageSubMatchNameIndexMap["NAME"]; hasIndex {
			packageName = strings.TrimSpace(subMatchList[index])
		}
	}

	packagePath, analyzeError = filepath.Rel(global.GoPathSrc, fileABS)
	if packageName != "main" {
		packagePath = path.Join(packagePath, packageName)
	}

	return packageName, packagePath, analyzeError
}

// analyzeGoScopeSingleLineImport
// @param content 待分析单行 import 域的内容
// @return
func analyzeGoScopeSingleLineImport(content string) *GoImportAnalysis {
	return analyzeGoScopeImportContent(strings.Replace(content, "import", "", 1))
}

// analyzeGoScopeMultiLineImport
// @param content 待分析多行 import 域的内容
// @return
func analyzeGoScopeMultiLineImport(content string) []*GoImportAnalysis {
	goImportAnalysisList := make([]*GoImportAnalysis, 0)
	for _, lineContent := range strings.Split(content, "\n") {
		if importAnalysis := analyzeGoScopeImportContent(lineContent); importAnalysis != nil {
			goImportAnalysisList = append(goImportAnalysisList, importAnalysis)
		}
	}
	return goImportAnalysisList
}

// analyzeGoScopeImport
// @param content 待分析 import 域的内容
// @return
func analyzeGoScopeImportContent(content string) *GoImportAnalysis {
	goImportAnalysis := &GoImportAnalysis{}
	for _, subMatchList := range regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopeImportContentTemplate).FindAllStringSubmatch(content, -1) {
		if index, hasIndex := goAnalyzerImportContentSubMatchNameIndexMap["ALIAS"]; hasIndex {
			goImportAnalysis.Alias = strings.TrimSpace(subMatchList[index])
		}
		if index, hasIndex := goAnalyzerImportContentSubMatchNameIndexMap["CONTENT"]; hasIndex {
			goImportAnalysis.Path = strings.TrimSpace(subMatchList[index])
		}
	}
	if len(goImportAnalysis.Path) == 0 {
		return nil
	}
	if len(goImportAnalysis.Alias) == 0 {
		goImportAnalysis.Alias = path.Base(goImportAnalysis.Path)
	}
	return goImportAnalysis
}

// analyzeGoScopePackageVariable
// @param content 待分析 variable 域的内容
// @return
func analyzeGoScopePackageVariable(content string) *GoVariableAnalysis {
	goPackageVariableAnalysis := &GoVariableAnalysis{}
	for _, subMatchList := range regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopePackageVariableTemplate).FindAllStringSubmatch(content, -1) {
		if index, hasIndex := goAnalyzerPackageVariableSubMatchNameIndexMap["NAME"]; hasIndex {
			goPackageVariableAnalysis.Name = strings.TrimSpace(subMatchList[index])
		}
		if index, hasIndex := goAnalyzerPackageVariableSubMatchNameIndexMap["TYPE"]; hasIndex {
			goPackageVariableAnalysis.Type, goPackageVariableAnalysis.TypeFrom = analyzeGoVariableType(strings.TrimSpace(subMatchList[index]))
		}
	}
	if len(goPackageVariableAnalysis.Name) == 0 || len(goPackageVariableAnalysis.Type) == 0 {
		return nil
	}
	return goPackageVariableAnalysis
}

// analyzeGoScopeInterface
// @param content 待分析 interface 域的内容
// @return
func analyzeGoScopeInterface(content string) *GoInterfaceAnalysis {
	goInterfaceAnalysis := &GoInterfaceAnalysis{
		Function: make(map[string]*GoFunctionAnalysis),
	}

	// utility2.TestOutput("interface content = |%v|", content)

	// var interfaceName string
	var interfaceBody string
	for _, subMatchList := range regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopeInterfaceTemplate).FindAllStringSubmatch(content, -1) {
		// if index, hasIndex := goAnalyzerInterfaceSubMatchNameIndexMap["NAME"]; hasIndex {
		// 	interfaceName = strings.TrimSpace(subMatchList[index])
		// }
		if index, hasIndex := goAnalyzerInterfaceSubMatchNameIndexMap["BODY"]; hasIndex {
			interfaceBody = strings.TrimSpace(subMatchList[index])
		}
	}

	// utility2.TestOutput("interfaceName = |%v|", interfaceName)
	// utility2.TestOutput("interfaceBody = |%v|", interfaceBody)

	if len(interfaceBody) == 0 {
		// utility2.TestOutput("interface %v is empty interface", interfaceName)
		return goInterfaceAnalysis
	}

	rootNode := utility3.TraitMultiPunctuationMarksContent(interfaceBody, global.GoAnalyzerScopePunctuationMarkList, 1)
	nodeList := []*utility.NewPunctuationContent{rootNode}
	for len(nodeList) != 0 {
		node := nodeList[0]
		if node == nil {
			break
		}
		nodeList = nodeList[1:]
		// utility2.TestOutput("content = |%v%v%v|", string(node.LeftPunctuationMark.PunctuationMark), node.Content, string(node.RightPunctuationMark.PunctuationMark))
		nodeList = append(nodeList, node.SubPunctuationContentList...)
	}

	interfaceScopeNode := utility3.RecursiveSplitUnderSameDeepPunctuationMarksContentNode(rootNode, global.GoSplitterStringEnter)
	if interfaceScopeNode == nil {
		// utility2.TestOutput("nil interface scope node")
		return goInterfaceAnalysis
	}

	// if len(interfaceScopeNode.ContentList) == 0 {
	// 	utility2.TestOutput("interfaceScopeNode = %+v", interfaceScopeNode)
	// }

	for _, functionDefinitionString := range interfaceScopeNode.ContentList {
		// utility2.TestOutput("index %v, func = |%v|", index, functionDefinitionString)
		replacedContent, _ := utility.ReplaceToUniqueString(functionDefinitionString, global.GoKeywordEmptyInterface)
		// utility2.TestOutput("replace %v to %v", global.GoKeywordEmptyInterface, replacedString)
		// utility2.TestOutput("replaced content = |%v|", replacedContent)

		goFunctionAnalysis := analyzeGoFunctionDefinition(replacedContent)
		goInterfaceAnalysis.Function[goFunctionAnalysis.Name] = goFunctionAnalysis
	}
	// utility2.TestOutput(ui.CommonNote2)

	return goInterfaceAnalysis
}

// analyzeGoScopeStruct
// @param structScope 待分析的 struct 域
// @return
func analyzeGoScopeStruct(structScope *scope) *GoStructAnalysis {
	goStructAnalysis := &GoStructAnalysis{
		Base:           make(map[string]*GoVariableAnalysis),
		MemberVariable: make(map[string]*GoVariableAnalysis),
		MemberFunction: make(map[string]*GoFunctionAnalysis),
	}

	if structScope.isOneLineScope() {
		rootNode := utility3.TraitPunctuationMarksContent(structScope.Content, global.PunctuationMarkCurlyBracket)
		if rootNode == nil || len(rootNode.SubPunctuationContentList) == 0 {
			return nil
		}
		structContent := strings.TrimSpace(rootNode.SubPunctuationContentList[0].Content)
		goVariableAnalysis := &GoVariableAnalysis{}
		if spaceIndex := strings.IndexRune(structContent, ' '); spaceIndex == -1 {
			goVariableAnalysis.Type, goVariableAnalysis.TypeFrom = analyzeGoVariableType(structContent)
			goVariableAnalysis.Name = goVariableAnalysis.Type
			goStructAnalysis.Base[goVariableAnalysis.Name] = goVariableAnalysis
		} else {
			goVariableAnalysis.Name = structContent[0:spaceIndex]
			goVariableAnalysis.Type, goVariableAnalysis.TypeFrom = analyzeGoVariableType(structContent[spaceIndex+1:])
			goStructAnalysis.MemberVariable[goVariableAnalysis.Name] = goVariableAnalysis
		}
	} else {
		lineContentList := strings.Split(structScope.Content, "\n")
		for lineIndex, lineContent := range lineContentList {
			// continue scope start and scope end
			if lineIndex == 0 || lineIndex == len(lineContentList)-1 {
				continue
			}
			goVariableAnalysis := &GoVariableAnalysis{}
			for _, subMatchList := range regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopeStructVariableTemplate).FindAllStringSubmatch(lineContent, -1) {
				if index, hasIndex := goAnalyzerStructVariableSubMatchNameIndexMap["NAME"]; hasIndex {
					goVariableAnalysis.Name = strings.TrimSpace(subMatchList[index])
				}
				if index, hasIndex := goAnalyzerStructVariableSubMatchNameIndexMap["TYPE"]; hasIndex {
					goVariableAnalysis.Type, goVariableAnalysis.TypeFrom = analyzeGoVariableType(strings.TrimSpace(subMatchList[index]))
				}
			}
			if len(goVariableAnalysis.Name) == 0 {
				continue
			}
			if len(goVariableAnalysis.Type) == 0 && len(goVariableAnalysis.TypeFrom) == 0 {
				goVariableAnalysis.Type, goVariableAnalysis.TypeFrom = analyzeGoVariableType(strings.TrimSpace(goVariableAnalysis.Name))
				goVariableAnalysis.Name = goVariableAnalysis.Type
				goStructAnalysis.Base[goVariableAnalysis.Name] = goVariableAnalysis
			} else {
				goStructAnalysis.MemberVariable[goVariableAnalysis.Name] = goVariableAnalysis
			}
		}
	}

	return goStructAnalysis
}

// analyzeGoScopeFunction
// @param functionScope 待分析的 function 域
// @return
func analyzeGoScopeFunction(functionScope *scope) *GoFunctionAnalysis {
	// goFunctionAnalysis := &GoFunctionAnalysis{}

	// rootNode := utility3.TraitPunctuationMarksContent(functionScope.Content,  global.PunctuationMarkBracket)

	// utility2.TestOutput(ui.CommonNote2)
	// utility2.TestOutput("functionScope.Content = \n|%v|", functionScope.Content)
	replacedContent, _ := utility.ReplaceToUniqueString(functionScope.Content, global.GoKeywordEmptyInterface)
	// utility2.TestOutput("replace %v to %v", global.GoKeywordEmptyInterface, replacedString)
	// utility2.TestOutput("replaced content = %v", replacedContent)

	contentRootNode := utility3.TraitMultiPunctuationMarksContent(replacedContent, global.GoAnalyzerScopePunctuationMarkList, 1)
	subNodeCount := len(contentRootNode.SubPunctuationContentList)
	// utility2.TestOutput("subNodeCount = %v", subNodeCount)

	// utility2.TestOutput("contentRootNode.Content = |%v|", contentRootNode.Content)
	// utility2.TestOutput("contentRootNode.Left = %+v", contentRootNode.LeftPunctuationMark)
	// utility2.TestOutput("contentRootNode.Right = %+v", contentRootNode.RightPunctuationMark)
	// utility2.TestOutput("contentRootNode.SubPunctuationContentList[0].Content = |%v|", contentRootNode.SubPunctuationContentList[0].Content)
	// utility2.TestOutput("contentRootNode.SubPunctuationContentList[0].LeftPunctuationMark = %+v", contentRootNode.SubPunctuationContentList[0].LeftPunctuationMark)
	// utility2.TestOutput("contentRootNode.SubPunctuationContentList[0].RightPunctuationMark = %+v", contentRootNode.SubPunctuationContentList[0].RightPunctuationMark)

	// utility2.TestOutput("contentRootNode.SubPunctuationContentList[1].Content = |%v|", contentRootNode.SubPunctuationContentList[1].Content)
	// utility2.TestOutput("contentRootNode.SubPunctuationContentList[1].LeftPunctuationMark = %+v", contentRootNode.SubPunctuationContentList[1].LeftPunctuationMark)
	// utility2.TestOutput("contentRootNode.SubPunctuationContentList[1].RightPunctuationMark = %+v", contentRootNode.SubPunctuationContentList[1].RightPunctuationMark)

	if contentRootNode == nil || subNodeCount < 2 {
		return nil
	}

	// var functionParamListString string
	// var functionReturnListString string
	// var functionBodyString string
	// var functionReturnTypeString string

	// functionParamListString = contentRootNode.SubPunctuationContentList[0].Content
	// utility2.TestOutput("functionParamListString = |%v|", functionParamListString)
	// functionParamListNode := utility3.RecursiveSplitUnderSameDeepPunctuationMarksContent(functionParamListString, global.GoAnalyzerScopePunctuationMarkList, global.GoSplitterStringComma)
	// utility2.TestOutput("param list:")
	// for index, param := range functionParamListNode.ContentList {
	// 	utility2.TestOutput("index %v param = |%v|", index, param)
	// }

	// utility2.TestOutput("function definition string = |%v|", replacedContent[:contentRootNode.SubPunctuationIndexMap[subNodeCount-1].Left])
	// utility2.TestOutput("subNodeCount = %v", subNodeCount)
	// utility2.TestOutput("contentRootNode.SubPunctuationContentList[subNodeCount-1] = %+v", contentRootNode.SubPunctuationContentList[subNodeCount-1])
	// utility2.TestOutput("function definition string = |%v|", replacedContent[:contentRootNode.SubPunctuationContentList[subNodeCount-1].LeftPunctuationMark.Index])
	// goFunctionAnalysis := analyzeGoFunctionDefinition(replacedContent[:contentRootNode.SubPunctuationIndexMap[subNodeCount-1].Left])
	goFunctionAnalysis := analyzeGoFunctionDefinition(replacedContent[:contentRootNode.SubPunctuationContentList[subNodeCount-1].LeftPunctuationMark.Index])
	goFunctionAnalysis.CallMap = analyzeGoFunctionBody(contentRootNode.SubPunctuationContentList[subNodeCount-1].Content)

	// utility2.TestOutput(ui.CommonNote2)

	return goFunctionAnalysis
}

// analyzeGoFunctionDefinition
// @param functionDefinitionContent 待分析的函数定义的内容
// @return
func analyzeGoFunctionDefinition(functionDefinitionContent string) *GoFunctionAnalysis {
	functionDefinitionContentRootNode := utility3.TraitMultiPunctuationMarksContent(functionDefinitionContent, global.GoAnalyzerScopePunctuationMarkList, 1)
	subNodeCount := len(functionDefinitionContentRootNode.SubPunctuationContentList)

	if functionDefinitionContentRootNode == nil || subNodeCount < 1 {
		return nil
	}

	// function name
	name := strings.TrimSpace(functionDefinitionContent[:functionDefinitionContentRootNode.SubPunctuationContentList[0].LeftPunctuationMark.Index])
	if len(name) == 0 {
		return nil
	}

	// param list
	paramListString := functionDefinitionContentRootNode.SubPunctuationContentList[0].Content
	paramListNode := utility3.RecursiveSplitUnderSameDeepPunctuationMarksContent(paramListString, global.GoAnalyzerScopePunctuationMarkList, global.GoSplitterStringComma)
	// utility2.TestOutput("param list:")
	// for index, param := range paramListNode.ContentList {
	// 	utility2.TestOutput("index %v param = |%v|", index, param)
	// }
	paramMap := analyzeGoFunctionDefinitionParamList(paramListNode.ContentList)
	// for paramName, analysis := range paramMap {
	// 	utility2.TestOutput("paramName = %v, index = %v, type = |%v|, type from = %v", paramName, analysis.Index, analysis.Type, analysis.TypeFrom)
	// }

	// return list
	returnListString := strings.TrimSpace(functionDefinitionContent[functionDefinitionContentRootNode.SubPunctuationContentList[0].RightPunctuationMark.Index+1:])
	returnListString = strings.TrimFunc(returnListString, func(r rune) bool {
		return r == global.PunctuationMarkLeftBracket || r == global.PunctuationMarkRightBracket
	})
	// utility2.TestOutput("return list = |%v|", returnListString)
	returnListNode := utility3.RecursiveSplitUnderSameDeepPunctuationMarksContent(returnListString, global.GoAnalyzerScopePunctuationMarkList, global.GoSplitterStringComma)
	// utility2.TestOutput("return list:")
	// for index, returnType := range returnListNode.ContentList {
	// 	utility2.TestOutput("index %v return = |%v|", index, returnType)
	// }
	returnMap := analyzeGoFunctionDefinitionReturnList(returnListNode.ContentList)
	// for returnName, analysis := range returnMap {
	// 	utility2.TestOutput("returnName = %v, index = %v, type = |%v|, type from = %v", returnName, analysis.Index, analysis.Type, analysis.TypeFrom)
	// }

	return &GoFunctionAnalysis{
		Name:      name,
		ParamsMap: paramMap,
		ReturnMap: returnMap,
	}
}

// analyzeGoFunctionDefinitionParamList
// @param paramListString 待分析的函数参数表
// @return
func analyzeGoFunctionDefinitionParamList(paramStringList []string) map[string]*GoFunctionVariable {
	paramMap := make(map[string]*GoFunctionVariable)
	unknownTypeParamList := make([]*GoFunctionVariable, 0)
	for index, paramString := range paramStringList {
		paramStringTimSpace := strings.TrimSpace(paramString)
		if len(paramStringTimSpace) == 0 {
			continue
		}
		splitterIndex := strings.Index(paramStringTimSpace, global.GoSplitterStringSpace)
		if splitterIndex == -1 {
			unknownTypeParamList = append(unknownTypeParamList, &GoFunctionVariable{
				GoVariableAnalysis: GoVariableAnalysis{
					Name: paramStringTimSpace,
				},
				Index: index,
			})
		} else {
			paramType, paramTypeFrom := analyzeGoVariableType(paramStringTimSpace[splitterIndex+1:])
			paramMap[paramStringTimSpace[:splitterIndex]] = &GoFunctionVariable{
				GoVariableAnalysis: GoVariableAnalysis{
					Name:     paramStringTimSpace[:splitterIndex],
					Type:     paramType,
					TypeFrom: paramTypeFrom,
				},
				Index: index,
			}
			if len(unknownTypeParamList) != 0 {
				for _, unknownTypeParam := range unknownTypeParamList {
					unknownTypeParam.GoVariableAnalysis.Type = paramType
					unknownTypeParam.GoVariableAnalysis.TypeFrom = paramTypeFrom
					paramMap[unknownTypeParam.Name] = unknownTypeParam
				}
				unknownTypeParamList = make([]*GoFunctionVariable, 0)
			}
		}
	}
	return paramMap
}

// analyzeGoFunctionDefinitionReturnList
// @param returnListString 待分析的返回值列表
// @return
func analyzeGoFunctionDefinitionReturnList(returnStringList []string) map[string]*GoFunctionVariable {
	returnMap := make(map[string]*GoFunctionVariable)

	for index, returnString := range returnStringList {
		returnStringTimSpace := strings.TrimSpace(returnString)
		if len(returnStringTimSpace) == 0 {
			continue
		}
		splitterIndex := strings.Index(returnStringTimSpace, global.GoSplitterStringSpace)
		if splitterIndex == -1 {
			returnType, returnTypeFrom := analyzeGoVariableType(returnStringTimSpace)
			returnMap[fmt.Sprintf("%v", len(returnMap))] = &GoFunctionVariable{
				GoVariableAnalysis: GoVariableAnalysis{
					Type:     returnType,
					TypeFrom: returnTypeFrom,
				},
				Index: index,
			}
		} else {
			// TODO: return string = |v int|
			// TODO: return string = |func(int, int) int|
			// TODO: return string = |f func(int, int) int|
			// TODO: return string = |struct{ int }|
			// TODO: return string = |s struct{ int }|
			// TODO: return string = |s struct {
			// 	v int
			// 	f func()
			// })|

			// TEMP:
			keywordIndex := strings.Index(returnStringTimSpace, global.GoKeywordFunc)
			if keywordIndex == -1 {
				keywordIndex = strings.Index(returnStringTimSpace, global.GoKeywordStruct)
			}
			var returnMapKey string
			var returnName string
			var returnType string
			var returnTypeFrom string
			if keywordIndex != -1 {
				returnName := strings.TrimSpace(returnStringTimSpace[:keywordIndex])
				returnType, returnTypeFrom = analyzeGoVariableType(returnStringTimSpace[keywordIndex:])
				if len(returnName) == 0 {
					returnMapKey = fmt.Sprintf("%v", len(returnMap))
				} else {
					returnMapKey = returnName
				}
			} else {
				returnType, returnTypeFrom = analyzeGoVariableType(returnStringTimSpace)
				returnMapKey = fmt.Sprintf("%v", len(returnMap))
			}

			returnMap[returnMapKey] = &GoFunctionVariable{
				GoVariableAnalysis: GoVariableAnalysis{
					Name:     returnName,
					Type:     returnType,
					TypeFrom: returnTypeFrom,
				},
				Index: index,
			}
		}
	}
	return returnMap
}

// analyzeGoFunctionBody
// @functionBodyContent 待分析的函数体的内容
// @return
func analyzeGoFunctionBody(functionBodyContent string) map[string][]*GoFunctionCallAnalysis {
	callMap := make(map[string][]*GoFunctionCallAnalysis)
	bodyContentRootNode := utility3.TraitMultiPunctuationMarksContent(functionBodyContent, global.GoAnalyzerScopePunctuationMarkList, -1)
	toSearchNodeList := []*utility.NewPunctuationContent{bodyContentRootNode}
	bracketScopeNodeList := make([]*utility.NewPunctuationContent, 0)
	// bracketScopeContentList := make([]string, 0)

	for len(toSearchNodeList) != 0 {
		toSearchNode := toSearchNodeList[0]
		if toSearchNode == nil {
			continue
		}

		// fmt.Printf("content |%v%v%v|, index testContent[%v:%v]\n", string(toSearchNode.LeftPunctuationMark.PunctuationMark), toSearchNode.Content, string(toSearchNode.RightPunctuationMark.PunctuationMark), toSearchNode.LeftPunctuationMark.Index, toSearchNode.RightPunctuationMark.Index)
		// if toSearchNode.LeftPunctuationMark.Index >= 0 && toSearchNode.RightPunctuationMark.Index < len(toSearchNode.Content) {
		// 	fmt.Printf("content by index = |%v|\n", functionBodyContent[toSearchNode.LeftPunctuationMark.Index+1:toSearchNode.RightPunctuationMark.Index])
		// }

		if toSearchNode.LeftPunctuationMark.PunctuationMark == global.PunctuationMarkLeftBracket {
			bracketScopeNodeList = append(bracketScopeNodeList, toSearchNode)
		}
		toSearchNodeList = toSearchNodeList[1:]
		toSearchNodeList = append(toSearchNodeList, toSearchNode.SubPunctuationContentList...)
	}

	for _, bracketScopeNode := range bracketScopeNodeList {
		// utility2.TestOutput("bracketScopeNode content = |%v|", bracketScopeNode.Content)
		// utility2.TestOutput("bracketScopeNode left index = |%v|", bracketScopeNode.LeftPunctuationMark.Index)
		searchIndex := bracketScopeNode.LeftPunctuationMark.Index
		for searchIndex != -1 {
			searchIndex--
			// utility2.TestOutput("functionBodyContent[%v] = %v", searchIndex, string(rune(functionBodyContent[searchIndex])))
			if isIdentifierRune(rune(functionBodyContent[searchIndex])) {
				continue
			}
			break
		}
		// utility2.TestOutput("searchIndex = %v", searchIndex)
		// utility2.TestOutput("bracketScopeNode.LeftPunctuationMark.Index = %v", bracketScopeNode.LeftPunctuationMark.Index)
		identifierString := functionBodyContent[searchIndex+1 : bracketScopeNode.LeftPunctuationMark.Index]
		// utility2.TestOutput("identifierString = |%v|", identifierString)
		callIdentifier := identifierString
		// utility2.TestOutput("call |%v(%v)|", callIdentifier, bracketScopeNode.Content)
		var callIdentifierFrom string
		punctuationMarkPointIndex := strings.Index(identifierString, global.GoSplitterStringPoint)
		if punctuationMarkPointIndex != -1 {
			callIdentifierFrom = identifierString[:punctuationMarkPointIndex]
			callIdentifier = identifierString[punctuationMarkPointIndex+1:]
		}
		// if len(callIdentifierFrom) != 0 {
		// 	utility2.TestOutput("call |%v.%v(%v)|", callIdentifierFrom, callIdentifier, bracketScopeNode.Content)
		// } else {
		// 	utility2.TestOutput("call |%v(%v)|", callIdentifier, bracketScopeNode.Content)
		// }
		if _, hasCall := callMap[identifierString]; !hasCall {
			// utility2.TestOutput("|%v| does not have list, call (%v)", identifierString, bracketScopeNode.Content)
			callMap[identifierString] = make([]*GoFunctionCallAnalysis, 0)
		}

		// for _, subNode := range bracketScopeNode.SubPunctuationContentList {
		// 	utility2.TestOutput("subNode Content = |%v|", subNode.Content)
		// }

		paramListRootNode := utility3.RecursiveSplitUnderSameDeepPunctuationMarksContent(bracketScopeNode.Content, global.GoAnalyzerScopePunctuationMarkList, global.GoSplitterStringComma)
		paramList := make([]string, 0, len(paramListRootNode.ContentList))
		for _, param := range paramListRootNode.ContentList {
			if len(strings.TrimSpace(param)) != 0 {
				paramList = append(paramList, strings.TrimSpace(param))
				// utility2.TestOutput("index %v param = %v", index, strings.TrimSpace(param))
			}
		}

		callMap[identifierString] = append(callMap[identifierString], &GoFunctionCallAnalysis{
			Content:   fmt.Sprintf("%v(%v)", identifierString, bracketScopeNode.Content),
			From:      callIdentifierFrom,
			Call:      callIdentifier,
			ParamList: paramList,
		})
	}

	return callMap
}

// isIdentifierRune
// @param r 待检查的字符
// @return
func isIdentifierRune(r rune) bool {
	return r == 46 || r == 95 || (48 <= r && r <= 57) || (65 <= r && r <= 90) || (97 <= r && r <= 123)
}

// analyzeGoVariableType
// @param variableTypeString 待分析变量的类型字符串
// @return
func analyzeGoVariableType(variableTypeString string) (string, string) {
	var typeName string
	var typeFrom string
	keywordIndex := strings.Index(variableTypeString, global.GoKeywordFunc)
	if keywordIndex == -1 {
		keywordIndex = strings.Index(variableTypeString, global.GoKeywordStruct)
	}
	if keywordIndex == -1 {
		punctuationMarkPointIndex := strings.Index(variableTypeString, global.GoSplitterStringPoint)
		if punctuationMarkPointIndex == -1 {
			typeName = variableTypeString
		} else {
			typeFrom = variableTypeString[:punctuationMarkPointIndex]
			typeName = variableTypeString[punctuationMarkPointIndex+1:]
		}
	} else {
		typeName = variableTypeString
	}
	return typeName, typeFrom
}

// func outputAnalyzeGoFileResult(goFileAnalysis *goFileAnalysis) string {
// 	templateStyleRegexp, hasTemplateStyleRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AETemplateStyle]
// 	if !hasTemplateStyleRegexp {
// 		ui.OutputWarnInfo(ui.CommonError16, global.AETemplateStyle)
// 		return ""
// 	}

// 	// file path
// 	resultContent := strings.Replace(ui.AnalyzeGoFileResultTemplate, global.AnalyzeRPFilePath, goFileAnalysis.FilePath, -1)

// 	// package content
// 	packageContent := ui.ParseStyleTemplate(templateStyleRegexp, ui.AnalyzeGoFilePackageContentTemplate)
// 	packageContent = strings.Replace(packageContent, global.AnalyzeRPPackageName, goFileAnalysis.PackageName, -1)
// 	packageContent = strings.Replace(packageContent, global.AnalyzeRPPackagePath, goFileAnalysis.PackagePath, -1)
// 	resultContent = strings.Replace(resultContent, global.AnalyzeRPPackageContent, packageContent, -1)

// 	// import 内容
// 	importPackageListContent := ""
// 	if len(goFileAnalysis.ImportAliasMap) != 0 {
// 		importPackageListString := ""
// 		for packageAlias, packagePath := range goFileAnalysis.ImportAliasMap {
// 			packageAliasPathContent := ui.AnalyzeGoFileImportPackageTemplate

// 			// style template
// 			packageAliasPathContent = ui.ParseStyleTemplate(templateStyleRegexp, packageAliasPathContent)

// 			// package alias
// 			packageAliasPathContent = strings.Replace(packageAliasPathContent, global.AnalyzeRPPackageAlias, packageAlias, -1)

// 			// package path
// 			packageAliasPathContent = strings.Replace(packageAliasPathContent, global.AnalyzeRPPackagePath, packagePath, -1)

// 			if importPackageListString == "" {
// 				importPackageListString = packageAliasPathContent
// 			} else {
// 				importPackageListString = fmt.Sprintf("%v\n%v", importPackageListString, packageAliasPathContent)
// 			}
// 		}
// 		importPackageListContent = strings.Replace(ui.AnalyzeGoFileImportPackageListTemplate, global.AnalyzeRPImportPackage, importPackageListString, -1)
// 	}
// 	resultContent = strings.Replace(resultContent, global.AnalyzeRPImportPackageList, importPackageListContent, -1)

// 	// function 内容
// 	functionDefinitionListContent := ""
// 	if len(goFileAnalysis.functionList) != 0 {
// 		functionDefinitionListString := ""
// 		for _, functionName := range goFileAnalysis.functionList {
// 			functionAnalysis := goFileAnalysis.FunctionMap[functionName]
// 			functionDefinitionContent := ui.AnalyzeGoFileFunctionDefinitionTemplate

// 			// sytle template
// 			functionDefinitionContent = ui.ParseStyleTemplate(templateStyleRegexp, functionDefinitionContent)

// 			// function name
// 			functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionName, functionAnalysis.Name, -1)

// 			// function class
// 			if functionAnalysis.Class != "" {
// 				functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionClass, strings.Replace(ui.AnalyzeGoFileFunctionClassTemplate, global.AnalyzeRPFunctionClassName, functionAnalysis.Class, -1), -1)
// 			} else {
// 				functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionClass, global.AnalyzeRPEmptyString, -1)
// 			}

// 			// function params
// 			functionParamListContent := ""
// 			if len(functionAnalysis.ParamsMap) != 0 {
// 				// functionParamListContent = parseGoFunctionParamOrReturnListContent(templateStyleRegexp, functionAnalysis.ParamsMap, functionParamList)
// 			}
// 			functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionParamList, functionParamListContent, -1)

// 			// function return
// 			functionReturnListContent := ""
// 			if len(functionAnalysis.ReturnMap) != 0 {
// 				// functionReturnListContent = parseGoFunctionParamOrReturnListContent(templateStyleRegexp, functionAnalysis.ReturnMap, functionReturnList)
// 			}
// 			functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionReturnList, functionReturnListContent, -1)

// 			// function body call map
// 			functionCallMapContent := ""
// 			if len(functionAnalysis.OuterPackageCallMap) != 0 {
// 				// functionCallMapContent = parseGoFunctionCallMapContent(templateStyleRegexp, functionAnalysis.OuterPackageCallMap)
// 			}
// 			functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionCallMap, functionCallMapContent, -1)

// 			// add to function definition list content
// 			if functionDefinitionListString == "" {
// 				functionDefinitionListString = functionDefinitionContent
// 			} else {
// 				functionDefinitionListString = fmt.Sprintf("%v\n%v", functionDefinitionListString, functionDefinitionContent)
// 			}
// 		}
// 		functionDefinitionListContent = strings.Replace(ui.AnalyzeGoFileFunctionDefinitionListTemplate, global.AnalyzeRPFunctionDefinition, functionDefinitionListString, -1)
// 	}
// 	resultContent = strings.Replace(resultContent, global.AnalyzeRPFunctionDefinitionList, functionDefinitionListContent, -1)

// 	// clear space line
// 	resultContent = utility3.TrimSpaceLine(resultContent)

// 	return resultContent
// }

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

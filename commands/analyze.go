package commands

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-worker/global"
	"github.com/go-worker/logger"
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
	Scope               map[string]*GoPackageScope              // file path : go package scope
	PackageName         string                                  // package name
	PackagePath         string                                  // package path
	ImportAnalysis      map[string]map[string]*GoImportAnalysis // file path : package path : import analysis
	VariableAnalysisMap map[string]*GoVariableAnalysis          // variable name : variable analysis
	ConstAnalysisMap    map[string]*GoVariableAnalysis          // const name : variable analysis
	InterfaceAnalysis   map[string]*GoInterfaceAnalysis         // interface name : interface analysis
	StructAnalysis      map[string]*GoStructAnalysis            // struct name : struct analysis
	FunctionAnalysisMap map[string]*GoFunctionAnalysis          // function name : function analysis
	TypeRename          map[string]*GoTypeRenameAnalysis        // rename : origin type
}

// GoImportAnalysis go 引入包的分析结果
type GoImportAnalysis struct {
	Alias string
	Path  string
}

type GoTypeAnalysis struct {
	Name            string
	From            string
	FromPackagePath string
}

// GoVariableAnalysis go 变量分析结果
type GoVariableAnalysis struct {
	Name string
	Type *GoTypeAnalysis
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
	VariableThis        *GoVariableAnalysis                                   // class variable this
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
	Content         string
	From            string
	FromPackagePath string
	Call            string
	ParamList       []string
	PreviousCall    *GoFunctionCallAnalysis
	PostCall        *GoFunctionCallAnalysis
}

// GoFunctionVariable go 函数内变量
type GoFunctionVariable struct {
	GoVariableAnalysis
	Index int
}

// GoTypeRenameAnalysis go 类型重命名分析结果
type GoTypeRenameAnalysis struct {
	Name string
	Type *GoTypeAnalysis
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

	goAnalysis := &GoAnalysis{
		RootPath:           rootPath,
		PackageAnalysisMap: make(map[string]*GoPackageAnalysis),
	}

	// 4.1.3.1.6.1
	for _, toAnalyzeFilePath := range toAnalyzePathList {
		// 4.1.3.1.6.1.1
		splitFileResult := SplitGoFile(toAnalyzeFilePath, false)

		// 4.1.3.1.6.1.2
		packageName, packagePath, analyzeGoScopePackageError := analyzeGoScopePackage(splitFileResult.Package.Content, rootPath)
		if analyzeGoScopePackageError != nil {
			return nil, analyzeGoScopePackageError
		}

		if _, hasPackagePath := goAnalysis.PackageAnalysisMap[packagePath]; !hasPackagePath {
			goAnalysis.PackageAnalysisMap[packagePath] = &GoPackageAnalysis{
				PackageName:         packageName,
				PackagePath:         packagePath,
				ImportAnalysis:      make(map[string]map[string]*GoImportAnalysis),
				VariableAnalysisMap: make(map[string]*GoVariableAnalysis),
				ConstAnalysisMap:    make(map[string]*GoVariableAnalysis),
				InterfaceAnalysis:   make(map[string]*GoInterfaceAnalysis),
				StructAnalysis:      make(map[string]*GoStructAnalysis),
				FunctionAnalysisMap: make(map[string]*GoFunctionAnalysis),
				TypeRename:          make(map[string]*GoTypeRenameAnalysis),
				Scope:               make(map[string]*GoPackageScope),
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

		fileImportAnalysisMap := goAnalysis.PackageAnalysisMap[packagePath].ImportAnalysis[toAnalyzeFilePath]

		// 4.1.3.1.6.1.4
		for _, packageVariableScope := range splitFileResult.PackageVariable {
			if goPackageVariableAnalysis := analyzeGoScopePackageVariable(packageVariableScope.Content, fileImportAnalysisMap); goPackageVariableAnalysis != nil {
				goAnalysis.PackageAnalysisMap[packagePath].VariableAnalysisMap[goPackageVariableAnalysis.Name] = goPackageVariableAnalysis
			}
		}

		// 4.1.3.1.6.1.5
		for interfaceName, interfaceScope := range splitFileResult.InterfaceDefinition {
			if goInterfaceAnalysis := analyzeGoScopeInterface(interfaceScope.Content, fileImportAnalysisMap); goInterfaceAnalysis != nil {
				goInterfaceAnalysis.Name = interfaceName
				goAnalysis.PackageAnalysisMap[packagePath].InterfaceAnalysis[interfaceName] = goInterfaceAnalysis
			}
		}

		// 4.1.3.1.6.1.6
		for structName, structScope := range splitFileResult.StructDefinition {
			if goStructAnalysis := analyzeGoScopeStruct(structScope, fileImportAnalysisMap); goStructAnalysis != nil {
				goStructAnalysis.Name = structName
				if _, hasAnalysis := goAnalysis.PackageAnalysisMap[packagePath].StructAnalysis[structName]; !hasAnalysis {
					goAnalysis.PackageAnalysisMap[packagePath].StructAnalysis[structName] = goStructAnalysis
				} else {
					goAnalysis.PackageAnalysisMap[packagePath].StructAnalysis[structName].Name = goStructAnalysis.Name
					goAnalysis.PackageAnalysisMap[packagePath].StructAnalysis[structName].Base = goStructAnalysis.Base
					goAnalysis.PackageAnalysisMap[packagePath].StructAnalysis[structName].MemberVariable = goStructAnalysis.MemberVariable
				}
			}
		}

		// 4.1.3.1.6.1.7
		for functionName, functionScope := range splitFileResult.FunctionDefinition {
			if functionAnalysis := analyzeGoScopeFunction(functionScope, false, fileImportAnalysisMap); functionAnalysis != nil {
				goAnalysis.PackageAnalysisMap[packagePath].FunctionAnalysisMap[functionName] = functionAnalysis
			}
		}

		// 4.1.3.1.6.1.8
		for className, functionMap := range splitFileResult.MemberFunctionDefinition {
			for functionName, functionScope := range functionMap {
				if functionAnalysis := analyzeGoScopeFunction(functionScope, true, fileImportAnalysisMap); functionAnalysis != nil {
					if _, hasAnalysis := goAnalysis.PackageAnalysisMap[packagePath].StructAnalysis[className]; !hasAnalysis {
						goAnalysis.PackageAnalysisMap[packagePath].StructAnalysis[className] = &GoStructAnalysis{
							Name:           className,
							MemberFunction: make(map[string]*GoFunctionAnalysis),
						}
					}
					goAnalysis.PackageAnalysisMap[packagePath].StructAnalysis[className].MemberFunction[functionName] = functionAnalysis
				}
			}
		}

		// 4.1.3.1.6.1.9
		for variableName, constScope := range splitFileResult.SingleLineConst {
			for _, constVariableAnalysis := range analyzeGoScopeConst(constScope, fileImportAnalysisMap) {
				goAnalysis.PackageAnalysisMap[packagePath].ConstAnalysisMap[variableName] = constVariableAnalysis
			}
		}
		for _, constScope := range splitFileResult.MultiLineConst {
			for variableName, constVariableAnalysis := range analyzeGoScopeConst(constScope, fileImportAnalysisMap) {
				goAnalysis.PackageAnalysisMap[packagePath].ConstAnalysisMap[variableName] = constVariableAnalysis
			}
		}

		// 4.1.3.1.6.1.10
		for _, renameMap := range splitFileResult.TypeRename {
			for rename, renameScope := range renameMap {
				if typeRenameAnalysis := analyzeGoTypeRename(renameScope, fileImportAnalysisMap); typeRenameAnalysis != nil {
					goAnalysis.PackageAnalysisMap[packagePath].TypeRename[rename] = typeRenameAnalysis
				}
			}
		}
	}

	for packagePath, packageAnalysis := range goAnalysis.PackageAnalysisMap {
		logger.OutputNoteInfo(ui.CommonNote2)
		logger.OutputNoteInfo("packageName = %v, packagePath = %v", packageAnalysis.PackageName, packagePath)

		logger.OutputNoteInfo(ui.CommonNote2)
		logger.OutputNoteInfo("Output package import analysis:")
		for filePath, importAnlysisMap := range goAnalysis.PackageAnalysisMap[packagePath].ImportAnalysis {
			logger.OutputNoteInfo("import list from file: %v", filePath)
			for alias, analysis := range importAnlysisMap {
				logger.OutputNoteInfo("alias = %v, importPath = %v", alias, analysis.Path)
			}
		}
		logger.OutputNoteInfo(ui.CommonNote2)

		logger.OutputNoteInfo(ui.CommonNote2)
		logger.OutputNoteInfo("Output package variable analysis:")
		for name, analysis := range goAnalysis.PackageAnalysisMap[packagePath].VariableAnalysisMap {
			logger.OutputNoteInfo("name = %v, type = %v, type from = |%v %v|", name, analysis.Type.Name, analysis.Type.From, analysis.Type.FromPackagePath)
		}
		logger.OutputNoteInfo(ui.CommonNote2)

		logger.OutputNoteInfo(ui.CommonNote2)
		logger.OutputNoteInfo("Output package interface analysis:")
		for interfaceName, interfaceAnalysis := range goAnalysis.PackageAnalysisMap[packagePath].InterfaceAnalysis {
			logger.OutputNoteInfo("interface %v function list: %v", interfaceName, len(interfaceAnalysis.Function))
			for functionName, functionAnalysis := range interfaceAnalysis.Function {
				logger.OutputNoteInfo("function: %v", functionName)
				logger.OutputNoteInfo("- param list:")
				for paramName, paramAnalysis := range functionAnalysis.ParamsMap {
					var paramTypeString string
					if len(paramAnalysis.Type.From) != 0 {
						paramTypeString = fmt.Sprintf("(%v).%v", paramAnalysis.Type.FromPackagePath, paramAnalysis.Type.Name)
					} else {
						paramTypeString = paramAnalysis.Type.Name
					}
					logger.OutputNoteInfo("\t- Index: %v, Name: %v, Type: %v", paramAnalysis.Index, paramName, paramTypeString)
				}
				logger.OutputNoteInfo("- return list:")
				for returnIndex, returnAnalysis := range functionAnalysis.ReturnMap {
					var returnTypeString string
					if len(returnAnalysis.Type.From) != 0 {
						returnTypeString = fmt.Sprintf("(%v).%v", returnAnalysis.Type.FromPackagePath, returnAnalysis.Type)
					} else {
						returnTypeString = returnAnalysis.Type.Name
					}
					logger.OutputNoteInfo("\t- Index: %v, Name: %v, Type: %v", returnIndex, returnAnalysis.Name, returnTypeString)
				}
			}
		}
		logger.OutputNoteInfo(ui.CommonNote2)

		logger.OutputNoteInfo(ui.CommonNote2)
		logger.OutputNoteInfo("Output package struct analysis:")
		for structName, structAnalysis := range goAnalysis.PackageAnalysisMap[packagePath].StructAnalysis {
			logger.OutputNoteInfo("struct %v base struct list: %v", structName, len(structAnalysis.Base))
			for baseType, baseVariableAnalysis := range structAnalysis.Base {
				logger.OutputNoteInfo("- base type: %v, from |%v %v|", baseType, baseVariableAnalysis.Type.From, baseVariableAnalysis.Type.FromPackagePath)
			}
			logger.OutputNoteInfo("struct %v variable list: %v", structName, len(structAnalysis.MemberVariable))
			for variableName, variableAnalysis := range structAnalysis.MemberVariable {
				logger.OutputNoteInfo("- variable name: %v, type: %v, from: |%v %v|", variableName, variableAnalysis.Type.Name, variableAnalysis.Type.From, variableAnalysis.Type.FromPackagePath)
			}
		}
		logger.OutputNoteInfo(ui.CommonNote2)

		logger.OutputNoteInfo(ui.CommonNote2)
		logger.OutputNoteInfo("Output package function analysis:")
		for functionName, functionAnalysis := range goAnalysis.PackageAnalysisMap[packagePath].FunctionAnalysisMap {
			logger.OutputNoteInfo("function %v", functionName)
			for paramName, functionVariableAnalysis := range functionAnalysis.ParamsMap {
				logger.OutputNoteInfo("- param index: %v, name: %v, type: %v, from: |%v %v|", functionVariableAnalysis.Index, paramName, functionVariableAnalysis.Type.Name, functionVariableAnalysis.Type.From, functionVariableAnalysis.Type.FromPackagePath)
			}
			for _, functionVariableAnalysis := range functionAnalysis.ReturnMap {
				logger.OutputNoteInfo("- return index: %v, name: %v, type: %v, from: |%v %v|", functionVariableAnalysis.Index, functionVariableAnalysis.Name, functionVariableAnalysis.Type.Name, functionVariableAnalysis.Type.From, functionVariableAnalysis.Type.FromPackagePath)
			}
			for call, callAnalysisList := range functionAnalysis.CallMap {
				logger.OutputNoteInfo("- call: %v, len = %v", call, len(callAnalysisList))
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
						logger.OutputNoteInfo("\t- analysis = |%v.%v(%v)|", callAnalysis.From, callAnalysis.Call, paramListString)
					} else {
						logger.OutputNoteInfo("\t- analysis = |%v(%v)|", callAnalysis.Call, paramListString)
					}
					logger.OutputNoteInfo("\t- content  = |%v|", callAnalysis.Content)
				}
			}
		}
		logger.OutputNoteInfo(ui.CommonNote2)

		logger.OutputNoteInfo(ui.CommonNote2)
		logger.OutputNoteInfo("Output package member function analysis:")
		for className, functionMap := range goAnalysis.PackageAnalysisMap[packagePath].StructAnalysis {
			logger.OutputNoteInfo("class %v", className)
			for functionName, functionAnalysis := range functionMap.MemberFunction {
				logger.OutputNoteInfo("\t- function %v", functionName)
				for paramName, functionVariableAnalysis := range functionAnalysis.ParamsMap {
					logger.OutputNoteInfo("\t\t- param index: %v, name: %v, type: %v, from: |%v %v|", functionVariableAnalysis.Index, paramName, functionVariableAnalysis.Type.Name, functionVariableAnalysis.Type.From, functionVariableAnalysis.Type.FromPackagePath)
				}
				for _, functionVariableAnalysis := range functionAnalysis.ReturnMap {
					logger.OutputNoteInfo("\t\t- return index: %v, name: %v, type: %v, from: |%v %v", functionVariableAnalysis.Index, functionVariableAnalysis.Name, functionVariableAnalysis.Type.Name, functionVariableAnalysis.Type.From, functionVariableAnalysis.Type.FromPackagePath)
				}
				for call, callAnalysisList := range functionAnalysis.CallMap {
					logger.OutputNoteInfo("\t\t- call: %v, len = %v", call, len(callAnalysisList))
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
							logger.OutputNoteInfo("\t\t\t- analysis = |%v.%v(%v)|", callAnalysis.From, callAnalysis.Call, paramListString)
						} else {
							logger.OutputNoteInfo("\t\t\t- analysis = |%v(%v)|", callAnalysis.Call, paramListString)
						}
						logger.OutputNoteInfo("\t\t\t- content = |%v|", callAnalysis.Content)
					}
				}
			}
		}
		logger.OutputNoteInfo(ui.CommonNote2)

		logger.OutputNoteInfo(ui.CommonNote2)
		logger.OutputNoteInfo("Output package const variable analysis:")
		for variableName, analysis := range goAnalysis.PackageAnalysisMap[packagePath].ConstAnalysisMap {
			logger.OutputNoteInfo("const variable %v, type %v, type from |%v %v|", variableName, analysis.Type.Name, analysis.Type.From, analysis.Type.FromPackagePath)
		}
		logger.OutputNoteInfo(ui.CommonNote2)

		logger.OutputNoteInfo(ui.CommonNote2)
		logger.OutputNoteInfo("Output package type rename analysis:")
		for rename, analysis := range goAnalysis.PackageAnalysisMap[packagePath].TypeRename {
			logger.OutputNoteInfo("type rename %v, type %v, type from |%v %v|", rename, analysis.Type.Name, analysis.Type.From, analysis.Type.FromPackagePath)
		}
		logger.OutputNoteInfo(ui.CommonNote2)
	}

	// 4.1.3.1.6.2.2
	for packagePath, packageAnalysis := range goAnalysis.PackageAnalysisMap {
		logger.OutputNoteInfo("merge function call for package: %v", packagePath)
		// non-member function
		for functionName, functionAnalysis := range packageAnalysis.FunctionAnalysisMap {
			logger.OutputNoteInfo("deal non-member function %v:", functionName)
			for callFunctionMame, callFunctionList := range functionAnalysis.CallMap {
				for index, functionCallAnalysis := range callFunctionList {
					logger.OutputNoteInfo("- call: %v, index: %v, param: %v, from: |%v %v|", index, callFunctionMame, functionCallAnalysis.ParamList, functionCallAnalysis.From, functionCallAnalysis.FromPackagePath)
					if len(functionCallAnalysis.From) == 0 || len(functionCallAnalysis.FromPackagePath) == 0 {
						utility2.TestOutput("callFunctionMame %v continue 1", callFunctionMame)
						continue
					}
					callFromPackageAnalysis, hasPackagePath := goAnalysis.PackageAnalysisMap[functionCallAnalysis.FromPackagePath]
					if !hasPackagePath {
						for key := range goAnalysis.PackageAnalysisMap {
							utility2.TestOutput("key = %v", key)
						}
						utility2.TestOutput("callFunctionMame %v continue 2: %v", callFunctionMame, functionCallAnalysis.FromPackagePath)
						continue
					}
					callFunctionAnalysis, hasFunction := callFromPackageAnalysis.FunctionAnalysisMap[callFunctionMame]
					if !hasFunction {
						utility2.TestOutput("callFunctionMame %v continue 3", callFunctionMame)
						continue
					}
					if _, hasCallerPackagePath := callFunctionAnalysis.CallerMap[packagePath]; !hasCallerPackagePath {
						callFunctionAnalysis.CallerMap[packagePath] = make(map[string]map[int]*GoFunctionCallAnalysis)
					}
					if _, hasCallerFunction := callFunctionAnalysis.CallerMap[packagePath]; !hasCallerFunction {
						callFunctionAnalysis.CallerMap[packagePath][functionName] = make(map[int]*GoFunctionCallAnalysis)
					}
					logger.OutputNoteInfo("package: %v, function: %v, No.%v calls package %v function %v with params: %v", packagePath, functionName, index, functionCallAnalysis.FromPackagePath, callFunctionMame, functionCallAnalysis.ParamList)
					callFunctionAnalysis.CallerMap[packagePath][functionName][index] = functionCallAnalysis
				}
			}
			logger.OutputNoteInfo(ui.CommonNote2)
		}
	}

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
func analyzeGoScopePackageVariable(content string, importAnalysisMap map[string]*GoImportAnalysis) *GoVariableAnalysis {
	goPackageVariableAnalysis := &GoVariableAnalysis{}
	for _, subMatchList := range regexps.GetRegexpByTemplateEnum(global.GoFileAnalyzerScopePackageVariableTemplate).FindAllStringSubmatch(content, -1) {
		if index, hasIndex := goAnalyzerPackageVariableSubMatchNameIndexMap["NAME"]; hasIndex {
			goPackageVariableAnalysis.Name = strings.TrimSpace(subMatchList[index])
		}
		if index, hasIndex := goAnalyzerPackageVariableSubMatchNameIndexMap["TYPE"]; hasIndex {
			goPackageVariableAnalysis.Type = analyzeGoVariableType(strings.TrimSpace(subMatchList[index]), importAnalysisMap)
		}
	}
	if len(goPackageVariableAnalysis.Name) == 0 || goPackageVariableAnalysis.Type == nil {
		return nil
	}
	return goPackageVariableAnalysis
}

// analyzeGoScopeConst
// @param constScope 待分析 const 域的内容
// @return
func analyzeGoScopeConst(constScope *scope, importAnalysisMap map[string]*GoImportAnalysis) map[string]*GoVariableAnalysis {
	constVariableMap := make(map[string]*GoVariableAnalysis)
	if constScope.isOneLineScope() {
		constStringList := strings.Split(constScope.Content, global.GoSplitterStringSpace)
		constVariableMap[constStringList[1]] = &GoVariableAnalysis{
			Name: constStringList[1],
		}
		constVariableMap[constStringList[1]].Type = analyzeGoVariableType(constStringList[2], importAnalysisMap)
	} else {
		constScopeRootNode := utility3.TraitMultiPunctuationMarksContent(constScope.Content, global.GoAnalyzerScopePunctuationMarkList, 1)
		if len(constScopeRootNode.SubPunctuationContentList) < 1 {
			return constVariableMap
		}
		var constVariableType *GoTypeAnalysis
		for _, constVariableString := range strings.Split(constScopeRootNode.SubPunctuationContentList[0].Content, global.GoSplitterStringEnter) {
			constVariableStringList := strings.Split(strings.TrimSpace(constVariableString), global.GoSplitterStringSpace)
			if len(constVariableStringList) == 0 {
				continue
			}
			if len(constVariableStringList) > 1 {
				constVariableType = analyzeGoVariableType(constVariableStringList[1], importAnalysisMap)
			}
			constVariableMap[constVariableStringList[0]] = &GoVariableAnalysis{
				Name: constVariableStringList[0],
				Type: constVariableType,
			}
		}
	}
	return constVariableMap
}

// analyzeGoScopeInterface
// @param content 待分析 interface 域的内容
// @return
func analyzeGoScopeInterface(content string, importAnalysisMap map[string]*GoImportAnalysis) *GoInterfaceAnalysis {
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

		goFunctionAnalysis := analyzeGoFunctionDefinition(replacedContent, false, importAnalysisMap)
		goInterfaceAnalysis.Function[goFunctionAnalysis.Name] = goFunctionAnalysis
	}
	// utility2.TestOutput(ui.CommonNote2)

	return goInterfaceAnalysis
}

// analyzeGoScopeStruct
// @param structScope 待分析的 struct 域
// @return
func analyzeGoScopeStruct(structScope *scope, importAnalysisMap map[string]*GoImportAnalysis) *GoStructAnalysis {
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
			goVariableAnalysis.Type = analyzeGoVariableType(structContent, importAnalysisMap)
			goVariableAnalysis.Name = goVariableAnalysis.Type.Name
			goStructAnalysis.Base[goVariableAnalysis.Name] = goVariableAnalysis
		} else {
			goVariableAnalysis.Name = structContent[0:spaceIndex]
			goVariableAnalysis.Type = analyzeGoVariableType(structContent[spaceIndex+1:], importAnalysisMap)
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
					goVariableAnalysis.Type = analyzeGoVariableType(strings.TrimSpace(subMatchList[index]), importAnalysisMap)
				}
			}
			if len(goVariableAnalysis.Name) == 0 {
				continue
			}
			if goVariableAnalysis.Type == nil {
				goVariableAnalysis.Type = analyzeGoVariableType(strings.TrimSpace(goVariableAnalysis.Name), importAnalysisMap)
				goVariableAnalysis.Name = goVariableAnalysis.Type.Name
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
func analyzeGoScopeFunction(functionScope *scope, isMemberFunction bool, importAnalysisMap map[string]*GoImportAnalysis) *GoFunctionAnalysis {
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
	goFunctionAnalysis := analyzeGoFunctionDefinition(replacedContent[:contentRootNode.SubPunctuationContentList[subNodeCount-1].LeftPunctuationMark.Index], isMemberFunction, importAnalysisMap)
	goFunctionAnalysis.CallMap = analyzeGoFunctionBody(contentRootNode.SubPunctuationContentList[subNodeCount-1].Content, importAnalysisMap)

	// utility2.TestOutput(ui.CommonNote2)

	return goFunctionAnalysis
}

// analyzeGoFunctionDefinition
// @param functionDefinitionContent 待分析的函数定义的内容
// @return
func analyzeGoFunctionDefinition(functionDefinitionContent string, isMemberFunction bool, importAnalysisMap map[string]*GoImportAnalysis) *GoFunctionAnalysis {
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

	// function class
	var functionClass string
	var variableThis *GoVariableAnalysis
	paramListNodeIndex := 0
	if isMemberFunction {
		paramListNodeIndex = 1
		if subNodeCount < 2 {
			return nil
		}
		variableThisString := functionDefinitionContentRootNode.SubPunctuationContentList[0].Content
		functionClass, variableThis = analyzeGoMemberFunctionVariableThis(variableThisString, importAnalysisMap)
		if len(functionClass) == 0 || variableThis == nil {
			return nil
		}
	}

	// param list
	paramListString := functionDefinitionContentRootNode.SubPunctuationContentList[paramListNodeIndex].Content
	paramListNode := utility3.RecursiveSplitUnderSameDeepPunctuationMarksContent(paramListString, global.GoAnalyzerScopePunctuationMarkList, global.GoSplitterStringComma)
	// utility2.TestOutput("param list:")
	// for index, param := range paramListNode.ContentList {
	// 	utility2.TestOutput("index %v param = |%v|", index, param)
	// }
	paramMap := analyzeGoFunctionDefinitionParamList(paramListNode.ContentList, importAnalysisMap)
	// for paramName, analysis := range paramMap {
	// 	utility2.TestOutput("paramName = %v, index = %v, type = |%v|, type from = %v", paramName, analysis.Index, analysis.Type, analysis.TypeFrom)
	// }

	// return list
	returnListString := strings.TrimSpace(functionDefinitionContent[functionDefinitionContentRootNode.SubPunctuationContentList[paramListNodeIndex].RightPunctuationMark.Index+1:])
	returnListString = strings.TrimFunc(returnListString, func(r rune) bool {
		return r == global.PunctuationMarkLeftBracket || r == global.PunctuationMarkRightBracket
	})
	// utility2.TestOutput("return list = |%v|", returnListString)
	returnListNode := utility3.RecursiveSplitUnderSameDeepPunctuationMarksContent(returnListString, global.GoAnalyzerScopePunctuationMarkList, global.GoSplitterStringComma)
	// utility2.TestOutput("return list:")
	// for index, returnType := range returnListNode.ContentList {
	// 	utility2.TestOutput("index %v return = |%v|", index, returnType)
	// }
	returnMap := analyzeGoFunctionDefinitionReturnList(returnListNode.ContentList, importAnalysisMap)
	// for returnName, analysis := range returnMap {
	// 	utility2.TestOutput("returnName = %v, index = %v, type = |%v|, type from = %v", returnName, analysis.Index, analysis.Type, analysis.TypeFrom)
	// }

	return &GoFunctionAnalysis{
		Class:        functionClass,
		VariableThis: variableThis,
		Name:         name,
		ParamsMap:    paramMap,
		ReturnMap:    returnMap,
	}
}

// analyzeGoMemberFunctionVariableThis
// @param variableThisString
// @return
func analyzeGoMemberFunctionVariableThis(variableThisString string, importAnalysisMap map[string]*GoImportAnalysis) (string, *GoVariableAnalysis) {
	variableThisStringList := strings.Split(strings.TrimSpace(variableThisString), " ")
	if len(variableThisStringList) < 2 {
		return "", nil
	}
	functionClass := utility.TraitStructName(variableThisStringList[1])
	variableThisAnalysis := &GoVariableAnalysis{
		Name: variableThisStringList[0],
	}
	variableThisAnalysis.Type = analyzeGoVariableType(variableThisStringList[1], importAnalysisMap)
	return functionClass, variableThisAnalysis
}

// analyzeGoFunctionDefinitionParamList
// @param paramListString 待分析的函数参数表
// @return
func analyzeGoFunctionDefinitionParamList(paramStringList []string, importAnalysisMap map[string]*GoImportAnalysis) map[string]*GoFunctionVariable {
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
			paramType := analyzeGoVariableType(paramStringTimSpace[splitterIndex+1:], importAnalysisMap)
			paramMap[paramStringTimSpace[:splitterIndex]] = &GoFunctionVariable{
				GoVariableAnalysis: GoVariableAnalysis{
					Name: paramStringTimSpace[:splitterIndex],
					Type: paramType,
				},
				Index: index,
			}
			if len(unknownTypeParamList) != 0 {
				for _, unknownTypeParam := range unknownTypeParamList {
					unknownTypeParam.GoVariableAnalysis.Type = paramType
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
func analyzeGoFunctionDefinitionReturnList(returnStringList []string, importAnalysisMap map[string]*GoImportAnalysis) map[string]*GoFunctionVariable {
	returnMap := make(map[string]*GoFunctionVariable)

	for index, returnString := range returnStringList {
		returnStringTimSpace := strings.TrimSpace(returnString)
		if len(returnStringTimSpace) == 0 {
			continue
		}
		splitterIndex := strings.Index(returnStringTimSpace, global.GoSplitterStringSpace)
		if splitterIndex == -1 {
			returnType := analyzeGoVariableType(returnStringTimSpace, importAnalysisMap)
			returnMap[fmt.Sprintf("%v", len(returnMap))] = &GoFunctionVariable{
				GoVariableAnalysis: GoVariableAnalysis{
					Type: returnType,
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
			var returnType *GoTypeAnalysis
			if keywordIndex != -1 {
				returnName := strings.TrimSpace(returnStringTimSpace[:keywordIndex])
				returnType = analyzeGoVariableType(returnStringTimSpace[keywordIndex:], importAnalysisMap)
				if len(returnName) == 0 {
					returnMapKey = fmt.Sprintf("%v", len(returnMap))
				} else {
					returnMapKey = returnName
				}
			} else {
				returnType = analyzeGoVariableType(returnStringTimSpace, importAnalysisMap)
				returnMapKey = fmt.Sprintf("%v", len(returnMap))
			}

			returnMap[returnMapKey] = &GoFunctionVariable{
				GoVariableAnalysis: GoVariableAnalysis{
					Name: returnName,
					Type: returnType,
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
func analyzeGoFunctionBody(functionBodyContent string, importAnalysisMap map[string]*GoImportAnalysis) map[string][]*GoFunctionCallAnalysis {
	callMap := make(map[string][]*GoFunctionCallAnalysis)
	bodyContentRootNode := utility3.TraitMultiPunctuationMarksContent(functionBodyContent, global.GoAnalyzerScopePunctuationMarkList, -1)
	toSearchNodeList := []*utility.NewPunctuationContent{bodyContentRootNode}
	bracketScopeNodeList := make([]*utility.NewPunctuationContent, 0)
	// bracketScopeContentList := make([]string, 0)

	// BF-Traversal
	for len(toSearchNodeList) != 0 {
		toSearchNode := toSearchNodeList[0]
		if toSearchNode == nil {
			continue
		}

		// utility2.TestOutput("content |%v%v%v|, index testContent[%v:%v]\n", string(toSearchNode.LeftPunctuationMark.PunctuationMark), toSearchNode.Content, string(toSearchNode.RightPunctuationMark.PunctuationMark), toSearchNode.LeftPunctuationMark.Index, toSearchNode.RightPunctuationMark.Index)
		if toSearchNode.LeftPunctuationMark.Index >= 0 && toSearchNode.RightPunctuationMark.Index < len(toSearchNode.Content) {
			// utility2.TestOutput("content by index = |%v|\n", functionBodyContent[toSearchNode.LeftPunctuationMark.Index+1:toSearchNode.RightPunctuationMark.Index])
		}

		if toSearchNode.LeftPunctuationMark.PunctuationMark == global.PunctuationMarkLeftBracket {
			bracketScopeNodeList = append(bracketScopeNodeList, toSearchNode)
		}
		toSearchNodeList = toSearchNodeList[1:]
		toSearchNodeList = append(toSearchNodeList, toSearchNode.SubPunctuationContentList...)
	}

	// merge same start
	sameStartIndexNodeMap := make(map[int][]*utility.NewPunctuationContent)
	for _, bracketScopeNode := range bracketScopeNodeList {
		// utility2.TestOutput("bracketScopeNode left index = |%v|", bracketScopeNode.LeftPunctuationMark.Index)
		startIndex := searchCallStartIndex(functionBodyContent, bracketScopeNode.LeftPunctuationMark.Index)
		if startIndex == -1 {
			// utility2.TestOutput("continue node: |%+v|", bracketScopeNode)
			continue
		}
		// utility2.TestOutput("call identifier functionBody[%v:%v]", startIndex, bracketScopeNode.LeftPunctuationMark.Index)
		// utility2.TestOutput("call identifier = |%v|", functionBodyContent[startIndex:bracketScopeNode.LeftPunctuationMark.Index])
		// utility2.TestOutput(ui.CommonNote2)
		if _, has := sameStartIndexNodeMap[startIndex]; !has {
			sameStartIndexNodeMap[startIndex] = make([]*utility.NewPunctuationContent, 0)
		}
		// utility2.TestOutput("append start index %v, node content = |%v|", startIndex, bracketScopeNode.Content)
		sameStartIndexNodeMap[startIndex] = append(sameStartIndexNodeMap[startIndex], bracketScopeNode)
	}

	for sameStartIndex, bracketScopeNodeList := range sameStartIndexNodeMap {
		// utility2.TestOutput("deal start index %v call, len = %v", sameStartIndex, len(bracketScopeNodeList))
		var rootCall *GoFunctionCallAnalysis
		var lastCall *GoFunctionCallAnalysis
		for index, bracketScopeNode := range bracketScopeNodeList {
			var identifierString string
			// utility2.TestOutput("index %v, full call chain = |%v|", index, functionBodyContent[sameStartIndex:bracketScopeNode.LeftPunctuationMark.Index])
			if index == 0 {
				identifierString = functionBodyContent[sameStartIndex:bracketScopeNode.LeftPunctuationMark.Index]
			} else {
				identifierString = functionBodyContent[bracketScopeNodeList[index-1].RightPunctuationMark.Index+2 : bracketScopeNode.LeftPunctuationMark.Index]
			}
			// utility2.TestOutput("index %v current call = |%v|", index, identifierString)
			if len(identifierString) == 0 {
				continue
			}

			callIdentifierStringList := strings.Split(identifierString, global.GoSplitterStringPoint)
			if len(callIdentifierStringList) == 0 {
				// utility2.TestOutput("callIdentifierStringList len is 0")
				continue
			}

			callIdentifier := callIdentifierStringList[len(callIdentifierStringList)-1]
			if callIdentifier == global.GoKeywordFunc {
				continue
			}
			// utility2.TestOutput("call |%v|", callIdentifier)

			var callIdentifierFrom string
			var callIdentifierFromPackagePath string
			if len(callIdentifierStringList) > 1 {
				callIdentifierFrom = strings.Join(callIdentifierStringList[:len(callIdentifierStringList)-1], global.GoSplitterStringPoint)
				if _, hasAlias := importAnalysisMap[callIdentifierFrom]; hasAlias {
					callIdentifierFromPackagePath = importAnalysisMap[callIdentifierFrom].Path
				}
				// for
				// utility2.TestOutput("From |%v|", callIdentifierFrom)
			}

			paramListRootNode := utility3.RecursiveSplitUnderSameDeepPunctuationMarksContent(bracketScopeNode.Content, global.GoAnalyzerScopePunctuationMarkList, global.GoSplitterStringComma)
			paramList := make([]string, 0, len(paramListRootNode.ContentList))
			for _, param := range paramListRootNode.ContentList {
				if len(strings.TrimSpace(param)) != 0 {
					paramList = append(paramList, strings.TrimSpace(param))
					// utility2.TestOutput("index %v param = %v", index, strings.TrimSpace(param))
				}
			}

			callAnalysis := &GoFunctionCallAnalysis{
				Content:         fmt.Sprintf("%v(%v)", identifierString, bracketScopeNode.Content),
				From:            callIdentifierFrom,
				FromPackagePath: callIdentifierFromPackagePath,
				Call:            callIdentifier,
				ParamList:       paramList,
			}
			callAnalysis.PreviousCall = lastCall
			if lastCall != nil {
				lastCall.PostCall = callAnalysis
			}
			lastCall = callAnalysis
			if rootCall == nil {
				rootCall = callAnalysis
			}

			callMap[identifierString] = append(callMap[identifierString], callAnalysis)
		}
		// utility2.TestOutput(ui.CommonNote2)
	}

	return callMap
}

// searchCallStartIndex
// @param searchString 待查找的字符串
// @param searchIndex 待搜索的起始下标
// @return
func searchCallStartIndex(searchString string, searchIndex int) int {
	// utility2.TestOutput("searchCallStartIndex, searchIndex = %v, rune = %v", searchIndex, string(searchString[searchIndex]))
	// utility2.TestOutput("left index content = |%v|", searchString[:searchIndex])
	if searchIndex > len(searchString) {
		return -1
	}
	invalidScope := false
	for index := searchIndex - 1; index > -1; index-- {
		r := rune(searchString[index])

		if invalidScope {
			continue
		}

		if isInvalidScopeRune(r) {
			invalidScope = !invalidScope
			continue
		}

		if isIdentifierScopeEndRune(r) {
			// utility2.TestOutput("into scope, index = %v, rune = %v", index, string(searchString[index])) // ...(...).Func() -> ')' means into scope
			// utility2.TestOutput("calculate from utility.ReverseString(searchString[:index] = |%v|", utility.ReverseString(searchString[:index]))
			length := utility3.CalculatePunctuationMarksContentLength(utility.ReverseString(searchString[:index]), r, utility.GetAnotherPunctuationMark(r), global.GoAnalyzerInvalidScopePunctuationMarkMap)
			// utility2.TestOutput("content length = %v", length)
			if length == -1 {
				// means the scope end rune is invalid, just continue
				continue
			}
			index--         // cut ')' it self length, index points to scope content end position // searchString[:index] = '...(...).Func()'[:index] -> ...(...
			index -= length // cut content length, index points to scope start rune // searchString[index] = '('
			// utility2.TestOutput("scope content = |%v|", searchString[index+1:index+length+1])
			// utility2.TestOutput("after calculate, index = %v, rune = |%v|", index, string(searchString[index]))
			continue // index-- by cycle logic
		}

		if isIdentifierRune(r) {
			continue
		}
		// utility2.TestOutput("return index = %v", index+1)
		return index + 1
	}
	// utility2.TestOutput("return index = 0")
	return 0
}

// isIdentifierRune
// @param r 待检查的字符
// @return
func isIdentifierRune(r rune) bool {
	// . || _ || A~Z || a~z
	return r == 46 || r == 95 || (48 <= r && r <= 57) || (65 <= r && r <= 90) || (97 <= r && r <= 123)
}

func isInvalidScopeRune(r rune) bool {
	// " || ' || `
	return r == 34 || r == 39 || r == 96
}

func isIdentifierScopeStartRune(r rune) bool {
	// ( || [ || {
	return r == 40 || r == 91 || r == 121
}

func isIdentifierScopeEndRune(r rune) bool {
	// ) || ] || }
	return r == 41 || r == 93 || r == 125
}

// analyzeGoVariableType
// @param variableTypeString 待分析变量的类型字符串
// @return
func analyzeGoVariableType(variableTypeString string, importAnalysisMap map[string]*GoImportAnalysis) *GoTypeAnalysis {
	var typeName string
	var typeFrom string
	var typeFromPackagePath string
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

	if len(typeName) == 0 {
		return nil
	}

	if len(typeFrom) != 0 {
		if _, hasAlias := importAnalysisMap[typeFrom]; hasAlias {
			typeFromPackagePath = importAnalysisMap[typeFrom].Path
		}
	}

	return &GoTypeAnalysis{
		Name:            typeName,
		From:            typeFrom,
		FromPackagePath: typeFromPackagePath,
	}
}

// analyzeGoTypeRename
// @param 待分析的
func analyzeGoTypeRename(renameScope *scope, importAnalysisMap map[string]*GoImportAnalysis) *GoTypeRenameAnalysis {
	renameScopeRootNode := utility3.RecursiveSplitUnderSameDeepPunctuationMarksContent(renameScope.Content, global.GoAnalyzerScopePunctuationMarkList, global.GoSplitterStringSpace)
	if renameScopeRootNode == nil || len(renameScopeRootNode.ContentList) < 3 {
		return nil
	}
	typeRenameAnalysis := &GoTypeRenameAnalysis{
		Name: renameScopeRootNode.ContentList[1],
	}
	typeRenameAnalysis.Type = analyzeGoVariableType(strings.Join(renameScopeRootNode.ContentList[2:], global.GoSplitterStringSpace), importAnalysisMap)
	return typeRenameAnalysis
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

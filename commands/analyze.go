package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/go-worker/config"
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
	utility2.TestOutput("analyze execute")
	// 解析指令的选项和参数
	parseCommandParamsError := command.parseCommandParams()
	if parseCommandParamsError != nil {
		return parseCommandParamsError
	}

	// 获取全局配置
	projectPath := config.GetCurrentProjectPath()
	fileType := config.WorkerConfig.ProjectSyntax
	if fileType == "" {
		ui.OutputWarnInfo(ui.CommonWarn1)
	}

	if command.Params.parentValue != "" {
		projectPath = fmt.Sprintf("%v/%v", projectPath, command.Params.parentValue)
	}

	if command.Params.sourceValue == "" {
		return fmt.Errorf(ui.CommonError1)
	}

	toAnalyzeFilePath := fmt.Sprintf("%v/%v.%v", projectPath, command.Params.sourceValue, fileType)
	if !utility.IsExist(toAnalyzeFilePath) {
		return fmt.Errorf(ui.CMDAnalyzeFileNotExist)
	}

	toWriteFilePath := fmt.Sprintf("%v/%v.%v.%v", projectPath, command.Params.sourceValue, fileType, global.SyntaxMarkdown)
	if command.Params.outputValue != "" {
		toWriteFilePath = fmt.Sprintf("%v/%v", projectPath, command.Params.outputValue)
	}

	var toWriteFile *os.File
	defer func() {
		if toWriteFile != nil {
			toWriteFile.Close()
		}
	}()
	if utility.IsExist(toWriteFilePath) {
		var openFileError error
		toWriteFile, openFileError = os.OpenFile(toWriteFilePath, os.O_RDWR|os.O_APPEND, 0644)
		if openFileError != nil {
			return openFileError
		}
	} else {
		var createFileError error
		toWriteFile, createFileError = utility.CreateFile(toWriteFilePath)
		if createFileError != nil {
			return createFileError
		}
	}

	if toWriteFile == nil {
		return fmt.Errorf("analyze to write file is nil")
	}

	toAnalyzeFile, inputError := os.Open(toAnalyzeFilePath)
	defer toAnalyzeFile.Close()
	if inputError != nil || toAnalyzeFile == nil {
		return fmt.Errorf(ui.CommonError5, toAnalyzeFilePath, inputError.Error())
	}

	var analyzeError error
	switch fileType {
	case global.SyntaxGo:
		analyzeError = analyzeGoFile(toAnalyzeFile, toWriteFile)
	case global.SyntaxCpp:
		analyzeError = analyzeCppFile(toAnalyzeFile, toWriteFile)
	}
	if analyzeError != nil {
		ui.OutputWarnInfo(ui.CMDAnalyzeOccursError, analyzeError)
	}

	return nil
}

type analyzeParam struct {
	sourceType  string
	sourceValue string
	parentValue string
	outputValue string
}

func (command *Analyze) parseCommandParams() error {
	optionValueString := ""
	if optionValueRegexp, hasOptionValueRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEAnalyzeOptionValue]; hasOptionValueRegexp {
		optionValueString = optionValueRegexp.FindString(command.CommandStruct.InputString)
	} else {
		ui.OutputWarnInfo(ui.CommonWarn2, "analyze", "file|directory|package")
	}
	if optionValueString == "" {
		return fmt.Errorf(ui.CommonError1)
	}
	optionValueList := strings.Split(optionValueString, " ")
	command.Params = &analyzeParam{
		sourceType:  optionValueList[0],
		sourceValue: optionValueList[1],
	}
	parentValue := ""
	if parentValueRegexp := regexps.GetRegexpByTemplateEnum(global.OptionParentValueTemplate); parentValueRegexp != nil {
		parentValue = parentValueRegexp.FindString(command.CommandStruct.InputString)
	} else {
		ui.OutputWarnInfo(ui.CommonWarn2, "analyze", "parent")
	}
	if parentValue != "" {
		parentValueList := strings.Split(parentValue, " ")
		command.Params.parentValue = parentValueList[1]
	}
	outputValue := ""
	if outputValueRegexp := regexps.GetRegexpByTemplateEnum(global.OptionOutputValueTemplate); outputValueRegexp != nil {
		outputValue = outputValueRegexp.FindString(command.CommandStruct.InputString)
	} else {
		ui.OutputWarnInfo(ui.CommonWarn2, "analyze", "output")
	}
	if outputValue != "" {
		outputValueList := strings.Split(outputValue, " ")
		command.Params.outputValue = outputValueList[1]
	}
	return nil
}

// 分析 GO 文件

// GoFileAnalysis go 文件分析结果
type GoFileAnalysis struct {
	FilePath     string
	PackageName  string
	ImportList   []string
	FunctionMap  map[string]*GoFunctionAnalysis
	functionList []string
}

// GoFunctionAnalysis go 函数分析结果
type GoFunctionAnalysis struct {
	Class     string
	Name      string
	ParamsMap map[string]string
	ReturnMap map[string]string
}

func analyzeGoFile(toAnalyzeFile, toWriteFile *os.File) error {
	goFileAnalysis := &GoFileAnalysis{
		ImportList:   make([]string, 0),
		FunctionMap:  make(map[string]*GoFunctionAnalysis),
		functionList: make([]string, 0),
	}

	toAnalyzeContent, readToAnalyzeContentError := ioutil.ReadAll(toAnalyzeFile)
	if readToAnalyzeContentError != nil {
		return readToAnalyzeContentError
	}

	// 文件路径
	goFileAnalysis.FilePath = toAnalyzeFile.Name()

	// 解析包名
	analyzeGoKeywordPackage(goFileAnalysis, toAnalyzeContent)

	// 解析依赖包
	analyzeGoImportPackage(goFileAnalysis, toAnalyzeContent)

	// 解析函数定义
	analyzeGoFunctionDefinition(goFileAnalysis, toAnalyzeContent)

	// 输出解析结果
	functionListContent := outputAnalyzeGoFileResult(goFileAnalysis)

	// 输出到文件
	_, writeError := toWriteFile.WriteString(functionListContent)
	if writeError != nil {
		return writeError
	}

	return nil
}

func analyzeGoKeywordPackage(goFileAnalysis *GoFileAnalysis, fileContentByte []byte) {
	if keywordPackageRegexp, hasKeywordPackageRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEGoKeywordPackageValue]; hasKeywordPackageRegexp {
		if packageValueContentByte := keywordPackageRegexp.Find(fileContentByte); len(packageValueContentByte) != 0 {
			goFileAnalysis.PackageName = strings.Split(string(packageValueContentByte), " ")[1]
		}
	} else {
		ui.OutputWarnInfo(ui.CMDAnalyzeGoKeywordRegexpNotExist, "package")
	}
}

func analyzeGoImportPackage(goFileAnalysis *GoFileAnalysis, fileContentByte []byte) {
	if keywordImportValueRegexp := regexps.GetRegexpByTemplateEnum(global.GoKeywordImportValueTemplate); keywordImportValueRegexp != nil {
		if doubleQuotesContentRegexp, hasDoubleQuotesContentRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEDoubleQuotesContent]; hasDoubleQuotesContentRegexp {
			for _, matchedKeywordImportValue := range keywordImportValueRegexp.FindAll(fileContentByte, -1) {
				for _, quotesContentByte := range doubleQuotesContentRegexp.FindAll(matchedKeywordImportValue, -1) {
					goFileAnalysis.ImportList = append(goFileAnalysis.ImportList, strings.Trim(string(quotesContentByte), "\""))
				}
			}
		} else {
			ui.OutputWarnInfo(ui.CommonWarn3, global.AEDoubleQuotesContent)
		}
	} else {
		ui.OutputWarnInfo(ui.CMDAnalyzeGoKeywordRegexpNotExist, "import")
	}
}

func analyzeGoFunctionDefinition(goFileAnalysis *GoFileAnalysis, fileContentByte []byte) {
	bracketsContentRegexp, hasBracketsContentRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEBracketsContent]
	if !hasBracketsContentRegexp {
		ui.OutputWarnInfo(ui.CommonWarn3, global.AEBracketsContent)
		return
	}
	functionDefinitionRegexp := regexps.GetRegexpByTemplateEnum(global.GoFunctionDefinitionTemplate)
	if functionDefinitionRegexp == nil {
		ui.OutputWarnInfo(ui.CMDAnalyzeGoKeywordRegexpNotExist, "function")
		return
	}
	functionDefinitionByteList := functionDefinitionRegexp.FindAll(fileContentByte, -1)
	// 解析函数定义
	for _, functionDefinitionByte := range functionDefinitionByteList {
		functionAnalysis := &GoFunctionAnalysis{
			ParamsMap: make(map[string]string),
			ReturnMap: make(map[string]string),
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
			functionAnalysis.Class = utility.TraitStructName(memberStringList[1])
		}
		// 解析函数名
		functionAnalysis.Name = functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$NAME")
		// 解析参数表
		paramsStringWithPunctuationMark := functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$PARAM")
		if paramsStringWithPunctuationMark != "" {
			paramsString := bracketsContentRegexp.ReplaceAllString(paramsStringWithPunctuationMark, "$CONTENT")
			typeString := ""
			unknownTypeNameList := make([]string, 0)
			for _, paramString := range strings.Split(paramsString, ",") {
				paramStringList := strings.Split(strings.TrimSpace(paramString), " ")
				if len(paramStringList) == 1 {
					unknownTypeNameList = append(unknownTypeNameList, paramStringList[0])
				} else {
					functionAnalysis.ParamsMap[paramStringList[0]] = paramStringList[1]
					if typeString == "" {
						for _, unknownTypeName := range unknownTypeNameList {
							functionAnalysis.ParamsMap[unknownTypeName] = paramStringList[1]
						}
					}
				}
			}
		}
		// 解析返回列表
		returnStringWithPunctuationMark := functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$RETURN")
		if returnStringWithPunctuationMark != "" {
			returnContent := bracketsContentRegexp.ReplaceAllString(returnStringWithPunctuationMark, "$CONTENT")
			for _, returnString := range strings.Split(returnContent, ",") {
				returnStringList := strings.Split(strings.TrimSpace(returnString), " ")
				if len(returnStringList) == 1 {
					functionAnalysis.ReturnMap[fmt.Sprintf("%v", len(functionAnalysis.ReturnMap))] = returnStringList[0]
				} else if len(returnStringList) == 2 {
					functionAnalysis.ReturnMap[returnStringList[0]] = returnStringList[1]
				}
			}
		}
		goFileAnalysis.FunctionMap[functionAnalysis.Name] = functionAnalysis
		goFileAnalysis.functionList = append(goFileAnalysis.functionList, functionAnalysis.Name)
	}

	// 解析函数体
	// functionDefinitionIndexList := functionDefinitionRegexp.FindAllIndex(fileContentByte, -1)
	// for index, functionDefinitionIndex := range functionDefinitionIndexList {
	// 	definitionLength := utility.CalculatePunctuationMarksContentLength(string(fileContentByte[functionDefinitionIndex[1]+1:]), global.PunctuationMarkCurlyBraces)
	// 	if definitionLength == -1 {
	// 		ui.OutputWarnInfo(ui.CMDAnalyzeGoFunctionContentSyntax)
	// 		continue
	// 	}
	// 	utility2.TestOutput("function[%v] end index = %v", index, functionDefinitionIndex[1]+1+definitionLength)
	// 	functionContent := fileContentByte[functionDefinitionIndex[0] : functionDefinitionIndex[1]+1+definitionLength]
	// 	utility2.TestOutput("function[%v] functionContent = |%v|", index, string(functionContent))
	// 	utility2.TestOutput("---------------- splitter ----------------")
	// }
}

func outputAnalyzeGoFileResult(goFileAnalysis *GoFileAnalysis) string {
	templateStyleRegexp, hasTemplateStyleRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AETemplateStyle]
	if !hasTemplateStyleRegexp {
		ui.OutputWarnInfo(ui.CommonWarn3, global.AETemplateStyle)
		return ""
	}

	// file path
	resultContent := strings.Replace(ui.AnalyzeGoFileResultTemplate, global.AnalyzeRPFilePath, goFileAnalysis.FilePath, -1)

	// package name
	resultContent = strings.Replace(resultContent, global.AnalyzeRPPackageName, goFileAnalysis.PackageName, -1)

	// import 内容
	importPackageListContent := ""
	if len(goFileAnalysis.ImportList) != 0 {
		importPackageListString := ""
		for _, packagePath := range goFileAnalysis.ImportList {
			packagePathContent := ui.AnalyzeGoFileImportPackageTemplate

			// style template
			packagePathContent = ui.ParseStyleTemplate(templateStyleRegexp, packagePathContent)

			// package path
			packagePathContent = strings.Replace(packagePathContent, global.AnalyzeRPPackagePath, packagePath, -1)

			if importPackageListString == "" {
				importPackageListString = packagePathContent
			} else {
				importPackageListString = fmt.Sprintf("%v\n%v", importPackageListString, packagePathContent)
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
				functionParamListContent = parseGoFunctionParamOrReturnListContent(templateStyleRegexp, functionAnalysis.ParamsMap, functionParamList)
			}
			functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionParamList, functionParamListContent, -1)

			// function return
			functionReturnListContent := ""
			if len(functionAnalysis.ReturnMap) != 0 {
				functionReturnListContent = parseGoFunctionParamOrReturnListContent(templateStyleRegexp, functionAnalysis.ReturnMap, functionReturnList)
			}
			functionDefinitionContent = strings.Replace(functionDefinitionContent, global.AnalyzeRPFunctionReturnList, functionReturnListContent, -1)

			// clear space line
			functionDefinitionContent = utility3.TrimSpaceLine(functionDefinitionContent)

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

	ui.OutputNoteInfo(resultContent)

	return resultContent
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

// 分析 CPP 文件

func analyzeCppFile(toAnalyzeFile, toWriteFile *os.File) error {
	return nil
}

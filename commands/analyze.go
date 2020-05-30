package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-worker/config"
	"github.com/go-worker/global"
	"github.com/go-worker/regexps"
	"github.com/go-worker/ui"
	"github.com/go-worker/utility"
)

type Analyze struct {
	*CommandStruct
	Params *analyzeParam
}

func (command *Analyze) Execute() error {
	utility.TestOutput("analyze execute")
	// 解析指令的选项和参数
	parseCommandParamsError := command.parseCommandParams()
	if parseCommandParamsError != nil {
		return parseCommandParamsError
	}
	utility.TestOutput("command.Params = %+v", command.Params)

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

	toWriteFilePath := fmt.Sprintf("%v/%v.%v.analyze", projectPath, command.Params.sourceValue, fileType)
	if command.Params.outputValue != "" {
		toWriteFilePath = fmt.Sprintf("%v/%v", projectPath, command.Params.outputValue)
	}

	utility.TestOutput("toAnalyzeFilePath = %v, toWriteFilePath = %v", toAnalyzeFilePath, toWriteFilePath)

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

	_, writeError := toWriteFile.WriteString(fmt.Sprintf("toAnalyzeFilePath = %v, toWriteFilePath = %v", toAnalyzeFilePath, toWriteFilePath))
	if writeError != nil {
		return writeError
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
	FilePath    string
	PackageName string
	ImportList  []string
	FunctionMap map[string]*GoFunctionAnalysis
}

// GoFunctionAnalysis go 函数分析结果
type GoFunctionAnalysis struct {
	Member    string
	Name      string
	ParamsMap map[string]string
	ReturnMap map[string]string
}

func analyzeGoFile(toAnalyzeFile, toWriteFile *os.File) error {
	goFileAnalysis := &GoFileAnalysis{
		ImportList:  make([]string, 0),
		FunctionMap: make(map[string]*GoFunctionAnalysis),
	}

	toAnalyzeContent, readToAnalyzeContentError := ioutil.ReadAll(toAnalyzeFile)
	if readToAnalyzeContentError != nil {
		return readToAnalyzeContentError
	}

	// 解析包名
	analyzeGoKeywordPackage(goFileAnalysis, toAnalyzeContent)

	// 解析依赖包
	analyzeGoImportPackage(goFileAnalysis, toAnalyzeContent)

	// 解析函数定义
	analyzeGoFunctionDefinition(goFileAnalysis, toAnalyzeContent)

	return nil
}

func analyzeGoKeywordPackage(goFileAnalysis *GoFileAnalysis, fileContentByte []byte) {
	if keywordPackageRegexp, hasKeywordPackageRegexp := regexps.AtomicExpressionEnumRegexpMap[global.AEGoKeywordPackageValue]; hasKeywordPackageRegexp {
		if packageValueContentByte := keywordPackageRegexp.Find(fileContentByte); len(packageValueContentByte) != 0 {
			goFileAnalysis.PackageName = strings.Split(string(packageValueContentByte), " ")[1]
			utility.TestOutput("analysis file package name = %v", goFileAnalysis.PackageName)
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
					goFileAnalysis.ImportList = append(goFileAnalysis.ImportList, string(quotesContentByte))
					utility.TestOutput("import package = %v", string(quotesContentByte))
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
		functionAnalysis := &GoFunctionAnalysis{}
		// utility.TestOutput("MEMBER = %v", functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$MEMBER"))
		for _, memberType := range utility.TraitPunctuationMarksContentAsNameTypeMap(functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$MEMBER"), bracketsContentRegexp, global.SyntaxGo) {
			functionAnalysis.Member = memberType
		}
		// utility.TestOutput("NAME = %v", functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$NAME"))
		functionAnalysis.Name = functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$NAME")
		// utility.TestOutput("PARAM = %v", functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$PARAM"))
		functionAnalysis.ParamsMap = utility.TraitPunctuationMarksContentAsNameTypeMap(functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$PARAM"), bracketsContentRegexp, global.SyntaxGo)
		// utility.TestOutput("RETURN = %v", functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$RETURN"))
		functionAnalysis.ReturnMap = utility.TraitPunctuationMarksContentAsNameTypeMap(functionDefinitionRegexp.ReplaceAllString(string(functionDefinitionByte), "$RETURN"), bracketsContentRegexp, global.SyntaxGo)
		utility.TestOutput("function analysis = %+v", functionAnalysis)
		goFileAnalysis.FunctionMap[functionAnalysis.Name] = functionAnalysis
		utility.TestOutput("---------------- splitter ----------------")
	}

	// 解析函数体
	// functionDefinitionIndexList := functionDefinitionRegexp.FindAllIndex(fileContentByte, -1)
	// for index, functionDefinitionIndex := range functionDefinitionIndexList {
	// 	definitionLength := utility.CalculatePunctuationMarksContentLength(string(fileContentByte[functionDefinitionIndex[1]+1:]), global.PunctuationMarkCurlyBraces)
	// 	if definitionLength == -1 {
	// 		ui.OutputWarnInfo(ui.CMDAnalyzeGoFunctionContentSyntax)
	// 		continue
	// 	}
	// 	utility.TestOutput("function[%v] end index = %v", index, functionDefinitionIndex[1]+1+definitionLength)
	// 	functionContent := fileContentByte[functionDefinitionIndex[0] : functionDefinitionIndex[1]+1+definitionLength]
	// 	utility.TestOutput("function[%v] functionContent = |%v|", index, string(functionContent))
	// 	utility.TestOutput("---------------- splitter ----------------")
	// }
}

// 分析 CPP 文件

func analyzeCppFile(toAnalyzeFile, toWriteFile *os.File) error {
	return nil
}

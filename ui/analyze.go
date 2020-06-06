package ui

// CMDAnalyzeFileOrDirectoryNotExist 待分析的文件不存在
var CMDAnalyzeFileOrDirectoryNotExist string = "File or Directory %v to analyze does not exist"

// CMDAnalyzeOccursError 分析出错
var CMDAnalyzeOccursError string = "File to analyze occurs error: %v"

// CMDAnalyzeGoKeywordRegexpNotExist 分析 go 文件，关键词的解析式不存在
var CMDAnalyzeGoKeywordRegexpNotExist string = "Analyze Go file, keyword %v regexp does not exist"

// CMDAnalyzeGoFunctionDefinitionSyntaxError 分析 go 文件，函数定义语法错误
var CMDAnalyzeGoFunctionDefinitionSyntaxError string = "Analyze Go file, function definition syntax error"

// CMDAnalyzeGoFunctionContentSyntaxError 分析 go 文件，函数体语法错误
var CMDAnalyzeGoFunctionContentSyntaxError string = "Analyze Go file, function content syntax error"

// CMDAnalyzeGoFunctionError 分析 go 文件，函数体语法错误
var CMDAnalyzeGoFunctionError string = "Analyze Go file, function definition does not match body on index %v"

// CMDAnalyzeGoFileAnalysisNotExist 分析 go 项目，文件的分析结果不存在
var CMDAnalyzeGoFileAnalysisNotExist string = "Analyze Go project, file No %v does not have analysis"

// CMDAnalyzeGoPackageAnalysisNotExist 分析 go 项目，包的分析结果不存在
var CMDAnalyzeGoPackageAnalysisNotExist string = "Analyze Go project, package %v does not have analysis"

// AnalyzeGoFileResultTemplate 分析 go 文件，输出结果的文本模板
var AnalyzeGoFileResultTemplate = `
## File: RP_FILE_PATH
RP_PACKAGE_CONTENT
RP_IMPORT_PACKAGE_LIST
RP_FUNCTION_DEFINITION_LIST
`

// AnalyzeGoFilePackageContentTemplate 分析 go 文件，文件定义的包的文本模板
var AnalyzeGoFilePackageContentTemplate = `
- Package:
(9,1)- Name: RP_PACKAGE_NAME
(9,1)- Path: RP_PACKAGE_PATH
`

// AnalyzeGoFileImportPackageListTemplate 分析 go 文件，依赖包列表的文本模板
var AnalyzeGoFileImportPackageListTemplate = `
- Import
RP_IMPORT_PACKAGE
`

// AnalyzeGoFileImportPackageTemplate 分析 go 文件，依赖包的文本模板
var AnalyzeGoFileImportPackageTemplate = `(9,1)- RP_PACKAGE_ALIAS: RP_PACKAGE_PATH`

// AnalyzeGoFileFunctionDefinitionListTemplate 分析 go 文件，函数列表的文本模板
var AnalyzeGoFileFunctionDefinitionListTemplate = `
- Function
RP_FUNCTION_DEFINITION
`

// AnalyzeGoFileFunctionDefinitionTemplate 分析 go 文件，函数的定义的文本模板
var AnalyzeGoFileFunctionDefinitionTemplate = `
(9,1)- RP_FUNCTION_NAME
(9,2)RP_FUNCTION_CLASS
RP_FUNCTION_PARAM_LIST
RP_FUNCTION_RETURN_LIST
RP_FUNCTION_CALL_MAP
`

// AnalyzeGoFileFunctionClassTemplate 分析 go 文件，函数的定义的类的文本模板
var AnalyzeGoFileFunctionClassTemplate = `- Class: RP_FUNCTION_CLASS_NAME`

// AnalyzeGoFileFunctionParamListTemplate 分析 go 文件，函数的定义的参数表的文本模板
var AnalyzeGoFileFunctionParamListTemplate = `
(9,2)- Params
RP_FUNCTION_PARAM_NAME_TYPE_LIST
`

// AnalyzeGoFileFunctionParamNameTypeTemplate 分析 go 文件，函数定义中的名称，类型的文本模板
var AnalyzeGoFileFunctionParamNameTypeTemplate = `(9,3)- RP_NAME: RP_TYPE`

// AnalyzeGoFileFunctionReturnListTemplate 分析 go 文件，函数的返回值列表的文本模板
var AnalyzeGoFileFunctionReturnListTemplate = `
(9,2)- Return
RP_FUNCTION_RETURN_NAME_TYPE_LIST
`

// AnalyzeGoFileFunctionReturnNameTypeTemplate 分析 go 文件，函数定义中的名称，类型的文本模板
var AnalyzeGoFileFunctionReturnNameTypeTemplate = `(9,3)- RP_NAME: RP_TYPE`

// AnalyzeGoFileFunctionCallMapTemplate 分析 go 文件，函数体内调用其他函数的列表的文本模板
var AnalyzeGoFileFunctionCallMapTemplate = `
(9,2)- Call
RP_FUNCTION_CALL_PACKAGE_MAP
`

// AnalyzeGoFileFunctionCallPackageMap 分析 go 文件，函数体内调用其他函数的文本模板
var AnalyzeGoFileFunctionCallPackageMap = `
(9,3)- RP_PACKAGE: RP_IDENTIFIER
`

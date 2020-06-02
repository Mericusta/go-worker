package ui

// CMDAnalyzeFileNotExist 待分析的文件不存在
var CMDAnalyzeFileNotExist string = "File to analyze does not exist"

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

// AnalyzeGoFileResultTemplate 分析 go 文件，输出结果的文本模板
var AnalyzeGoFileResultTemplate = `
## File: RP_FILE_PATH
- Package: RP_PACKAGE_NAME
RP_IMPORT_PACKAGE_LIST
RP_FUNCTION_DEFINITION_LIST
`

// AnalyzeGoFileImportPackageListTemplate 分析 go 文件，依赖包列表的文本模板
var AnalyzeGoFileImportPackageListTemplate = `
- Import
RP_IMPORT_PACKAGE
`

// AnalyzeGoFileImportPackageTemplate 分析 go 文件，依赖包的文本模板
var AnalyzeGoFileImportPackageTemplate = `(9,1)- PACKAGE_PATH`

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

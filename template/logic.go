package template

// TEOptionParentValue 指令选项 parent 的模板表达式
var TEOptionParentValue TemplateExpression = `parent(\s+TK_PATH){1}`

// TEOptionOutputValue 指令选项 output 的模板表达式
var TEOptionOutputValue TemplateExpression = `output(\s+TK_PATH){1}`

// TEGoKeywordImportValue go import 关键词的模板表达式
var TEGoKeywordImportValue TemplateExpression = `(?ms)^import\s+(?P<VALUE>\(.*?\)|TK_GoKeywordImportAlias)`

// TEGoKeywordImportAlias go import 关键词的重命名形式的模板表达式
var TEGoKeywordImportAlias TemplateExpression = `(?P<ALIAS>\w+\s+)??TK_DoubleQuotesContent`

// TEGoFunctionDefinition go function 定义的模板表达式
var TEGoFunctionDefinition TemplateExpression = `(?ms)^func\s*(?P<MEMBER>\(.*?\))?\s*(?P<NAME>\w+)\s*(?P<PARAM>\(.*?\)){1}\s*(?P<RETURN>\(.*?\)|[\.\*\w]+)?\s*\{`

// TEGoFunctionCall go function 调用的模板表达式
var TEGoFunctionCall TemplateExpression = `(?ms)((?P<CALL>TK_IDENTIFIER)\.)?(?P<NAME>TK_IDENTIFIER)\((?P<PARAM>[^\n]*)\)`

// TEOptionIgnoreValue 指令选项 ignore 的模板表达式
var TEOptionIgnoreValue TemplateExpression = `ignore(\s+TK_PATH){1}`

// TEGoTypeConvert go 内建类型转换的模板表达式
var TEGoTypeConvert TemplateExpression = `^(?P<IDENTIFIER>TK_IDENTIFIER)\(.*\)$`

// TEGoVariableDeclaration go variable 声明的模板表达式
var TEGoVariableDeclaration TemplateExpression = `var\s+(?P<NAME>TK_IDENTIFIER)\s+(?P<TYPE>TK_IDENTIFIER(\.TK_IDENTIFIER)?)`

// TEGoVariableInitialization go variable 初始化的模板表达式
var TEGoVariableInitialization TemplateExpression = `(?P<LIST>[^\n]*?)\s*:=\s*(?P<INIT>.*?)\n`

// TEGoFileSplitterScopePackage Go 语言文件切分器 package 域的模板表达式
var TEGoFileSplitterScopePackage TemplateExpression = `^package\s+(?P<NAME>TK_IDENTIFIER)`

// TEGoFileSplitterScopeImport Go 语言文件切分器 import 域的模板表达式
var TEGoFileSplitterScopeImport TemplateExpression = `^import(\s*\(|\s+TK_GoKeywordImportAlias)`

// TEGoLineImportMultiLineAliasPackage Go 语言 import 多行内包名与路径的表达式
var TEGoLineImportMultiLineAliasPackage TemplateExpression = `^\s*TK_GoKeywordImportAlias`

// // TEGoLineImportOneLine Go 语言 import 单行的模板表达式
// var TEGoLineImportOneLine TemplateExpression = `^(import\s+)?(?P<ALIAS>\w+\s+)??TK_DoubleQuotesContent`

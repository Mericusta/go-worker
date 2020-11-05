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

// TEGoVariableInitialization go variable 初始化的模板表达式
var TEGoVariableInitialization TemplateExpression = `(?P<LIST>[^\n]*?)\s*:=\s*(?P<INIT>.*?)\n`

// TEGoFileSplitterScopePackage Go 语言文件切分器 package 域的模板表达式
var TEGoFileSplitterScopePackage TemplateExpression = `^package\s+(?P<NAME>TK_IDENTIFIER)`

// TEGoFileAnalyzerScopePackage Go 语言文件分析器 package 域的模板表达式
var TEGoFileAnalyzerScopePackage TemplateExpression = TEGoFileSplitterScopePackage

// TEGoFileSplitterScopeMultiLineImportStart Go 语言文件切分器多行 import 域的起始的模板表达式
var TEGoFileSplitterScopeMultiLineImportStart TemplateExpression = `^import\s*\(`

// TEGoFileSplitterScopeMultiLineImportContent Go 语言文件切分器多行 import 域的内容的模板表达式
var TEGoFileSplitterScopeMultiLineImportContent TemplateExpression = `^\s*TK_GoKeywordImportAlias`

// TEGoFileAnalyzerScopeImportContent Go 语言文件分析器 import 域的内容的模板表达式
var TEGoFileAnalyzerScopeImportContent TemplateExpression = `^\s*TK_GoKeywordImportAlias`

// TEGoFileSplitterScopeSingleLineImport Go 语言文件切分器单行 import 域的模板表达式
var TEGoFileSplitterScopeSingleLineImport TemplateExpression = `^import\s+TK_GoKeywordImportAlias`

// TEGoFileSplitterScopePackageVariable Go 语言文件切分器包变量域的模板表达式
var TEGoFileSplitterScopePackageVariable TemplateExpression = `^var\s+(?P<NAME>TK_IDENTIFIER)\s+(?P<TYPE>TK_IDENTIFIER(\.TK_IDENTIFIER)?)`

// TEGoFileAnalyzerScopePackageVariable Go 语言文件分析器包变量域的模板表达式
var TEGoFileAnalyzerScopePackageVariable TemplateExpression = TEGoFileSplitterScopePackageVariable

// TEGoFileSplitterScopeInterface Go 语言文件切分器接口域的模板表达式
var TEGoFileSplitterScopeInterface TemplateExpression = `^type\s+(?P<NAME>TK_IDENTIFIER)\s+interface\s*\{(.*(?P<SCOPE_END>\}))?`

// TEGoFileAnalyzerScopeInterfaceFunction Go 语言文件分析器接口域的函数的模板表达式
var TEGoFileAnalyzerScopeInterfaceFunction TemplateExpression = `\s*(?P<NAME>TK_IDENTIFIER)(\((?P<PARAM>.*?)\)){1}\s*(?P<RETURN>\(.*?\)|[\.\*\w]+)?`

// TEGoFileAnalyzerScopeInterface Go 语言文件分析器接口域的内容的模板表达式
var TEGoFileAnalyzerScopeInterface TemplateExpression = `(?ms)^type\s+(?P<NAME>TK_IDENTIFIER)\s+interface\s*\{(?P<BODY>.*)\}`

// TEGoFileSplitterScopeStruct Go 语言文件切分器结构体域的模板表达式
var TEGoFileSplitterScopeStruct TemplateExpression = `^type\s+(?P<NAME>TK_IDENTIFIER)\s+struct\s*\{(?P<CONTENT>[^\}]*)*(?P<SCOPE_END>\})?`

// TEGoFileAnalyzerScopeStructVariable Go 语言文件分析器结构体域的变量的模板表达式
var TEGoFileAnalyzerScopeStructVariable TemplateExpression = `\s*(?P<NAME>TK_IDENTIFIER)(\s+(?P<TYPE>.*))?`

// TEGoFileSplitterScopeFunction Go 语言文件切分器函数域的模板表达式
var TEGoFileSplitterScopeFunction TemplateExpression = `^func\s+(?P<NAME>TK_IDENTIFIER)\s*(?P<DEFINITION>.*)\s*\{(?P<BODY>.*?)(?P<SCOPE_END>\})?$`

// var TEGoFileSplitterScopeFunction TemplateExpression = `^func\s+(?P<NAME>\w+)\s*(?P<PARAM>\(.*?\)){1}\s*(?P<RETURN>\(.*?\)|[\.\*\w]+)?\s*\{(?P<CONTENT>[^\}]*)*(?P<SCOPE_END>\})?(?P<COMMENT>\s*//.*)?`

// TEGoFileAnalyzerScopeFunction Go 语言文件分析器函数域的模板表达式
var TEGoFileAnalyzerScopeFunction TemplateExpression = `(?ms)^func\s+(\((?P<THIS>[^\(\)]*)\))?\s*(?P<NAME>TK_IDENTIFIER)\s*(?P<DEFINITION>.*)\s*\{(?P<BODY>.*)\}$`

// var TEGoFileAnalyzerFunctionType TemplateExpression = `(?ms)^func\s*(?P<DEFINITION>.*)`
// var TEGoFileAnalyzerScopeFunction TemplateExpression = `^func\s+(?P<NAME>\w+)\s*(\((?P<PARAM>.*?)\)){1}\s*(\((?P<RETURN_LIST>.*?)\)|(?P<RETURN_TYPE>TK_IDENTIFIER(\.TK_IDENTIFIER)?))?\s*\{(?P<CONTENT>[^\}]*)*(?P<SCOPE_END>\})?(?P<COMMENT>\s*//.*)?`

// TEGoFileSplitterScopeMemberFunction Go 语言文件切分器成员函数域的模板表达式
var TEGoFileSplitterScopeMemberFunction TemplateExpression = `^func\s*(?P<MEMBER>\(.*?\))?\s*(?P<NAME>\w+)\s*(?P<PARAM>\(.*?\)){1}\s*(?P<RETURN>\(.*?\)|[\.\*\w]+)?\s*\{(?P<CONTENT>[^\}]*)*(?P<SCOPE_END>\})?(?P<COMMENT>\s*//.*)?`

//
// var TEGoFileAnalyzerScopeFunctionCall TemplateExpression = `(?P<FROM>TK_IDENTIFIER\.)?(?P<NAME>TK_IDENTIFIER){1}\(.*`

// TEGoFileSplitterScopeTypeRename Go 语言文件切分器类型重命名域的模板表达式
var TEGoFileSplitterScopeTypeRename TemplateExpression = `^type\s+(?P<NAME>TK_IDENTIFIER)\s+((?P<FROM>TK_IDENTIFIER)\.)?(?P<TYPE>TK_IDENTIFIER)$`

// TEGoFileSplitterScopeMultiLineConstStart Go 语言文件切分器多行 const 域的起始的模板表达式
var TEGoFileSplitterScopeMultiLineConstStart TemplateExpression = `^const\s*\(`

// TEGoFileSplitterScopeMultiLineConstContent Go 语言文件切分器多行 const 域的内容的模板表达式
var TEGoFileSplitterScopeMultiLineConstContent TemplateExpression = `^\s*(?P<NAME>TK_IDENTIFIER)\s*(?P<TYPE>TK_IDENTIFIER)?.*`

// TEGoFileSplitterScopeSingleLineConst Go 语言文件切分器单行 const 域的模板表达式
var TEGoFileSplitterScopeSingleLineConst TemplateExpression = `^const\s+(?P<NAME>TK_IDENTIFIER)\s*(?P<TYPE>TK_IDENTIFIER)?\s*=(?P<VALUE>.*)`

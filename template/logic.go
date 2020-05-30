package template

// TEOptionParentValue 指令选项 parent 的模板表达式
var TEOptionParentValue TemplateExpression = `parent(\s+TK_PATH){1}`

// TEOptionOutputValue 指令选项 output 的模板表达式
var TEOptionOutputValue TemplateExpression = `output(\s+TK_PATH){1}`

// TEGoKeywordImportValue go import 关键词的模板表达式
var TEGoKeywordImportValue TemplateExpression = `(?ms)import\s+(\(.*?\)|TK_DoubleQuotesContent)`

// TEGoFunctionDefinition go function 定义的模板表达式
var TEGoFunctionDefinition TemplateExpression = `(?ms)^func\s*(?P<MEMBER>\(.*?\))?\s*(?P<NAME>\w+)\s*(?P<PARAM>\(.*?\)){1}\s*(?P<RETURN>\(.*?\)|[\*\w]+)?\s*\{`

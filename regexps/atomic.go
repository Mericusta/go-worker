package regexps

// AtomicExpression 原子表达式类型
type AtomicExpression string

// AETemplateCommonKeyword 通用模板关键词表达式，用于匹配模板字符串中的通用模板文本
var AETemplateCommonKeyword AtomicExpression = `TK_[\w_]+`

// AEPath 路径表达式
var AEPath AtomicExpression = `[/_\.\w-]+`

// AEDoubleQuotesContent 双引号的内容的表达式
var AEDoubleQuotesContent AtomicExpression = `(?:")(?P<CONTENT>.*?)(?:")`

// AEBracketsContent 括号的内容的表达式
var AEBracketsContent AtomicExpression = `(?:\()(?ms)(?P<CONTENT>.*?)(?:\))`

// AESquareBracketsContent 方括号的内容的表达式
var AESquareBracketsContent AtomicExpression = `(?:\[)(?P<CONTENT>.*?)(?:\])`

// AECurlyBracesContent 大括号的内容的表达式
var AECurlyBracesContent AtomicExpression = `(?:\{)(?P<CONTENT>.*?)(?:\})`

// AEBindOptionValue 解析 bind 指令参数的表达式
var AEBindOptionValue AtomicExpression = `(project|syntax)(\s+[-:\.~\\/\w]+)?`

// AECreateOptionValue 解析 create 指令参数的表达式
var AECreateOptionValue AtomicExpression = `(package|file)(\s+[_\.\w-]+){1}`

// AEConvertOptionValue 解析 convert 指令参数的表达式
var AEConvertOptionValue AtomicExpression = `csv(\s+[_\w-]+){1}`

// AEConvertACOptionValue 解析 convert 指令参数的表达式
var AEConvertACOptionValue AtomicExpression = `(create|append)(\s+[_\.\w-]+){1}`

// AEAnalyzeOptionValue 解析 analyze 指令参数的表达式
var AEAnalyzeOptionValue AtomicExpression = `(file|directory)(\s+[_\.\w-]+)+`

// AEGoKeywordPackageValue go package 关键词的表达式
var AEGoKeywordPackageValue AtomicExpression = `^package\s+[_\w-]+`

// AETemplateStyle 格式模板关键词表达式，用于匹配文本中指定的格式模板
var AETemplateStyle AtomicExpression = `\((?P<CHAR>\d+),(?P<NUM>\d+)\)`

// AESpaceLine 空白行的表达式
var AESpaceLine AtomicExpression = `(?ms)\n\s*\n`

// AEIdentifier 标识符的表达式，以字母开头，字母，数字或下划线组合
var AEIdentifier AtomicExpression = `[[:alpha:]][_\.\w]*`

// AEFileNameType 文件名与类型的表达式
var AEFileNameType AtomicExpression = `(?P<NAME>.*)\.(?P<TYPE>\w+)$`

// AERemoveOptionValue 解析 remove 指令参数的表达式
var AERemoveOptionValue AtomicExpression = `(file|type)(\s+[_\.\w-]+)+`

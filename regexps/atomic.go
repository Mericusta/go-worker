package regexps

// AtomicExpression 原子表达式类型
type AtomicExpression string

// AETemplateCommonKeyword 通用模板关键词表达式，用于匹配模板字符串中的通用模板文本
var AETemplateCommonKeyword AtomicExpression = `TK_[\w_]+`

// AEPath 路径表达式
var AEPath AtomicExpression = `[/_\.\w-]+`

// AEBindOptionValue 解析 bind 指令参数的表达式
var AEBindOptionValue AtomicExpression = `(project|syntax)(\s+[-:\.~\\/\w]+)?`

// AECreateOptionValue 解析 create 指令参数的表达式
var AECreateOptionValue AtomicExpression = `(package|file)(\s+[_\.\w-]+){1}`

// AEConvertOptionValue 解析 convert 指令参数的表达式
var AEConvertOptionValue AtomicExpression = `csv(\s+[_\w-]+){1}`

// AEConvertACOptionValue 解析 convert 指令参数的表达式
var AEConvertACOptionValue AtomicExpression = `(create|append)(\s+[_\.\w-]+){1}`

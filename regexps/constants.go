package regexps

// 内部逻辑表达式

// AETemplateCommonKeyword 通用模板关键词表达式，用于匹配模板字符串中的通用模板文本
var AETemplateCommonKeyword Expression = `\$[\w_]+`

// AEPath 路径表达式
var AEPath Expression = `[/_\.\w-]+`

// 指令解析辅助表达式

// AEBindOptionValue 解析 bind 指令参数的表达式
var AEBindOptionValue Expression = `(project|syntax)(\s+[-:\.~\\/\w]+)`

// 模板匹配式

// 原子匹配式

// MTEPath 路径的模板匹配式
var MTEPath Expression = `\$PATH`

// 组合匹配式

// MTEOptionParentValue 指令选项 parent 的模板匹配式
var MTEOptionParentValue Expression = `\$OPVExpression`

// MTEOptionOutputValue 指令选项 output 的模板匹配式
var MTEOptionOutputValue Expression = `\$OOVExpression`

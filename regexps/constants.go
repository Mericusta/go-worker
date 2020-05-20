package regexps

// 指令表达式

// ExpExit 解析指令 exit
var ExpExit Expression = `^exit$`

// ExpBind 解析指令 bind [operation] [option] [value]
var ExpBind Expression = `^bind\s+(project|syntax)(\s+[-:\.~\\/\w]+)?$`

// ExpCreate 解析指令 create [option] [value] [option value]
var ExpCreate Expression = `^create\s+(package|file)(\s+[_\.\w-]+){1}(\s+parent(\s+[_\.\w-]+){1})?$`

// ExpConvert 解析指令 convert [option] [value] [option value] [option value]
var ExpConvert Expression = `^convert\s+csv(\s+[_\w-]+){1}(\s+parent(\s+[_\.\w-]+){1})?(\s+(create|append)(\s+[_\.\w-]+){1})?$`

// 内部逻辑表达式

// ExpTemplateKeyword 模板关键词表达式
var ExpTemplateKeyword Expression = `\$[\w_]+`

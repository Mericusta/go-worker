package regexps

// 外部逻辑表达式

// 指令表达式

// AECmdExit 解析指令 exit
var AECmdExit Expression = `^exit$`

// AECmdBind 解析指令 bind [operation] [option] [value]
var AECmdBind Expression = `^bind\s+(project|syntax)(\s+[-:\.~\\/\w]+)?$`

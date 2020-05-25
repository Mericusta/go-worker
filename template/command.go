package template

// TEBind 指令 bind 的模板表达式
var TEBind TemplateExpression = `^bind\s+$BOVExpression$`

// TECreate 指令 create 的模板表达式
var TECreate TemplateExpression = `^create\s+(package|file)(\s+[_\.\w-]+){1}(\s+$OPVExpression)?$`

// TEConvert 指令 convert 的模板表达式
var TEConvert TemplateExpression = `^convert\s+csv(\s+[_\w-]+){1}(\s+$OPVExpression)?(\s+(create|append)(\s+[_\.\w-]+){1})?$`

// TEAnalyze 指令 analyze 的模板表达式
var TEAnalyze TemplateExpression = `^analyze\s+(file|directory|package)(\s+[_\.\w-]+)+(\s+$OPVExpression)?(\s+$OOVExpression)?$`

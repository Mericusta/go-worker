package template

// TECommandBind 指令 bind 的模板表达式
var TECommandBind TemplateExpression = `^bind\s+TK_BOVExpression$`

// TECommandCreate 指令 create 的模板表达式
var TECommandCreate TemplateExpression = `^create\s+TK_CreateOVExpression(\s+TK_OPVExpression)?$`

// TECommandConvert 指令 convert 的模板表达式
var TECommandConvert TemplateExpression = `^convert\s+TK_ConvertOVExpression(\s+TK_OPVExpression)?(\s+(create|append)(\s+[_\.\w-]+){1})?$`

// TECommandAnalyze 指令 analyze 的模板表达式
var TECommandAnalyze TemplateExpression = `^analyze\s+(file|directory|package)(\s+[_\.\w-]+)+(\s+TK_OPVExpression)?(\s+TK_OOVExpression)?$`

package template

// TECommandBind 指令 bind 的模板表达式
var TECommandBind TemplateExpression = `^bind\s+TK_BOVExpression$`

// TECommandCreate 指令 create 的模板表达式
var TECommandCreate TemplateExpression = `^create\s+TK_CreateOVExpression(\s+TK_OPVExpression)?$`

// TECommandConvert 指令 convert 的模板表达式
var TECommandConvert TemplateExpression = `^convert\s+TK_ConvertOVExpression(\s+TK_OPVExpression)?(\s+TK_ConvertACOptionExpression)?$`

// TECommandAnalyze 指令 analyze 的模板表达式
// var TECommandAnalyze TemplateExpression = `^analyze\s+TK_AnalyzeOVExpression(\s+TK_OPVExpression)?(\s+TK_OOVExpression)?$`
var TECommandAnalyze TemplateExpression = `^analyze\s+(?P<PATH>TK_PATH)?$`

// TECommandRemove 指令 remove 的模板表达式
var TECommandRemove TemplateExpression = `^remove\s+TK_RemoveOVExpression(\s+TK_OPVExpression)?(\s+TK_OIVExpression)?$`

// TECommandTemplater 指令 templater 的模板表达式
var TECommandTemplater TemplateExpression = `^templater\s+TK_PATH(\s+TK_OPVExpression)?(\s+TK_OOVExpression)?$`

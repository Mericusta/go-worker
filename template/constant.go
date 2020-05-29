package template

// TemplateKeyword 模板关键词
type TemplateKeyword string

var (
	// 通用关键词
	// TK_Path path 的关键词
	TK_Path TemplateKeyword = "TK_PATH"

	// 通用选项关键词
	// TK_OPVExpression 选项 parent 的关键词
	TK_OPVExpression TemplateKeyword = "TK_OPVExpression"
	// TK_OOVExpression 选项 output 的关键词
	TK_OOVExpression TemplateKeyword = "TK_OOVExpression"
	// TK_DoubleQuotesContent 双引号的内容的关键词
	TK_DoubleQuotesContent TemplateKeyword = "TK_DoubleQuotesContent"

	// 指令特有关键词
	// TK_BOVExpression bind 指令的 project|syntax 选项关键词
	TK_BOVExpression TemplateKeyword = "TK_BOVExpression"
	// TK_CreateOVExpression create 指令的 package|file 选项关键词
	TK_CreateOVExpression TemplateKeyword = "TK_CreateOVExpression"
	// TK_ConvertOVExpression convert 指令的 csv 选项关键词
	TK_ConvertOVExpression TemplateKeyword = "TK_ConvertOVExpression"
	// TK_ConvertACOptionExpression convert 指令的 append|create 选项关键词
	TK_ConvertACOptionExpression TemplateKeyword = "TK_ConvertACOptionExpression"
	// TK_AnalyzeOVExpression analyze 指令的 file|directory|package 选项关键词
	TK_AnalyzeOVExpression TemplateKeyword = "TK_AnalyzeOVExpression"
)

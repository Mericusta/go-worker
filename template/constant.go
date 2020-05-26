package template

// TemplateKeyword 模板关键词
type TemplateKeyword string

var (
	// TK_BOVExpression bind 指令的 option 关键词
	TK_BOVExpression TemplateKeyword = "TK_BOVExpression"
	// TK_CreateOVExpression create 指令的 option 关键词
	TK_CreateOVExpression TemplateKeyword = "TK_CreateOVExpression"
	// TK_ConvertOVExpression convert 指令的 option 关键词
	TK_ConvertOVExpression TemplateKeyword = "TK_ConvertOVExpression"
	// TK_Path path 的关键词
	TK_Path TemplateKeyword = "TK_PATH"

	// TK_OPVExpression 选项 parent 的关键词
	TK_OPVExpression TemplateKeyword = "TK_OPVExpression"
	// TK_OOVExpression 选项 output 的关键词
	TK_OOVExpression TemplateKeyword = "TK_OOVExpression"
)

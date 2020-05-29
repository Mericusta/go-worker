package global

// TemplateEnum 模板枚举类型
type TemplateEnum int

const (
	// 内部模板

	// OptionParentValueTemplate option 选项模板
	OptionParentValueTemplate TemplateEnum = 1
	// OptionOutputValueTemplate output 选项模板
	OptionOutputValueTemplate TemplateEnum = 2
	// GoKeywordImportValueTemplate go import 关键词的模板
	GoKeywordImportValueTemplate TemplateEnum = 3
	// GoFunctionDefinitionTemplate go function 定义的模板
	GoFunctionDefinitionTemplate TemplateEnum = 4

	// 外部模板

	// CommandBindTemplate bind 指令模板
	CommandBindTemplate TemplateEnum = 101
	// CommandCreateTemplate create 指令模板
	CommandCreateTemplate TemplateEnum = 102
	// CommandConverTemplate convert 指令模板
	CommandConverTemplate TemplateEnum = 103
)

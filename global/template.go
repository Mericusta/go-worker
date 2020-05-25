package global

// TemplateEnum 模板枚举类型
type TemplateEnum int

// 内部模板
const (
	// CommonKeywordTemplate 通用关键词模板
	CommonKeywordTemplate TemplateEnum = 1
	// PathTemplate 路径模板
	PathTemplate TemplateEnum = 2
	// OptionParentValueTemplate option 选项模板
	OptionParentValueTemplate TemplateEnum = 3
	// OptionOutputValueTemplate output 选项模板
	OptionOutputValueTemplate TemplateEnum = 4
	// BindOptionValueTemplate bind 选项模板
	BindOptionValueTemplate TemplateEnum = 5
	// CreateOptionValueTemplate create 选项模板
	CreateOptionValueTemplate TemplateEnum = 6
	// ConvertOptionValueTemplate convert 选项模板
	ConvertOptionValueTemplate TemplateEnum = 7
)

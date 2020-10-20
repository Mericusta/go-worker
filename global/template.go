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
	// GoKeywordImportAliasTemplate go import 关键词的重命名形式的模板表达式
	GoKeywordImportAliasTemplate TemplateEnum = 4
	// GoFunctionDefinitionTemplate go function 定义的模板
	GoFunctionDefinitionTemplate TemplateEnum = 5
	// GoFunctionCallTemplate go function 调用的模板
	GoFunctionCallTemplate TemplateEnum = 6
	// OptionIgnoreValueTemplate ignore 选项模板
	OptionIgnoreValueTemplate TemplateEnum = 7
	// GoTypeConvertTemplate go 内建类型转换的模板
	GoTypeConvertTemplate TemplateEnum = 8
	// GoVariableDeclarationTemplate go variable 声明的模板
	GoVariableDeclarationTemplate TemplateEnum = 9
	// GoVariableInitializationTemplate go variable 初始化的模板
	GoVariableInitializationTemplate TemplateEnum = 10

	// GoFileSplitterScopePackageTemplate Go 语言文件切分器 package 域的模板表达式
	GoFileSplitterScopePackageTemplate TemplateEnum = 11
	// GoFileSplitterScopeImportTemplate Go 语言文件切分器 import 域的模板表达式
	GoFileSplitterScopeImportTemplate TemplateEnum = 12

	// // GoFileSplitterScopeImportTemplate Go 语言 import 多行内包名与路径的表达式
	// GoFileSplitterScopeImportTemplate TemplateEnum = 12
	// // GoLineImportOneLine Go 语言 import 单行的模板表达式
	// GoLineImportOneLine TemplateEnum = 13

	// 外部模板

	// CommandBindTemplate bind 指令模板
	CommandBindTemplate TemplateEnum = 101
	// CommandCreateTemplate create 指令模板
	CommandCreateTemplate TemplateEnum = 102
	// CommandConverTemplate convert 指令模板
	CommandConverTemplate TemplateEnum = 103
	// CommandRemoveTemplate remove 指令模板
	CommandRemoveTemplate TemplateEnum = 104
)

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
	// GoVariableInitializationTemplate go variable 初始化的模板
	GoVariableInitializationTemplate TemplateEnum = 10

	// GoFileSplitterScopePackageTemplate Go 语言文件切分器 package 域的模板表达式
	GoFileSplitterScopePackageTemplate TemplateEnum = 11
	// GoFileSplitterScopeMultiLineImportStartTemplate Go 语言文件切分器多行 import 域的起始的模板表达式
	GoFileSplitterScopeMultiLineImportStartTemplate TemplateEnum = 12
	// GoFileSplitterScopeMultiLineImportContentTemplate Go 语言文件切分器多行 import 域的内容的模板表达式
	GoFileSplitterScopeMultiLineImportContentTemplate TemplateEnum = 13
	// GoFileSplitterScopeSingleLineImportTemplate Go 语言文件切分器单行 import 域的模板表达式
	GoFileSplitterScopeSingleLineImportTemplate TemplateEnum = 14
	// GoFileSplitterScopePackageVariableTemplate Go 语言文件切分器包变量域的模板表达式
	GoFileSplitterScopePackageVariableTemplate TemplateEnum = 15
	// GoFileSplitterScopeInterfaceTemplate Go 语言文件切分器接口域的模板表达式
	GoFileSplitterScopeInterfaceTemplate TemplateEnum = 16
	// GoFileSplitterScopeStructTemplate Go 语言文件切分器结构体域的模板表达式
	GoFileSplitterScopeStructTemplate TemplateEnum = 17
	// GoFileSplitterScopeFunctionTemplate Go 语言文件切分器函数域的模板表达式
	GoFileSplitterScopeFunctionTemplate TemplateEnum = 18
	// GoFileSplitterScopeMemberFunctionTemplate Go 语言文件切分器成员函数域的模板表达式
	GoFileSplitterScopeMemberFunctionTemplate TemplateEnum = 19
	// GoFileSplitterScopeTypeRenameTemplate Go 语言文件切分器类型重命名域的模板表达式
	GoFileSplitterScopeTypeRenameTemplate TemplateEnum = 20
	// GoFileSplitterScopeMultiLineConstStartTemplate Go 语言文件切分器多行 const 域的起始的模板表达式
	GoFileSplitterScopeMultiLineConstStartTemplate TemplateEnum = 21
	// GoFileSplitterScopeMultiLineConstContentTemplate Go 语言文件切分器多行 const 域的内容的模板表达式
	GoFileSplitterScopeMultiLineConstContentTemplate TemplateEnum = 22
	// GoFileSplitterScopeSingleLineConstTemplate Go 语言文件切分器单行 const 域的模板表达式
	GoFileSplitterScopeSingleLineConstTemplate TemplateEnum = 23

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

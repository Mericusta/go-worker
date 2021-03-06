package template

import "github.com/go-worker/global"

// TemplateString 模板字符串类型
type TemplateString string

// TemplateExpression 模板表达式类型
type TemplateExpression string

// CommandTemplateExpressionMap 指令枚举与模板表达式的映射
var CommandTemplateExpressionMap map[global.CommandEnum]TemplateExpression

// TemplateExpressionMap 模板枚举与模板表达式的映射
var TemplateExpressionMap map[global.TemplateEnum]TemplateExpression

func init() {
	CommandTemplateExpressionMap = map[global.CommandEnum]TemplateExpression{
		global.CmdBind:      TECommandBind,
		global.CmdCreate:    TECommandCreate,
		global.CmdConvert:   TECommandConvert,
		global.CmdAnalyze:   TECommandAnalyze,
		global.CmdRemove:    TECommandRemove,
		global.CmdTemplater: TECommandTemplater,
	}
	TemplateExpressionMap = map[global.TemplateEnum]TemplateExpression{
		global.OptionParentValueTemplate:                         TEOptionParentValue,
		global.OptionOutputValueTemplate:                         TEOptionOutputValue,
		global.OptionIgnoreValueTemplate:                         TEOptionIgnoreValue,
		global.GoKeywordImportValueTemplate:                      TEGoKeywordImportValue,
		global.GoKeywordImportAliasTemplate:                      TEGoKeywordImportAlias,
		global.GoFunctionDefinitionTemplate:                      TEGoFunctionDefinition,
		global.GoFunctionCallTemplate:                            TEGoFunctionCall,
		global.GoTypeConvertTemplate:                             TEGoTypeConvert,
		global.GoVariableInitializationTemplate:                  TEGoVariableInitialization,
		global.GoFileSplitterScopePackageTemplate:                TEGoFileSplitterScopePackage,
		global.GoFileSplitterScopeMultiLineImportStartTemplate:   TEGoFileSplitterScopeMultiLineImportStart,
		global.GoFileSplitterScopeMultiLineImportContentTemplate: TEGoFileSplitterScopeMultiLineImportContent,
		global.GoFileSplitterScopeSingleLineImportTemplate:       TEGoFileSplitterScopeSingleLineImport,
		global.GoFileSplitterScopePackageVariableTemplate:        TEGoFileSplitterScopePackageVariable,
		global.GoFileSplitterScopeInterfaceTemplate:              TEGoFileSplitterScopeInterface,
		global.GoFileSplitterScopeStructTemplate:                 TEGoFileSplitterScopeStruct,
		global.GoFileSplitterScopeFunctionTemplate:               TEGoFileSplitterScopeFunction,
		global.GoFileSplitterScopeMemberFunctionTemplate:         TEGoFileSplitterScopeMemberFunction,
		global.GoFileSplitterScopeTypeRenameTemplate:             TEGoFileSplitterScopeTypeRename,
		global.GoFileSplitterScopeMultiLineConstStartTemplate:    TEGoFileSplitterScopeMultiLineConstStart,
		global.GoFileSplitterScopeMultiLineConstContentTemplate:  TEGoFileSplitterScopeMultiLineConstContent,
		global.GoFileSplitterScopeSingleLineConstTemplate:        TEGoFileSplitterScopeSingleLineConst,
		global.GoFileAnalyzerScopePackageTemplate:                TEGoFileAnalyzerScopePackage,
		global.GoFileAnalyzerScopeImportContentTemplate:          TEGoFileAnalyzerScopeImportContent,
		global.GoFileAnalyzerScopePackageVariableTemplate:        TEGoFileAnalyzerScopePackageVariable,
		global.GoFileAnalyzerScopeInterfaceFunctionTemplate:      TEGoFileAnalyzerScopeInterfaceFunction,
		global.GoFileAnalyzerScopeInterfaceTemplate:              TEGoFileAnalyzerScopeInterface,
		global.GoFileAnalyzerScopeStructVariableTemplate:         TEGoFileAnalyzerScopeStructVariable,
		global.GoFileAnalyzerScopeFunctionTemplate:               TEGoFileAnalyzerScopeFunction,
	}
}

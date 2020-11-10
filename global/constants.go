package global

// 内部常量
var (
	LogMarkTest  = "TEST"
	LogMarkNote  = "Note"
	LogMarkError = "Error"
	LogMarkWarn  = "Warn"

	ConvertRPStructName    = "RP_STRUCT_NAME"
	ConvertRPStructMember  = "RP_STRUCT_MEMBER"
	ConvertRPMemberName    = "RP_MEMBER_NAME"
	ConvertRPMemberType    = "RP_MEMBER_TYPE"
	ConvertRPMemberComment = "RP_MEMBER_COMMENT"
	ConvertRPCsvHead       = "RP_CSV_HEAD"

	AnalyzeRPFilePath                   = "RP_FILE_PATH"
	AnalyzeRPPackageName                = "RP_PACKAGE_NAME"
	AnalyzeRPPackageContent             = "RP_PACKAGE_CONTENT"
	AnalyzeRPImportPackageList          = "RP_IMPORT_PACKAGE_LIST"
	AnalyzeRPFunctionDefinitionList     = "RP_FUNCTION_DEFINITION_LIST"
	AnalyzeRPImportPackage              = "RP_IMPORT_PACKAGE"
	AnalyzeRPPackagePath                = "RP_PACKAGE_PATH"
	AnalyzeRPPackageAlias               = "RP_PACKAGE_ALIAS"
	AnalyzeRPFunctionDefinition         = "RP_FUNCTION_DEFINITION"
	AnalyzeRPFunctionName               = "RP_FUNCTION_NAME"
	AnalyzeRPFunctionClass              = "RP_FUNCTION_CLASS"
	AnalyzeRPFunctionParamList          = "RP_FUNCTION_PARAM_LIST"
	AnalyzeRPFunctionReturnList         = "RP_FUNCTION_RETURN_LIST"
	AnalyzeRPFunctionCallMap            = "RP_FUNCTION_CALL_MAP"
	AnalyzeRPFunctionClassName          = "RP_FUNCTION_CLASS_NAME"
	AnalyzeRPFunctionParamNameTypeList  = "RP_FUNCTION_PARAM_NAME_TYPE_LIST"
	AnalyzeRPFunctionReturnNameTypeList = "RP_FUNCTION_RETURN_NAME_TYPE_LIST"
	AnalyzeRPFunctionCallPackageMap     = "RP_FUNCTION_CALL_PACKAGE_MAP"
	AnalyzeRPName                       = "RP_NAME"
	AnalyzeRPType                       = "RP_TYPE"
	AnalyzeRPPackage                    = "RP_PACKAGE"
	AnalyzeRPIdentifier                 = "RP_IDENTIFIER"
	AnalyzeRPEmptyString                = ""
	AnalyzeRPPackageLevel               = "RP_PACKAGE_LEVEL"
	AnalyzeRPPackageList                = "RP_PACKAGE_LIST"

	PunctuationMarkLeftBracket        rune = '('
	PunctuationMarkRightBracket       rune = ')'
	PunctuationMarkLeftCurlyBracket   rune = '{'
	PunctuationMarkRightCurlyBracket  rune = '}'
	PunctuationMarkLeftSquareBracket  rune = '['
	PunctuationMarkRightSquareBracket rune = ']'
	PunctuationMarkLeftDoubleQuotes   rune = '"'
	PunctuationMarkRightDoubleQuotes  rune = '"'
	PunctuationMarkLeftInverseQuotes  rune = '`'
	PunctuationMarkRightInverseQuotes rune = '`'
	PunctuationMarkLeftSingleQuotes   rune = '\''
	PunctuationMarkRightSingleQuotes  rune = '\''
	PunctuationMarkPoint              rune = '.'
	ASCIISpace                        rune = ' '

	GoKeywordEmptyInterface            = "interface{}"
	GoKeywordConst                     = "const"
	GoKeywordFunc                      = "func"
	GoKeywordStruct                    = "struct"
	GoSplitterStringEnter              = "\n"
	GoSplitterStringComma              = ","
	GoSplitterStringSpace              = " "
	GoSplitterStringPoint              = "."
	GoAnalyzerScopePunctuationMarkList = []int{
		PunctuationMarkBracket,
		PunctuationMarkCurlyBracket,
		PunctuationMarkSquareBracket,
	}
	GoAnalyzerInvalidScopePunctuationMarkMap = map[rune]rune{
		PunctuationMarkLeftDoubleQuotes:  PunctuationMarkRightDoubleQuotes,
		PunctuationMarkLeftInverseQuotes: PunctuationMarkRightInverseQuotes,
		PunctuationMarkLeftSingleQuotes:  PunctuationMarkRightSingleQuotes,
	}
)

// 逻辑常量
const (
	PunctuationMarkQuote         = 1
	PunctuationMarkBracket       = 2
	PunctuationMarkCurlyBracket  = 3
	PunctuationMarkSquareBracket = 4
)

// 配置常量
var (
	SyntaxGo          = "go"
	SyntaxCSV         = "csv"
	SyntaxCpp         = "cpp"
	SyntaxMarkdown    = "md"
	ResourceDirectory = "resources"
)

// 外部常量
var (
	ConfigProjectPathKey    = "path"
	ConfigProjectSyntaxKey  = "syntax"
	ConfigConvertCsvHeadKey = "csv_head"
)

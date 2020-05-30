package global

// 内部常量

// UI 常量
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

	PunctuationMarkLeftQuote        = '"'
	PunctuationMarkRightQuote       = '"'
	PunctuationMarkLeftBracket      = '('
	PunctuationMarkRightBracket     = ')'
	PunctuationMarkLeftCurlyBraces  = '{'
	PunctuationMarkRightCurlyBraces = '}'
)

// 逻辑常量
const (
	PunctuationMarkQuote       = 1
	PunctuationMarkBracket     = 2
	PunctuationMarkCurlyBraces = 3
)

// 配置常量
var (
	SyntaxGo  = "go"
	SyntaxCSV = "csv"
	SyntaxCpp = "cpp"
)

// 外部常量
var (
	ConfigProjectPathKey    = "path"
	ConfigProjectSyntaxKey  = "syntax"
	ConfigConvertCsvHeadKey = "csv_head"
)

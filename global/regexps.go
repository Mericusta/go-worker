package global

// AtomicExpressionEnum 原子表达式枚举类型
type AtomicExpressionEnum int

const (
	// AETemplateCommonKeyword 通用模板关键词的表达式
	AETemplateCommonKeyword AtomicExpressionEnum = 0

	// AEPath 路径的表达式
	AEPath AtomicExpressionEnum = 1
	// AEDoubleQuotesContent 双引号的内容的表达式
	AEDoubleQuotesContent AtomicExpressionEnum = 2
	// AEBracketsContent 括号的内容的表达式
	AEBracketsContent AtomicExpressionEnum = 3
	// AESquareBracketsContent 中括号的内容的表达式
	AESquareBracketsContent AtomicExpressionEnum = 4
	// AECurlyBracesContent 大括号的内容的表达式
	AECurlyBracesContent AtomicExpressionEnum = 5
	// AEGoKeywordPackageValue go package 关键词的表达式
	AEGoKeywordPackageValue AtomicExpressionEnum = 6
	// AETemplateStyle 格式模板关键词表达式，用于匹配文本中指定的格式模板
	AETemplateStyle AtomicExpressionEnum = 7
	// AESpaceLine 空白行的表达式
	AESpaceLine AtomicExpressionEnum = 8
	// AEIdentifier 标识符的表达式
	AEIdentifier AtomicExpressionEnum = 9
	// AEFileNameType 文件名与类型的表达式
	AEFileNameType AtomicExpressionEnum = 10

	// AEBindOptionValue bind 指令参数的表达式
	AEBindOptionValue AtomicExpressionEnum = 101
	// AECreateOptionValue create 指令参数的表达式
	AECreateOptionValue AtomicExpressionEnum = 102
	// AEConvertOptionValue convert 指令参数的表达式
	AEConvertOptionValue AtomicExpressionEnum = 103
	// AEConvertACOptionValue convert 指令参数的表达式
	AEConvertACOptionValue AtomicExpressionEnum = 104
	// AEAnalyzeOptionValue analyze 指令参数的表达式
	AEAnalyzeOptionValue AtomicExpressionEnum = 105
	// AERemoveOptionValue remove 指令参数得表达式
	AERemoveOptionValue AtomicExpressionEnum = 106
)

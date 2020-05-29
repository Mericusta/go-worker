package global

// AtomicExpressionEnum 原子表达式枚举类型
type AtomicExpressionEnum int

const (
	// AETemplateCommonKeyword 通用模板关键词的表达式
	AETemplateCommonKeyword AtomicExpressionEnum = 0

	// AEPath 路径的表达式
	AEPath AtomicExpressionEnum = 1
	// AEDoubleQuotesContent 双引号的表达式
	AEDoubleQuotesContent AtomicExpressionEnum = 2
	// AEGoKeywordPackageValue go package 关键词的表达式
	AEGoKeywordPackageValue AtomicExpressionEnum = 3

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
)

package global

// AtomicExpressionEnum 原子表达式枚举类型
type AtomicExpressionEnum int

const (
	// AETemplateCommonKeyword 通用模板关键词表达式
	AETemplateCommonKeyword AtomicExpressionEnum = 1
	// AEPath 路径表达式
	AEPath AtomicExpressionEnum = 2
	// AEBindOptionValue bind 指令参数的表达式
	AEBindOptionValue AtomicExpressionEnum = 3
	// AECreateOptionValue create 指令参数的表达式
	AECreateOptionValue AtomicExpressionEnum = 4
	// AEConvertOptionValue convert 指令参数的表达式
	AEConvertOptionValue AtomicExpressionEnum = 5
)

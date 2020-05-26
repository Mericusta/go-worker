package global

// AtomicRegexpEnum 原子表达式枚举类型
type AtomicRegexpEnum int

const (
	// AETemplateCommonKeyword 通用模板关键词表达式
	AETemplateCommonKeyword AtomicRegexpEnum = 1
	// AEPath 路径表达式
	AEPath AtomicRegexpEnum = 2
	// AEBindOptionValue bind 指令参数的表达式
	AEBindOptionValue AtomicRegexpEnum = 3
	// AECreateOptionValue create 指令参数的表达式
	AECreateOptionValue AtomicRegexpEnum = 4
	// AEConvertOptionValue convert 指令参数的表达式
	AEConvertOptionValue AtomicRegexpEnum = 5
)

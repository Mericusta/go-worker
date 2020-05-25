package global

// RegexpEnum 正则表达式枚举类型
type RegexpEnum int

const (
	// AETemplateCommonKeyword 通用模板关键词表达式
	AETemplateCommonKeyword RegexpEnum = 1
	// AEPath 路径表达式
	AEPath RegexpEnum = 2
	// AEBindOptionValue bind 指令参数的表达式
	AEBindOptionValue RegexpEnum = 3
	// AECreateOptionValue create 指令参数的表达式
	AECreateOptionValue RegexpEnum = 4
	// AEConvertOptionValue convert 指令参数的表达式
	AEConvertOptionValue RegexpEnum = 5
)

package global

// CommandEnum 指令枚举类型
type CommandEnum int

const (
	// CmdUnknown 未知指令
	CmdUnknown CommandEnum = -1
	// CmdExit 退出指令
	CmdExit CommandEnum = 1
	// CmdBind 绑定指令
	CmdBind CommandEnum = 2
	// CmdCreate 创建指令
	CmdCreate CommandEnum = 3
	// CmdConvert 转化指令
	CmdConvert CommandEnum = 4
	// CmdAnalyze 分析指令
	CmdAnalyze CommandEnum = 5
)

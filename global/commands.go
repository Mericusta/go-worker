package global

// CommandEnum 指令枚举类型
type CommandEnum int

const (
	// CmdUnknown 未知指令
	CmdUnknown CommandEnum = -1
	// CmdCustom 自定义指令
	CmdCustom CommandEnum = 0
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
	// CmdRemove 删除指令
	CmdRemove CommandEnum = 6
	// CmdSpider 爬虫指令
	CmdSpider CommandEnum = 7
)

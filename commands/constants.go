package commands

const (
	// Unknown 未知指令
	Unknown Command = -1
	// Exit 退出指令
	Exit Command = 1
	// Create 创建表格指令
	Create Command = 2
	// Insert 插入表格数据
	Insert Command = 3
	// Spider 爬虫工具：BLHX 数据爬取
	Spider Command = 4
	// Simulator 模拟器：MHW 配装模拟器
	Simulator Command = 5
	// Worker 辅助工作器
	Worker Command = 6
	// Alpha 阿尔法机
	Alpha Command = 7
)
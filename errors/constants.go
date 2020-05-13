package errors

const (
	// Undefined 未定义错误
	Undefined Error = -1
	// UnknownCommand 未知指令
	UnknownCommand Error = 1
	// UnknownFSMState 未知状态
	UnknownFSMState Error = 2
	// CreateFSMState 创建表格状态
	CreateFSMState Error = 3
	// InsertFSMState 插入数据状态
	InsertFSMState Error = 4
	// NoneCommandExecutor 指令执行方法不存在
	NoneCommandExecutor Error = 5
)
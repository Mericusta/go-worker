package fsm

const (
	// Unknown 未知状态
	Unknown FSMState = -1
	// Error 错误状态
	Error FSMState = 1
	// Waitting 等待指令状态
	Waitting FSMState = 2
	// Exiting 退出状态
	Exiting FSMState = 3
	// Executing 执行指令状态
	Executing FSMState = 4
)

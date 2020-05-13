package commands

import "fmt"

// Command 指令枚举类型
type Command int

// CommandList 指令枚举列表
var CommandList []Command

// CommandExecutorList 指令执行列表
var CommandExecutorList map[Command]func(string) error

func init() {
	CommandList = []Command{
		Exit,
		Worker,
	}
	CommandExecutorList = map[Command]func(string) error{
		Exit:   ExitExecutor,
		Worker: WorkerExecutor,
	}
}

func Execute(inputString string, command Command) error {
	executor, hasExecutor := CommandExecutorList[command]
	if !hasExecutor {
		return fmt.Errorf("command[%v] does not have executor", command)
	}
	if executor != nil {
		return executor(inputString)
	}
	return nil
}

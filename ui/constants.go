package ui

// 状态机文本

// FSMWaitting 等待指令
var FSMWaitting string = "Waitting Orders: "

// FSMUnknownCommand 未知指令
var FSMUnknownCommand string = "Unknown command"

// CommonError1 通用错误显示文本1：参数不足
var CommonError1 string = "Params is not enough"

// CommonError2 通用错误提示文本2：未定义显示内容
var CommonError2 string = "UI output is undefined"

// CommonError3 通用错误提示文本3：文件夹创建失败
var CommonError3 string = "Create directory failed, %v"

// CommonError4 通用错误提示文本4：文件创建失败
var CommonError4 string = "Create file failed, %v"

// CommonError5 通用错误提示文本：文件打开失败
var CommonError5 string = "Open file %v error, file is nil or %v"

// CommonError6 通用错误提示文本：解析数值错误
var CommonError6 string = "Parse alpha to integer error: %v"

// CommonNote1 通用提示文本1：未知选项
var CommonNote1 string = "Unknown command option"

// CommonWarn1 通用警告提示文本：待创建文件未指定后缀名
var CommonWarn1 string = "File to create does not have suffix"

// CommonWarn2 通用警告提示文本：指令选项的解析式不存在
var CommonWarn2 string = "Command %v option %v does not regist regexp"

// CommonWarn3 通用警告提示文本：原子表达式的解析式不存在
var CommonWarn3 string = "Atomic %v does not regist regexp"

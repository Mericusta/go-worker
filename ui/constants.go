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

// CommonError7 通用错误提示文本：文件或目录状态获取失败
var CommonError7 string = "Get file or directory %v stat error: %v"

// CommonError8 通用错误提示文本：%v 不是一个文件夹
var CommonError8 string = "%v is not a directory"

// CommonError9 通用错误提示文本：读取文件夹的内容时发生错误
var CommonError9 string = "Read directory %v occurs error: %v"

// CommonError10 通用错误提示文本：删除文件时发生错误
var CommonError10 string = "Remove file %v occurs error: %v"

// CommonError11 通用错误提示文本：文件或目录的绝对路径获取失败
var CommonError11 string = "Get file or directory %v absolute path error: %v"

// CommonError12 通用错误提示文本：读取目录下的文件失败
var CommonError12 string = "Read directory %v files error: %v"

// CommonError13 通用错误提示文本：读取文件内容错误
var CommonError13 string = "Read file %v content error: %v"

// CommonNote1 通用提示文本1：未知选项
var CommonNote1 string = "Unknown command option: %v"

// CommonNote2 通用显示文本1：分隔线
var CommonNote2 string = "----------------------------------------------------------------"

// CommonWarn1 通用警告提示文本：待创建文件未指定后缀名
var CommonWarn1 string = "File to create does not have suffix"

// CommonWarn2 通用警告提示文本：指令选项的解析式不存在
var CommonWarn2 string = "Command %v option %v does not regist regexp"

// CommonWarn3 通用警告提示文本：原子表达式的解析式不存在
var CommonWarn3 string = "Atomic Expression %v does not regist regexp"

// CommonWarn4 通用警告提示文本：文件不存在
var CommonWarn4 string = "File %v does not exist"

// CommonWarn5 通用警告提示文本：模板表达式的解析式不存在
var CommonWarn5 string = "Template Expression %v does not regist regexp"

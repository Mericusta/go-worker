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

// CommonError3 通用错误提示文本3：文件夹 %v 创建失败，错误：%v
var CommonError3 string = "Create directory %v failed, %v"

// CommonError4 通用错误提示文本4：文件 %v 创建失败，错误：%v
var CommonError4 string = "Create file %v failed, %v"

// CommonError5 通用错误提示文本：文件打开 %v 失败，错误：%v
var CommonError5 string = "Open file %v error, file is nil or %v"

// CommonError6 通用错误提示文本：解析数值错误，错误：%v
var CommonError6 string = "Parse alpha to integer error: %v"

// CommonError7 通用错误提示文本：文件或目录 %v 状态获取失败，错误：%v
var CommonError7 string = "Get file or directory %v stat error: %v"

// CommonError8 通用错误提示文本：%v 不是一个文件夹
var CommonError8 string = "%v is not a directory"

// CommonError9 通用错误提示文本：读取文件夹 %v 的内容时发生错误，错误：%v
var CommonError9 string = "Read directory %v occurs error: %v"

// CommonError10 通用错误提示文本：删除文件 %v 时发生错误，错误：%v
var CommonError10 string = "Remove file %v occurs error: %v"

// CommonError11 通用错误提示文本：文件或目录 %v 的绝对路径获取失败，错误：%v
var CommonError11 string = "Get file or directory %v absolute path error: %v"

// CommonError12 通用错误提示文本：读取目录 %v 下的文件失败，错误：%v
var CommonError12 string = "Read directory %v files error: %v"

// CommonError13 通用错误提示文本：读取文件 %v 内容失败，错误：%v
var CommonError13 string = "Read file %v content error: %v"

// CommonError14 通用错误提示文本：待创建文件未指定后缀名
var CommonError14 string = "File to create does not have suffix"

// CommonError15 通用错误提示文本：指令选项的解析式不存在
var CommonError15 string = "Command %v option %v does not regist regexp"

// CommonError16 通用错误提示文本：原子表达式 %v 的解析式不存在
var CommonError16 string = "Atomic Expression %v does not regist regexp"

// CommonError17 通用错误提示文本：文件 %v 不存在
var CommonError17 string = "File %v does not exist"

// CommonError18 通用错误提示文本：模板表达式的解析式不存在
var CommonError18 string = "Template Expression %v does not regist regexp"

// CommonError19 通用错误提示文本：变量 %v 的 %v 类型断言失败
var CommonError19 string = "Interface variable %v type %v assert failed"

// ----------------------------------------------------------------

// CommonNote1 通用提示文本1：未知选项：%v
var CommonNote1 string = "Unknown command option: %v"

// CommonNote2 通用显示文本2：分隔线
var CommonNote2 string = "----------------------------------------------------------------"

// CommonNote3 通用显示文本3：执行指令：%v
var CommonNote3 string = "Execute command: %v"

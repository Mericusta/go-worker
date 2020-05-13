package regexps

// Exit 解析指令 exit
var Exit Expression = `^exit$`

// Create 解析指令 create table [name] [column] [row]
var Create Expression = `^create table [\w]+ [\d]+ [\d]+$`

// Spider 解析指令 spider [project]
var Spider Expression = `^spider [\w]+$`

// Simulator 解析指令 simulator [project]
var Simulator Expression = `^simulator [\w]+$`

// Worker 解析指令 worker operation[ project]
var Worker Expression = `^worker [\w]+( [\w]+)*$`

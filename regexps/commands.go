package regexps

// AECommandExit 指令 exit 的原子表达式
var AECommandExit AtomicExpression = `^exit$`

// AECommandCustom 指令 custom 的原子表达式
var AECommandCustom AtomicExpression = `^custom\s+execute\s+\d+.*$`

// AECommandCSVChecker 指令 csv checker 的原子表达式
var AECommandCSVChecker AtomicExpression = `^csv\s+checker$`

package utility

import (
	"fmt"
	"strings"
)

// TestOutput 测试输出
func TestOutput(format string, content ...interface{}) {
	formatContent := fmt.Sprintf(format, content...)
	fmt.Printf("TEST: %v\n", formatContent)
}

// NoteOutput 提示输出
func NoteOutput(format string, content ...interface{}) {
	formatContent := fmt.Sprintf(format, content...)
	fmt.Printf("Note: %v\n", formatContent)
}

// ErrorOutput 错误输出
func ErrorOutput(format string, content ...interface{}) {
	formatContent := fmt.Sprintf(format, content...)
	fmt.Printf("Error: %v\n", formatContent)
}

// TraitStructName 从含有结构体类型的 GO 组合类型中萃取结构体的名称，如：*Name -> Name
func TraitStructName(structString string) string {
	structName := strings.TrimLeft(structString, "*")
	return structName
}

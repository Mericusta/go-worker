package utility

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-worker/global"
	"github.com/go-worker/regexps"
)

// FormatOutput 格式化输出
func FormatOutput(logMark, format string, content ...interface{}) {
	formatContent := fmt.Sprintf(format, content...)
	fmt.Printf("%v: %v\n", logMark, formatContent)
}

// TestOutput 测试输出
func TestOutput(format string, content ...interface{}) {
	FormatOutput(global.LogMarkTest, format, content...)
}

// IsExist 检查文件或文件夹是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// CreateDir 创建文件夹
func CreateDir(directoryPath string) error {
	err := os.Mkdir(directoryPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// CreateFile 创建文件
func CreateFile(filePath string) (*os.File, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// ExtractSetFromSlice 从切片中提取集合
func ExtractSetFromSlice(slice []string) map[string]int {
	set := make(map[string]int)
	for _, str := range slice {
		if _, hasStr := set[str]; hasStr {
			continue
		}
		set[str]++
	}
	return set
}

// SlicePop 取头部
func SlicePop(slice []string) (string, []string) {
	return slice[0], slice[1:]
}

// Convert2CamelStyle 将特定格式字符串转换为驼峰样式：xxx_yyy_zzz -> XxxYyyZzz,
func Convert2CamelStyle(otherStyleString string, capitalize bool) string {
	camelStyleString := ""
	for _, singleString := range strings.Split(otherStyleString, "_") {
		capitalizeSingleString := fmt.Sprintf("%v%v", strings.ToUpper(singleString[:1]), singleString[1:])
		camelStyleString = fmt.Sprintf("%v%v", camelStyleString, capitalizeSingleString)
	}
	if !capitalize {
		camelStyleString = fmt.Sprintf("%v%v", strings.ToLower(camelStyleString[:1]), camelStyleString[1:])
	}
	return camelStyleString
}

// ConvertTemplate2Format 将模板字符串转换为格式字符串
func ConvertTemplate2Format(templateString string) string {
	if templateKeywordRegexp, hasTemplateKeywordRegexp := regexps.ExpressionRegexpMap[global.TemplateKeywordExpression]; hasTemplateKeywordRegexp {
		return templateKeywordRegexp.ReplaceAllString(templateString, "%v")
	}
	return ""
}

// TraitStructName 从含有结构体类型的 GO 组合类型中萃取结构体的名称，如：*Name -> Name
func TraitStructName(structString string) string {
	structName := strings.TrimLeft(structString, "*")
	return structName
}

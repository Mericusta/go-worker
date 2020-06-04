package utility

import (
	"fmt"
	"os"
	"strings"
)

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

// TraitStructName 从含有结构体类型的组合类型中萃取结构体的名称，如：*Name -> Name，packageName.Name -> Name，*packageName.Name -> Name
func TraitStructName(structString string) string {
	structName := strings.TrimLeft(structString, "*")
	structNameList := strings.Split(structName, ".")
	if len(structNameList) == 1 {
		return structNameList[0]
	} else if len(structNameList) == 2 {
		return structNameList[1]
	}
	return structName
}

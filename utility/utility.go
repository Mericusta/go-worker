package utility

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
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

// TraverseDirectorySpecificFile 遍历文件夹获取所有绑定类型的文件
func TraverseDirectorySpecificFile(directory, syntax string) []string {
	traverseFileList := make([]string, 0)
	syntaxExt := fmt.Sprintf(".%v", syntax)
	filepath.Walk(directory, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if path.Ext(filePath) == syntaxExt {
			traverseFileList = append(traverseFileList, strings.Replace(filePath, "\\", "/", -1))
		}
		return nil
	})
	return traverseFileList
}

// TraverseDirectorySpecificFileWithFunction 遍历文件夹获取所有绑定类型的文件
func TraverseDirectorySpecificFileWithFunction(directory, syntax string, operate func(string, os.FileInfo)) {
	syntaxExt := fmt.Sprintf(".%v", syntax)
	filepath.Walk(directory, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if path.Ext(filePath) == syntaxExt {
			operate(filePath, info)
			// traverseFileList = append(traverseFileList, strings.Replace(filePath, "\\", "/", -1))
		}
		return nil
	})
}

// NTreeNode N 叉树节点
type NTreeNode struct {
	No       int
	Children []int
}

// NTreeHierarchicalMergeAlgorithmImproved N 叉树分层归并算法改进版
func NTreeHierarchicalMergeAlgorithmImproved(nTreeNodeChildrenMap map[int][]int) map[int]map[int]int {
	noNodeMap := make(map[int]*NTreeNode)
	levelNoMap := make(map[int]map[int]int)
	noMostLevelMap := make(map[int]int)
	for no := range nTreeNodeChildrenMap {
		noNodeMap[no] = &NTreeNode{
			No:       no,
			Children: make([]int, 0),
		}
	}

	rootNode := noNodeMap[0]
	if rootNode == nil {
		rootNode = noNodeMap[1]
	}

	if rootNode == nil {
		return levelNoMap
	}

	level := 0
	currentLevelNodeMap := make(map[int]int, 0)
	currentLevelNodeMap[rootNode.No] = rootNode.No
	for len(currentLevelNodeMap) != 0 {
		levelNoMap[level] = currentLevelNodeMap
		nextLevelNodeMap := make(map[int]int, 0)
		for _, currentLevelNode := range currentLevelNodeMap {
			noMostLevelMap[currentLevelNode] = level
			for _, subNode := range nTreeNodeChildrenMap[currentLevelNode] {
				nextLevelNodeMap[subNode] = subNode
			}
		}
		currentLevelNodeMap = nextLevelNodeMap
		level++
	}

	for level := 0; level != len(levelNoMap); level++ {
		for no := range levelNoMap[level] {
			noMostLevel := noMostLevelMap[no]
			if noMostLevel > level {
				delete(levelNoMap[level], no)
			}
		}
	}

	return levelNoMap
}

// NTreeHierarchicalMergeAlgorithm N 叉树分层归并算法
func NTreeHierarchicalMergeAlgorithm(nTreeNodeChildrenMap map[int][]int) map[int]map[int]int {
	noNodeMap := make(map[int]*NTreeNode)
	levelNoMap := make(map[int]map[int]int)
	for no := range nTreeNodeChildrenMap {
		noNodeMap[no] = &NTreeNode{
			No:       no,
			Children: make([]int, 0),
		}
	}
	rootNode := noNodeMap[0]
	if rootNode == nil {
		rootNode = noNodeMap[1]
	}

	if rootNode == nil {
		return levelNoMap
	}

	level := 0
	currentLevelNodeMap := make(map[int]int, 0)
	currentLevelNodeMap[rootNode.No] = rootNode.No
	for len(currentLevelNodeMap) != 0 {
		levelNoMap[level] = currentLevelNodeMap
		nextLevelNodeMap := make(map[int]int, 0)
		for _, currentLevelNode := range currentLevelNodeMap {
			for _, subNode := range nTreeNodeChildrenMap[currentLevelNode] {
				nextLevelNodeMap[subNode] = subNode
			}
		}
		currentLevelNodeMap = nextLevelNodeMap
		level++
	}

	for level := 0; level != len(levelNoMap); level++ {
		for no := range levelNoMap[level] {
			found := false
			for checkLevel := level + 1; checkLevel < len(levelNoMap); checkLevel++ {
				if _, hasNo := levelNoMap[checkLevel][no]; hasNo {
					found = true
					break
				}
			}
			if found {
				delete(levelNoMap[level], no)
			}
		}
	}

	return levelNoMap
}

// ReadFileLineOneByOne 逐行读取文件内容，执行函数返回 true 则继续读取，返回 false 则结束读取
func ReadFileLineOneByOne(filename string, f func(string) bool) error {
	file, openError := os.Open(filename)
	if openError != nil {
		return openError
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if !f(scanner.Text()) {
			break
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	return nil
}

// ReadContentLineOneByOne 逐行读取指定内容，执行函数返回 true 则继续读取，返回 false 则结束读取
func ReadContentLineOneByOne(reader io.Reader, f func(string) bool) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		if !f(scanner.Text()) {
			break
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	return nil
}

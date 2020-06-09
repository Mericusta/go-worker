## 2020.5.13

- 初始化
- 构建基础架构
    - FSM
    - Main Loop
    - CMD UI Package

## 2020.5.14

- NOTE: bind 指令
    - 绑定当前 worker 应用的项目路径或项目语法
    - 格式：bind [project|syntax] value
    - 指令可选参数
        - project
        - syntax
- NOTE: create 指令
    - 在绑定目录下创建包或文件
    - 格式：create [package|file] value [parent value]
    - 指令可选参数
        - package 创建包
        - file 创建文件
        - parent 指定父级包，不指定则默认在根目录

## 2020.5.20

- NOTE: convert 指令
    - 在绑定目录下通过指定格式生成绑定语法的结构体
    - 格式：convert [csv] value [parent value] [create|append] value
    - 指令可选参数
        - csv 通过 csv 生成
        - parent 指定 csv 的父目录
        - create 创建模式，创建一个绑定语法类型的文件，将生成内容添加到其中
        - append 追加模式，将生成内容追加到一个绑定语法类型的文件的末尾
- 层级结构
    - main 主模块
        - regexpscommands 正则表达式与指令模块
            - commands 指令模块
                - ui UI文本模块
                - config 配置模块
                    - utility 工具包模块
                        - regexps 正则表达式模块
                            - global 全局变量模块
                                - fsm 状态机模块

## 2020.5.24

- Note: analyze 指令
    - 通过绑定的语法分析绑定目录下指定的文件或语法内指定结构的内容，内容包括：函数定义，结构定义，非局部变量定义，其他定义，依赖或引用
    - 格式：analyze [file|directory|package] value [parent value] [output value]
    - 指令可选参数
        - file 分析文件
        - directory 文件夹
        - package go 语法指定结构
        - parent 指定相对绑定目录的父级路径
        - output 输出分析结果至
- RESTRUCT: regexps
    - 重构 regexps 模块
        - 通过细粒度的正则表达式组合的正则表达式解析输入的指令
        - 采用细粒度的正则表达式解析指令的参数
    - 为了更方便管理越来越多，越来越复杂的正则表达式，引入“模板文本”等概念，新建 template 模块
        - 模板文本格式：**$文本**
        - 模板表达式：含有**模板文本**的正则表达式，以 TExxx 命名
        - 模板匹配式：以纯**模板文本**构建的正则表达式，以 MTExxx 命名
        - 模板字符串：含有**模板文本**的非正则表达式，以 TSxxx 命名
        - 原子表达式：不含有**模板文本**的正则表达式，以 AExxx 命名
        - 转换函数：在正式使用模板表达式与模板字符串之前，必须将模板表达式与模板字符串中的模板文本转换为预定义的指定文本
    - 确定正则表达式的**设计原则**
        - 以尽量可重用的细粒度正则表达式和模板关键词进行排列组合构成文本模板
        - QUESTION: 格式上达到统一，概念上用模板关键词表达，如：
            - output value 与 append|create value 可统一为 $OPTION value
- RESTURCT: commands
    - 重构 commands 模块
        - 将 command 的枚举类型移到 global 中作为一个子模块，避免循环引用
        - 修改所有指令的 parseCommandParams 方法，采用正则表达式解析指令的参数
- 移动 README 内容至 doc/log.md，记录开发日志

## 2020.5.25

- 模板解析式：通过正则表达式编译生成的 *regexp.Regexp 对象

## 2020.5.26

- RESTURCT: regexps
  - 重构 regexp 模块
    - review 时发现随着原子匹配式和组合匹配式的增多，单条指令的解析会变得非常耗时，做了很多无用匹配操作
    - 模板文本格式，修改为**TK_文本**，因为 '$' 符号在正则表达式中存在歧义，某些函数无法处理
    - 指令模板表达式，以**TECommand指令**命名
    - 逻辑模板表达式，以**TE逻辑项描述**命名
    - 原子表达式中的特殊匹配式：模板关键词表达式`TK_[\w_]+`

## 2020.5.27

- 根据重构之后的 regexps 模块重构了 bind，create，convert 指令解析参数的逻辑

## 2020.5.28

- 根据重构之后的 regxps 模块重构了 analyze 指令解析参数的逻辑

## 2020.5.31

- 完善了 analyze 指令解析 Go 语法函数的问题

## 2020.6.1

- 在 analyze 指令输出分析结果时，发现一个问题：文本格式的继承问题，如分析函数，函数其下的替换文本的“格式”是需要定义在替换文本外还是替换文本内？如下：

```go
// AnalyzeGoFileFunctionDefinitionListTemplate 分析 go 文件，函数列表的文本模板
// “一个 \t+内容+换行”这种格式，在替换文本外
var AnalyzeGoFileFunctionDefinitionListTemplate = `
Function List:
\tRP_FUNCTION_DEFINITION
`

// AnalyzeGoFileFunctionDefinitionTemplate 分析 go 文件，函数的定义的文本模板
// “一个 \t+内容+换行”这种格式，在替换文本外，最终和外部 \t 组合起来，在 class 等内容就变成了“两个 \t+内容+换行”这种格式
var AnalyzeGoFileFunctionDefinitionTemplate = `
- RP_FUNC_NAME
\tRP_FUNCTION_CLASS
\tRP_FUNCTION_PARAM_LIST
\tRP_FUNCTION_RETURN_LIST
`
```

- RP_FUNCTION_DEFINITION 会替换成定义模板，模板中的内容也会被替换成对应的内容，结果示例：
```markdown
Function List:
    - Function Name（存在类定义的情况）
        - Class: Class Name
        - Params: ...
        - ...
    - Function Name（不存在类定义的情况）
        - Params: ...
        - ...
```

- 想要在不存在可替换值的情况下不显示对应的文本模板
    - 若在替换文本外定义替换文本的格式，则**需要侵入式的对替换文本处进行解析处理**生成最终的格式（多行替换文本的情况下，外层 \t 与内层 \t 组合的情况）
    - 若在替换文本内定义替换文本的格式，则会导致定义的模板可读性极差，如：
        ```go
        // AnalyzeGoFileResultTemplate 分析 go 文件，输出结果的文本模板
        var AnalyzeGoFileResultTemplate = `
        File Path:    RP_FILE_PATH
        Package Name: RP_PACKAGE_NAMERP_IMPORT_PACKAGE_LISTRP_FUNCTION_DEFINITION_LIST
        `
        // 若 RP_IMPORT_PACKAGE_LIST 或 RP_FUNCTION_DEFINITION_LIST 不存在替换内容，则可以直接使用空字符串而不用考虑空白行的问题

        // AnalyzeGoFileImportPackageListTemplate 分析 go 文件，依赖包列表的文本模板
        var AnalyzeGoFileImportPackageListTemplate = `
        \nImport:
        RP_IMPORT_PACKAGE
        `

        // AnalyzeGoFileImportPackageTemplate 分析 go 文件，依赖包的文本模板
        var AnalyzeGoFileImportPackageTemplate = `\t- RP_IMPORT_PACKAGE\n`
        ```

- **想要在不存在可替换值的情况下不显示对应的文本模板，并且不做侵入式处理，仅通过文本文替换实现**，但对于固有的格式，文件路径，包名，不需要刻意处理格式问题，因为这些是固有数据，若不存在则属于语法级别错误，如：
```go
// AnalyzeGoFileResultTemplate 分析 go 文件，输出结果的文本模板
var AnalyzeGoFileResultTemplate = `
File Path:    RP_FILE_PATH
Package Name: RP_PACKAGE_NAME
RP_IMPORT_PACKAGE_LIST
RP_FUNCTION_DEFINITION_LIST
`
```

- 由于在模板中，使用了 换行符，若无替换内容，则会出现空白行的情况，考虑一种做法：
    - 全部替换后，使用方法去除空白行（非侵入式的做法）
    - 将 \t 格式定义在替换文本外，但若内容存在多行文本，则将之移动到更底层的定义中（此举在层级越来越深之后会手写大量 \t）

- **想要在上述做法的基础上优化在层级嵌套越来越深的情况下多行替换文本需要大量定义 \t 的做法**
- 设计一种机制，可以让 \t 这种格式在多行文本中，针对每一行都应用
- 需要注意一点，代码编写或运行时一定是先从最底层的开始生成文本，然后再替换回父级，这样底层文本就无法得知父级文本指定的样式
- 另外，描述 \t 也需要一种模板，如定义三个 \t，以 `(\t,3)` 的形式定义而非 `\t\t\t`，同时要注意模板中多处定义格式模板的情况
- 参考之前空白行的处理方式：在替换时若父级存在格式定义，则针对替换文本的每一行都应用，**此为侵入式的做法**，需要解析模板的定义
- 最终决定的做法：
    - 全部替换后，使用方法去除空白行
    - 将 \t 格式定义在替换文本外，但若内容存在多行文本，则将之移动到更底层的定义中
    - 格式以`(格式字符,循环次数)`的形式定义，称为格式模板，添加原子表达式`\((?P<CHAR>(?:\\t)*),(?P<NUM>\d+)\)`用于匹配格式模板，该原子表达式限定了当前支持的格式字符

## 2020.6.2

- 把 utility 拆分成两部分，一部分只依赖官方标准包，一部分依赖项目内的其他包，但此举在将来**预计针对每一个层级都要一个 utility**，导致结构复杂
    - 想在 utility 包中添加一个通用的函数，功能是去除字符串中的空白行，用到了原子表达式，但是 regexps 中又使用了 utility 中的内容，如此造成了循环引用
    - 看了 [Go语言(golang)包设计哲学/原则与项目结构组织最佳实践](https://www.jianshu.com/p/92d5b9d96343) 之后有启发。现将仅依赖官方标准库的内容放在 utility 中，将其他工具函数放置在其依赖的核心包中，以文件的形式存在
    - 但是将 utiltiy 全部拆散开又不好做测试以及管理，暂时还是用层级 utility 的做法
- TODO: 观察到由于 ui 包的引用，导致整个包结构显得比较混乱，考虑将整个 ui 模块拆分成一个协程实现
- TODO: 使用 utility 的输出调试信息也可以考虑移除
- 实现了 analyze 指令分析 go 文件的基本功能，还需完善的细节
    - 函数名排序
    - 函数所属类标识中的指针类型萃取
    - 生成结果所用的 UI 模板的外部支持（外部 DSL）
- 修正：函数所属类标识中的指针类型萃取
- 修正：函数名排序，采用文件中定义的顺序

## 2020.6.3

- 添加函数体调用其他函数分析的功能
- Note: go 文件包引用时的别名处理
- Note: go 函数体中调用其他函数时是通过变量还是通过包引用的形式的判断

## 2020.6.4

- Note: go 文件的 package 需要知晓其全路径，以便于项目整体的分析
- Note: analyze 指令，去掉 package 选项
- Note: analyze 目录级分析，递归分析指定目录及其所有子目录下的所有绑定类型的文件
- TODO: analyze 项目级分析
- Note: go 函数分析，参数表，返回值类型萃取
- Note: go 函数分析，参数表，返回值，不萃取类型，完全保留包和指针
- Note: 紧急开发指令 remove，删除指定路径下的指令类型或名称的文件，用于解决 analyze 根目录生成过多文件分析文件的问题

## 2020.6.5

- Note: remove 添加 ignore 选项，添加排除文件夹
- Note: 修复 go 函数分析，返回值类型中含有 . 字符时匹配错误
- TODO: go 包级有向图
- Note: remove 的 ignore 选项排除的文件夹，不排除其子文件夹
- Note: 递归遍历改用标准库 filepath 实现，不再依赖 regexps
- Note: 所有拼装路径的地方全都改用标准库 filepath 实现

## 2020.6.6

- Note：go 包级有向图，层级归并算法
- Note: 修改包的输出样式
- Note: 修复包分析结果的错误，分析结果中出现循环引用导致层级归并算法出现错误：由于解析了注释 import 的代码引起的错误

## 2020.6.9

- Note: go 包级有向图，层级归并算法
    - 问题描述：已知某 N 叉树，以及该树中若干父子节点关系，所有父子节点之间可能存在直系父子或跨越若干深度之间的孙子关系，现要求用已知数据生成**合理层级关系的最简 N 叉树**，要求包含所有节点，示例
        ```txt
        - 0: {2, 9, 11, 3}
        - 1: {9, 6, 7, 3, 10, 4, 8, 2}
        - 2: {}
        - 3: {4, 9}
        - 4: {9}
        - 5: {9}
        - 6: {9, 10}
        - 7: {5, 3, 9}
        - 8: {9, 7, 3}
        - 9: {2}
        - 10: {}
        - 11: {1, 7, 3}
        ```
    - 实现思路1：用现有数据生成一个完整 N 叉树结构，再通过比较剪枝的方式去除多余的节点
        - 问题：由于该思路，需要针对非根节点的每一个节点与其他非根节点做对比，对于深度为 L 的完全 N 叉树，其时间复杂度为 (N!+(N-1)!+(N-2)!+...+2!)^2，考虑到实际中，L 的数量级为 10^2，N 的数量级为 10^3，有较大时间代价
        - 结果示例
            ```txt
            level 0: 0
            level 1: 11
            level 2: 1
            level 3: 6 8
            level 4: 10 7
            level 5: 3 5
            level 6: 4
            level 7: 9
            level 8: 2
            ```
            [](20200609180201.png)
        - 实现代码
            ```go
            type Node struct {
                No      int
                SubNode []*Node
            }

            type NodeNo struct {
                No            int
                SubNodeNoList []int
            }
            func main() {
                noListMap := map[int][]int{
                    0:  []int{2, 9, 11, 3},
                    2:  []int{},
                    9:  []int{2},
                    11: []int{1, 7, 3},
                    3:  []int{4, 9},
                    4:  []int{9},
                    5:  []int{9},
                    7:  []int{5, 3, 9},
                    8:  []int{9, 7, 3},
                    10: []int{},
                    12: []int{6, 10, 4, 8, 9, 7, 3},
                    1:  []int{9, 6, 7, 3, 10, 4, 8, 2},
                    6:  []int{9, 10},
                }

                noNodeMap := make(map[int]*NodeNo)
                for no := range noListMap {
                    noNodeMap[no] = &NodeNo{
                        No:            no,
                        SubNodeNoList: make([]int, 0),
                    }
                }
                rootNode := noNodeMap[0]

                levelNoMap := make(map[int]map[int]int)

                level := 0
                currentLevelNodeMap := make(map[int]int, 0)
                currentLevelNodeMap[rootNode.No] = rootNode.No
                for len(currentLevelNodeMap) != 0 {
                    levelNoMap[level] = currentLevelNodeMap
                    nextLevelNodeList := make(map[int]int, 0)
                    for _, currentLevelNode := range currentLevelNodeMap {
                        for _, subNode := range noListMap[currentLevelNode] {
                            nextLevelNodeList[subNode] = subNode
                        }
                    }
                    currentLevelNodeMap = nextLevelNodeList
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

                for level := 0; level != len(levelNoMap); level++ {
                    fmt.Println(strings.Repeat("\t", level), "- level", level, "node list:", levelNoMap[level])
                }
            }
            ```
    - 实现思路2：边生成 N 叉树边去除
        - 问题：该思路实现并不能保证结果正确，并且目前的实现在最终叶子节点的结果处理上存在问题，但由于所有节点只迭代一遍并且在迭代过程中改变了节点的父子关系减少了判定依据，所以其复杂度为给出条件中的节点总和，如该条件：父子节点总共 40 个
        - 结果示例
            ```txt
            level 0: 0
            level 1: 11
            level 2: 1
            level 3: 6 8
            level 4: 10 7
            level 5: 3 5
            level 6: 9 4 9
            level 7: 2
            ```
            [](20200609180243.png)
        - 实现代码
            ```go
            type Node struct {
                No      int
                SubNode []*Node
            }

            type NodeNo struct {
                No            int
                SubNodeNoList []int
            }
            // 由于变换过程中改变了原节点的父子关系，所以叶子节点存在无法消去的同层次相同节点
            // - level 0 node 0 sub node [11]
            // - level 1 node 11 sub node [1]
            // - level 2 node 1 sub node [6 8]
            // - level 3 node 6 sub node [10]
            // - level 3 node 8 sub node [7]
            // - level 4 node 10 sub node []
            // - level 4 node 7 sub node [5 3]
            // - level 5 node 5 sub node [9]
            // - level 5 node 3 sub node [4 9]
            // - level 6 node 9 sub node [2] -> 9: [2] 相同节点
            // - level 6 node 4 sub node [] -> 4: [] 原本 4 是 4: [9]，可以消除同层次的 9: [2]，但是变换过程中改变了了该节点的父子关系
            // - level 6 node 9 sub node [2] -> 9: [2] 相同节点
            // - level 7 node 2 sub node [] -> 2: [] 相同节点
            // - level 7 node 2 sub node [] -> 2: [] 相同节点
            // 改变前：
            // - 0: {2, 9, 11, 3}
            // - 1: {9, 6, 7, 3, 10, 4, 8, 2}
            // - 2: {}
            // - 3: {4, 9}
            // - 4: {9}
            // - 5: {9}
            // - 6: {9, 10}
            // - 7: {5, 3, 9}
            // - 8: {9, 7, 3}
            // - 9: {2}
            // - 10: {}
            // - 11: {1, 7, 3}
            // - 12: {6, 10, 4, 8, 9, 7, 3}
            // 改变后：
            // - 0: {11}
            // - 1: {6, 8}
            // - 2: {}
            // - 3: {4, 9}
            // - 4: {}
            // - 5: {9}
            // - 6: {10}
            // - 7: {5, 3}
            // - 8: {7}
            // - 9: {2}
            // - 10: {}
            // - 11: {1}
            // - 12: {6, 10, 4, 8, 9, 7, 3}
            func endVersion() {
                noListMap := map[int][]int{
                    0:  []int{2, 9, 11, 3},
                    2:  []int{},
                    9:  []int{2},
                    11: []int{1, 7, 3},
                    3:  []int{4, 9},
                    4:  []int{9},
                    5:  []int{9},
                    7:  []int{5, 3, 9},
                    8:  []int{9, 7, 3},
                    10: []int{},
                    12: []int{6, 10, 4, 8, 9, 7, 3},
                    1:  []int{9, 6, 7, 3, 10, 4, 8, 2},
                    6:  []int{9, 10},
                }

                noNodeMap := make(map[int]*NodeNo)
                for no := range noListMap {
                    noNodeMap[no] = &NodeNo{
                        No:            no,
                        SubNodeNoList: make([]int, 0),
                    }
                }
                rootNode := noNodeMap[0]

                forLevel := 0
                currentLevel := 0
                nextLevel := currentLevel + 1
                lastLevel := currentLevel - 1
                levelNodeListMap := make(map[int][]int, 0)
                levelNodeListMap[currentLevel] = []int{
                    rootNode.No,
                }
                for forLevel != 9 {
                    currentLevel = forLevel
                    nextLevel = currentLevel + 1
                    lastLevel = currentLevel - 1
                    // 当前层级
                    if _, hasCurrentLevel := levelNodeListMap[currentLevel]; !hasCurrentLevel {
                        levelNodeListMap[currentLevel] = make([]int, 0)
                    }
                    lastLevelNodeList := levelNodeListMap[lastLevel]
                    fmt.Println("last level:", lastLevel, "node list:", lastLevelNodeList)
                    currentLevelNodeList := levelNodeListMap[currentLevel]
                    fmt.Println("current level:", currentLevel, currentLevelNodeList)

                    // 迭代终止
                    if len(currentLevelNodeList) == 0 {
                        break
                    }

                    // 开始计算
                    fmt.Println("---------------- start calculate ----------------")
                    removedCurrentLevelNodeList := make([]int, 0)
                    for _, currentLevelNode := range currentLevelNodeList {
                        fmt.Println("current level", currentLevel, "node", currentLevelNode, "sub node", noListMap[currentLevelNode])
                        currentNodeRemoved := false
                        for _, removedCurrentLevelNode := range removedCurrentLevelNodeList {
                            if removedCurrentLevelNode == currentLevelNode {
                                currentNodeRemoved = true
                                removedCurrentLevelNodeSubNode := noListMap[removedCurrentLevelNode]
                                fmt.Println("removedCurrentLevelNodeSubNode", removedCurrentLevelNodeSubNode)
                                break
                            }
                        }
                        if !currentNodeRemoved {
                            for _, subNode := range noListMap[currentLevelNode] {
                                fmt.Println("-- check if subNode", subNode, "need remove from parent node")
                                // fromCurrentNode := false
                                for checkLevel := 0; checkLevel != currentLevel; checkLevel++ {
                                    fmt.Println("---- check level", checkLevel, "node list", levelNodeListMap[checkLevel])
                                    afterRemoveNextLevelNode := make([]int, 0)
                                    for _, checkNode := range levelNodeListMap[checkLevel] {
                                        fmt.Println("------ check level", checkLevel, "node", checkNode, "sub node", noListMap[checkNode])
                                        afterRemoveCheckNodeSubNode := make([]int, 0)
                                        for _, checkNodeSubNode := range noListMap[checkNode] {
                                            if checkNodeSubNode == subNode {
                                                fmt.Println("-------- node", subNode, "need remove from leve", checkLevel, "node", checkNode, "sub node", noListMap[checkNode])
                                                fmt.Println("-------- node", subNode, "need remove from check next leve", checkLevel+1, "node list", levelNodeListMap[checkLevel+1])
                                                if checkLevel+1 == currentLevel {
                                                    // 记录当前层级由于被高层级节点的子节点移除的节点，减少计算过程
                                                    fmt.Println("-------- node", subNode, "in current level", currentLevel, "removed")
                                                    // fromCurrentNode = true
                                                    removedCurrentLevelNodeList = append(removedCurrentLevelNodeList, subNode)
                                                }
                                            } else {
                                                afterRemoveCheckNodeSubNode = append(afterRemoveCheckNodeSubNode, checkNodeSubNode)
                                                afterRemoveNextLevelNode = append(afterRemoveNextLevelNode, checkNodeSubNode)
                                            }
                                        }
                                        noListMap[checkNode] = afterRemoveCheckNodeSubNode
                                        fmt.Println("------ after remove subNode", subNode, "check level", checkLevel, "check node", checkNode, "sub node is", noListMap[checkNode])
                                    }
                                    levelNodeListMap[checkLevel+1] = afterRemoveNextLevelNode
                                    fmt.Println("---- after remove subNode", subNode, "next check level", checkLevel+1, "node list", levelNodeListMap[checkLevel+1])
                                }
                                // TODO: 添加到层级节点列表中，添加毫无意义，下一层及总是要重新计算的
                                // fmt.Println("-- add subNode", subNode, "to next level", nextLevel, "node list", levelNodeListMap[nextLevel])
                                // levelNodeListMap[nextLevel] = append(levelNodeListMap[nextLevel], subNode)

                                // 处理由于移除当前层级的节点而影响的子节点
                                // if fromCurrentNode {
                                // fmt.Println("-- subNode", subNode, "children", noListMap[subNode], "need remove from next level", nextLevel, "node list", levelNodeListMap[nextLevel])
                                // for _, subNodeSubNode := range noListMap[subNode] {
                                // 	removeIndex := -1
                                // 	for index, no := range levelNodeListMap[nextLevel] {
                                // 		if no == subNodeSubNode {
                                // 			removeIndex = index
                                // 			break
                                // 		}
                                // 	}
                                // 	if removeIndex != -1 {
                                // 		if removeIndex < len(levelNodeListMap[nextLevel]) {
                                // 			levelNodeListMap[nextLevel] = append(levelNodeListMap[nextLevel][0:removeIndex], levelNodeListMap[nextLevel][removeIndex+1:]...)
                                // 		}
                                // 	}
                                // }
                                // fmt.Println("-- after remove subNode", subNode, "children", noListMap[subNode], "removed from next level", nextLevel, "node list", levelNodeListMap[nextLevel])
                                // }
                            }
                        } else {
                            fmt.Println("current level", currentLevel, "node", currentLevelNode, "has been removed, continue calculate")
                        }
                        fmt.Println("after calculate next level", nextLevel, "node list", levelNodeListMap[nextLevel])
                        fmt.Println("---------------- node", currentLevelNode, "calculate done ----------------")
                        wideFirstTraverseNodeNo(rootNode, noListMap)
                    }
                    nextLevelNodeListFromCurrentLevelNodeList := make([]int, 0)
                    for _, currentLevelNode := range levelNodeListMap[currentLevel] {
                        for _, currentLevelNodeSubNode := range noListMap[currentLevelNode] {
                            nextLevelNodeListFromCurrentLevelNodeList = append(nextLevelNodeListFromCurrentLevelNodeList, currentLevelNodeSubNode)
                        }
                    }
                    fmt.Println("after calculate current level", currentLevel, "next level", nextLevel, "node list", levelNodeListMap[nextLevel], "should regenerate by current level", currentLevel, "node list", levelNodeListMap[currentLevel], "sub node list", nextLevelNodeListFromCurrentLevelNodeList)
                    levelNodeListMap[nextLevel] = nextLevelNodeListFromCurrentLevelNodeList
                    fmt.Println("---------------- stop calculate ----------------")

                    // 层级递增
                    forLevel++

                    fmt.Println("--------------------------------", forLevel, "--------------------------------")
                }

                fmt.Println("no list map")
                for no, list := range noListMap {
                    fmt.Println("no", no, "sub node", list)
                }
            }

            func wideFirstTraverseNodeNo(rootNode *NodeNo, noListMap map[int][]int) {
                level := 0
                currentLevelNodeList := make([]int, 0)
                currentLevelNodeList = append(currentLevelNodeList, rootNode.No)
                for len(currentLevelNodeList) != 0 {
                    nextLevelNodeList := make([]int, 0)
                    for _, currentLevelNode := range currentLevelNodeList {
                        fmt.Println(strings.Repeat("\t", level), "- level", level, "node", currentLevelNode, "sub node", noListMap[currentLevelNode])
                        nextLevelNodeList = append(nextLevelNodeList, noListMap[currentLevelNode]...)
                    }
                    currentLevelNodeList = nextLevelNodeList
                    level++
                }
            }
            ```
- Note: 
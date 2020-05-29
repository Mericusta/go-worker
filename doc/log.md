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

- TODO: analyze 指令
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

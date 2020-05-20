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


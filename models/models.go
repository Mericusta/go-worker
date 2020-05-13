package models

import (
	"github.com/go-worker/global"
)

// Table 表格
type Table struct {
	No    int64
	Name  string
	Head  []string
	Value []map[string]string
}

// Plan 计划
type Plan struct {
	No        int64
	MainType  global.MainType
	SubType   global.SubType
	Target    Object
	Operation Operation
	Param     string
}

// Object 对象
type Object struct {
	No        int
	MainType  global.MainType
	SubType   global.SubType
	Name      string
	Attribute interface{}
}

// Operation 操作
type Operation struct {
	No      global.OperationNo
	Content string
}

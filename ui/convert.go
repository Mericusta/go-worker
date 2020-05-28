package ui

// GoStructTemplate 生成 go 结构体的文本模板
var GoStructTemplate = `
type RP_STRUCT_NAME struct {
	RP_STRUCT_MEMBER
}
`

// GoStructMemberTemplate 生成 go 结构体成员的文本模板
var GoStructMemberTemplate = "RP_MEMBER_NAME RP_MEMBER_TYPE RP_MEMBER_COMMENT"

// GoMemberCommentByCSV 由 csv 生成 go 结构体成员的标识文本模板
var GoMemberCommentByCSV = "`csv:\"RP_CSV_HEAD\"`"

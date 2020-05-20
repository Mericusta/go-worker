package ui

// GoStructTemplate 生成 go 结构体的模板
var GoStructTemplate = `
type $STRUCT_NAME struct {
	$STRUCT_MEMBER
}
`

// GoStructMemberTemplate 生成 go 结构体成员的模板
var GoStructMemberTemplate = "$MEMBER_NAME $MEMBER_TYPE $MEMBER_COMMENT"

// GoMemberCommentByCSV 由 csv 生成 go 结构体成员的标识模板
var GoMemberCommentByCSV = "`csv:\"$CSV_HEAD\"`"

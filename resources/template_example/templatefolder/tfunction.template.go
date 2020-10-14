package templatefolder

import "github.com/go-worker/resources/template"

// template function deduction

// template <typename T1, typename T2, typename T3>
// func templateOperatorPlus(t1 T1, t2 T2) T3 {
func TemplateOperatorPlus(t1 template.TypeName, t2 template.TypeName) template.TypeName {
	// Go 中没有运算符重载，所以使用该 templateAdd 的的 template.TypeName 只能被推导为内建支持 + 运算符的类型
	var t3 template.TypeName
	t3 = t1 + t2
	return t3
}

func TemplateSelfIncrease(t1 template.TypeName) template.TypeName {
	return func(template.TypeName) template.TypeName {
		return t1 + 1
	}
}

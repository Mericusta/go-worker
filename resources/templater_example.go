package resources

import (
	"github.com/go-worker/resources/template"
)

func templateOperatorPlus(t1 template.TypeName, t2 template.TypeName) template.TypeName {
	// Go 中没有运算符重载，所以使用该 templateAdd 的的 template.TypeName 只能被推导为内建支持 + 运算符的类型
	return t1 + t2
}

func templateOperator

type exampleStruct struct {
	v int
}

// Plus is a function
func (es exampleStruct) Plus(oes exampleStruct) {

}

func main() {
	templateAdd(1, 2)

	templateAdd(1, 2.0)

	templateAdd([]int{1}, []int{2})

	templateAdd(exampleStruct{v: 1}, exampleStruct{v: 2})

	templateAdd(&exampleStruct{v: 1}, &exampleStruct{v: 2})
}

package resources

import (
	"github.com/go-worker/resources/template"
)

type exampleStruct struct {
	v int // This is a struct member
}

// Plus is a function
func (es exampleStruct) Plus(oes exampleStruct) {

}

// template function deduction

func templateOperatorPlus(t1 template.TypeName, t2 template.TypeName) template.TypeName {
	// Go 中没有运算符重载，所以使用该 templateAdd 的的 template.TypeName 只能被推导为内建支持 + 运算符的类型
	var t3 template.TypeName
	t3 = t1 + t2
	return t3
}

// template struct dedection

type templateStruct struct {
	tV template.TypeName
	iV int
}

func (ts templateStruct) Result() template.TypeName {
	return ts.tV + ts.iV
}

func main() {
	// T -> func(int, int) int
	v1 := templateOperatorPlus(1, 2)

	// T -> func(int, float) float -> implicit type conversion
	v2 := templateOperatorPlus(1, 2.0)

	// T -> func([]int, []int) []int -> maybe STL Vector
	v3 := templateOperatorPlus([]int{1}, []int{2})

	// T -> func(exampleStruct, exampleStruct) exampleStruct
	v4 := templateOperatorPlus(exampleStruct{v: 1}, exampleStruct{v: 2})

	// T -> func(*exampleStruct, *exampleStruct) *exampleStruct
	v5 := templateOperatorPlus(&exampleStruct{v: 1}, &exampleStruct{v: 2})

	// T -> struct { int, int }
	v6 := templateStruct{tV: 1, iV: 2}

	// T -> struct { float64, int }
	v7 := templateStruct{tV: 1.0, iV: 2}

	// T -> struct { exampleStruct, int }
	v8 := templateStruct{tV: exampleStruct{}, iV: 2}

	// T -> func() exampleStruct
	v9 := v8.Result()

	// no call
	v10 := templateOperatorPlus
}

package templateexample

import "github.com/go-worker/resources/template"

type ExampleStruct struct {
	v int // This is a struct member
}

// Plus is a function
func (es ExampleStruct) Plus(oes ExampleStruct) {

}

// template function deduction

// template <typename T1, typename T2, typename T3>
// func templateOperatorPlus(t1 T1, t2 T2) T3 {
func TemplateOperatorPlus(t1 template.TypeName, t2 template.TypeName) (t3 template.TypeName, t4 template.TypeName) {
	// Go 中没有运算符重载，所以使用该 templateAdd 的的 template.TypeName 只能被推导为内建支持 + 运算符的类型
	t3 = t1 + t2
	return t3, t1 - t2
}

// // 不支持匿名函数推导
// func TemplateSelfIncrease(t1 template.TypeName) template.TypeName {
// 	return func(template.TypeName) template.TypeName {
// 		return t1 + 1
// 	}
// }

// template struct dedection

type TemplateStruct struct {
	tV template.TypeName
	iV int
}

func (ts TemplateStruct) Result() template.TypeName {
	return ts.tV + ts.iV
}

func templateCaller() {
	// T -> func(int, int) int
	v1 := TemplateOperatorPlus(1, 2)

	// T -> func(int, float64) float64 -> implicit type conversion
	v2 := TemplateOperatorPlus(1, 2.0)

	// T -> func([]int, []int) []int -> maybe STL Vector
	v3 := TemplateOperatorPlus([]int{1}, []int{2})

	// T -> func(ExampleStruct, ExampleStruct) ExampleStruct
	v4 := TemplateOperatorPlus(ExampleStruct{v: 1}, ExampleStruct{v: 2})

	// T -> func(*ExampleStruct, *ExampleStruct) *ExampleStruct
	v5 := TemplateOperatorPlus(&ExampleStruct{v: 1}, &ExampleStruct{v: 2})

	// T -> struct { int, int }
	v6 := TemplateStruct{tV: 1, iV: 2}

	// T -> struct { float64, int }
	v7 := TemplateStruct{tV: 1.0, iV: 2}

	// T -> struct { exampleStruct, int }
	v8 := TemplateStruct{tV: ExampleStruct{}, iV: 2}

	// T -> func() exampleStruct
	v9 := v8.Result()

	// no call
	v10 := TemplateOperatorPlus

	// T -> func(int, float32) float32 -> explicit specify type to float32
	v11 := TemplateOperatorPlus(float32(1), float32(2.0))

	// T -> func(complex128, complex128) -> complex128
	v12 := TemplateOperatorPlus(1+2i, 1.1+2.2i)
}

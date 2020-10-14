package template

// TypeName 类型关键词
type TypeName interface{}

// Go 中没有运算符重载，借助接口来让使用者自己实现重载运算

// 内建支持 + - * / % ++ -- 的类型：数字类型
// uint8, uint16, uint32, uint64
// int8, int16, int32, int64
// float32, float64
// complex64, complex128
// byte, rune, uint, int, uintptr

// OperatorPlus 运算符 + 的接口
type OperatorPlus interface{ Plus(v TypeName) TypeName }

// T -> func(int, int) int
// templateOperatorPlus(1, 2)

// T -> func(int, float) float -> implicit type conversion
// templateOperatorPlus(1, 2.0)

// T -> func([]int, []int) []int -> maybe STL Vector
// templateOperatorPlus([]int{1}, []int{2})

// T -> func(exampleStruct, exampleStruct) exampleStruct
// templateOperatorPlus(exampleStruct{v: 1}, exampleStruct{v: 2})

// T -> func(*exampleStruct, *exampleStruct) *exampleStruct
// templateOperatorPlus(&exampleStruct{v: 1}, &exampleStruct{v: 2})

// T -> func(int, float32) float64 -> explicit specify type to float64
// ambiguous type deduction:
// T -> func(int, float32) float64
// T -> func(int, float64) float64
// ...
// T -> func(float64, float64) float64
// v11 := templateOperatorPlus(1, 2.0)
// to explicit type deduction, should using:
// T -> func(int, float32) float32
// v11 := templateOperatorPlus(1, float32(2.0))

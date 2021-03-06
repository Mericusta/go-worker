// This is template example file

package templateexample

import (
	f "fmt"

	_ "math"

	"github.com/go-worker/resources/template"
)

var e0 int = 1
var e1 EmptyInterface
var e2 template.TypeName

const c1 int = 1
const ti template.TemplateInt = 1
const (
	c2 int = iota
	c3
	c4
)

const (
	c5 int = iota
	c6
	t2 template.TemplateInt = 1
	t3
)

type ri interface{}
type r2ifunc func() (
	interface{},
	interface{},
)

// common interface

type EmptyInterface interface{}
type OneLineInterface interface{ OneLineFunction() template.TypeName }
type AnotherOneLineInterface interface{ AnotherOneLineFunction() interface{} }
type ExampleInterface interface {
	Example()
	FuncParamExample(
		v template.TypeName,
		f func(t template.TypeName, sf func() interface{}) func() interface{},
	) func() interface{}
	FuncReturnExample() (int, func() int)
}

// common struct

type OneLineStruct struct{ v int }
type ExampleStruct struct {
	OneLineStruct
	v interface{} // This is a struct member
	f func(int, int) int
}

// common struct define template function
func (es ExampleStruct) Set(t1 template.TypeName) { es.v = t1 } // return T

// common struct define template function
func (es ExampleStruct) Get() template.TypeName {
	return es.v
}

// ----------------------------------------------------------------

// template interface
type TemplateInterface interface {
	TExample(t1 template.TypeName, t2 template.TypeName)
}

// template struct
type TemplateStruct struct {
	tV template.TypeName
	iV int
}

// template struct define template function
func (ts TemplateStruct) Set(v template.TypeName) {
	ts.tV = v
}

// template struct define template function
func (ts TemplateStruct) Get() template.TypeName {
	return ts.tV
}

// ----------------------------------------------------------------

// common struct implement common interface
func (es ExampleStruct) Example() {
	f.Println("This is Example from ExampleStruct")
}

// template struct implement common interface
func (ts TemplateStruct) Example() {
	f.Println("This is Example from TemplateStruct")
}

// ----------------------------------------------------------------

// common struct implement template interface
func (es ExampleStruct) TExample(t1 template.TypeName, t2 template.TypeName) {
	f.Println("This is TExample from ExampleStruct")
}

// template struct implement template interface
func (ts TemplateStruct) TExample(t1 template.TypeName, t2 template.TypeName) (template.TypeName, template.TypeName) {
	f.Println("This is TExample from TemplateStruct")
}

// ----------------------------------------------------------------

// define template function
func TemplateOperatorADD(t1 template.TypeName, t2 template.TypeName) (t3 template.TypeName, t4 template.TypeName) {
	// Go 中没有运算符重载，所以使用该 templateAdd 的的 template.TypeName 只能被推导为内建支持 + 运算符的类型
	return t1, t2
}

func TemplateOperatorSUB(t1 template.TypeName, t2 template.TypeName) template.TypeName {
	return t2
}

func TemplateOperatorMUL(t1 template.TypeName, t2 interface{}, t3 struct{ t int }) interface{} {
	return t3
}

func TemplateOperatorDIV(t1 template.TypeName, t2 interface{}, t3 struct{ t int }) interface{} {
	return nil
}

func TemplateOperatorMOD(t1 template.TypeName, t2 interface{}, t3 struct{ t int }) {}

func TemplateOperatorREM(t1 int, t2 interface{}, t3 struct{ t int }) struct{ v int } {
	return struct{ v int }{v: 1}
}

func TemplateOperatorCustom(t1 template.TypeName, t2 template.TypeName, f func(template.TypeName, template.TypeName) template.TypeName) template.TypeName {
	return f(t1, t2)
}

// // 不支持匿名函数推导
// func TemplateSelfIncrease(t1 template.TypeName) template.TypeName {
// 	return func(template.TypeName) template.TypeName {
// 		return t1 + 1
// 	}
// }

func templateCaller() {
	// T -> func(int, int) int
	v1, _ := TemplateOperatorADD(1, 2)

	// T -> func(int, float64) float64 -> implicit type conversion
	_, v2 := TemplateOperatorADD(1, 2.0)

	// T -> func(string, string) string -> maybe STL Vector
	TemplateOperatorADD("1", "2")

	// T -> func([]int, []int) []int -> maybe STL Vector
	TemplateOperatorADD([]int{1}, []int{2})

	// T -> func(ExampleStruct, ExampleStruct) ExampleStruct
	es, _ := TemplateOperatorADD(ExampleStruct{v: 1}, ExampleStruct{v: v1})

	// T -> func(String)
	es.Set(v2)

	// T -> func()interface{}
	es.Get()

	// call interface Example
	es.Example()

	// call template interface TExample
	es.TExample()

	// T -> func(*ExampleStruct, *ExampleStruct) *ExampleStruct
	TemplateOperatorADD(&ExampleStruct{v: 1}, &ExampleStruct{v: 2})

	var ts TemplateStruct

	// T -> struct { int, int }
	ts = TemplateStruct{tV: 1, iV: 2}

	// T -> struct { float64, int }
	ts = TemplateStruct{tV: 1.0, iV: 2}

	// T -> struct { exampleStruct, int }
	ts = TemplateStruct{tV: ExampleStruct{}, iV: 2}

	// T -> func(ExampleStruct)
	ts.Set(es)

	// T -> func() ExampleStruct
	v4 := ts.Get()

	// call interface Example
	ts.Example()

	// call template interface TExample
	ts.TExample(v4, v4)

	// no call
	v5 := TemplateOperatorADD

	// T -> func(int, float32) float32 -> explicit specify type to float32
	v5(float32(1), float32(2.0))

	// T -> func(complex128, complex128) -> complex128
	v5(1+2i, 1.1+2.2i)

	f.Println("This is Println call from alias f -> fmt")

	v6 := v5

	v5(
		ts.Example(),
		ts.TExample(
			v4,
			v4,
		),
	)

	v7 := ts.Example().Get()
	v8 := ts.Ts.Get()
	v9 := ts.Ts.Get(ts.Example('(').Get(')')).Set(TExample(
		ts.Example().es.Get(),
		ts.Ts.Get(),
	))
}

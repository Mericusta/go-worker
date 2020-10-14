package templateexample

import (
	"github.com/go-worker/resources/template_example/nontemplatefolder"
	"github.com/go-worker/resources/template_example/templatefolder"
)

func templateCaller() {
	// T -> func(int, int) int
	v1 := templatefolder.TemplateOperatorPlus(1, 2)

	// T -> func(int, float64) float64 -> implicit type conversion
	v2 := templatefolder.TemplateOperatorPlus(1, 2.0)

	// T -> func([]int, []int) []int -> maybe STL Vector
	v3 := templatefolder.TemplateOperatorPlus([]int{1}, []int{2})

	// T -> func(ExampleStruct, ExampleStruct) ExampleStruct
	v4 := templatefolder.TemplateOperatorPlus(nontemplatefolder.ExampleStruct{v: 1}, nontemplatefolder.ExampleStruct{v: 2})

	// T -> func(*ExampleStruct, *ExampleStruct) *ExampleStruct
	v5 := templatefolder.TemplateOperatorPlus(&nontemplatefolder.ExampleStruct{v: 1}, &nontemplatefolder.ExampleStruct{v: 2})

	// T -> struct { int, int }
	v6 := templatefolder.TemplateStruct{tV: 1, iV: 2}

	// T -> struct { float64, int }
	v7 := templatefolder.TemplateStruct{tV: 1.0, iV: 2}

	// T -> struct { exampleStruct, int }
	v8 := templatefolder.TemplateStruct{tV: nontemplatefolder.ExampleStruct{}, iV: 2}

	// T -> func() exampleStruct
	v9 := v8.Result()

	// no call
	v10 := templatefolder.TemplateOperatorPlus

	// T -> func(int, float32) float32 -> explicit specify type to float32
	v11 := templatefolder.TemplateOperatorPlus(float32(1), float32(2.0))

	// T -> func(complex128, complex128) -> complex128
	v12 := templatefolder.TemplateOperatorPlus(1+2i, 1.1+2.2i)
}

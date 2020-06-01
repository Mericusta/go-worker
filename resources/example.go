package main

type ExampleStruct struct {
	value int
}

func (example ExampleStruct) ExampleFunc7(ex *ExampleStruct) *ExampleStruct {
	return nil
}

func (example ExampleStruct) ExampleFunc9(
	str string,
	value int,
	example2 ExampleStruct,
) (int, int) {
	return 0, 0
}

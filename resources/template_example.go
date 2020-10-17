
func (ts TemplateStruct) ResultType0() template.TypeName {
	return ts.tV + ts.iV
}

func TemplateOperatorPlusType4(t2 template.TypeName, t1 template.TypeName) (t3 template.TypeName, t4 template.TypeName) {
	t3 = t1 + t2
	return t3, t1 - t2
}

func TemplateOperatorPlusType5(t1 float32, t2 float32) (t3 template.TypeName, t4 template.TypeName) {
	t3 = t1 + t2
	return t3, t1 - t2
}

func TemplateOperatorPlusType6(t2 complex128, t1 complex128) (t3 template.TypeName, t4 template.TypeName) {
	t3 = t1 + t2
	return t3, t1 - t2
}

func TemplateOperatorPlusType0(t1 int, t2 int) (t3 template.TypeName, t4 template.TypeName) {
	t3 = t1 + t2
	return t3, t1 - t2
}

func TemplateOperatorPlusType1(t1 int, t2 float64) (t3 template.TypeName, t4 template.TypeName) {
	t3 = t1 + t2
	return t3, t1 - t2
}

func TemplateOperatorPlusType2(t1 template.TypeName, t2 template.TypeName) (t3 template.TypeName, t4 template.TypeName) {
	t3 = t1 + t2
	return t3, t1 - t2
}

func TemplateOperatorPlusType3(t1 template.TypeName, t2 template.TypeName) (t3 template.TypeName, t4 template.TypeName) {
	t3 = t1 + t2
	return t3, t1 - t2
}

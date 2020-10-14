package templatefolder

import "github.com/go-worker/resources/template"

// template struct dedection

type TemplateStruct struct {
	tV template.TypeName
	iV int
}

func (ts TemplateStruct) Result() template.TypeName {
	return ts.tV + ts.iV
}

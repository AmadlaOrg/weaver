package template

import (
	"fmt"
	"github.com/AmadlaOrg/weaver/hery"
	"os"
	"text/template"
)

type ITemplate interface {
}

type STemplate struct {
	Hery hery.IHery
}

// ListTemplates
func (s *STemplate) ListTemplates() {

}

// Weave
func (s *STemplate) Weave(templatePaths string, data any) string {

	funcMap := template.FuncMap{
		"hery": s.Hery.HeryFunc,
	}

	// Doc: https://pkg.go.dev/text/template#FuncMap

	tmpl, err := template.New("config").Funcs(funcMap).ParseFiles(templatePaths)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return ""
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
	}
}

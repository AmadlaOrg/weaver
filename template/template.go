package template

import (
	"fmt"
	"github.com/AmadlaOrg/weaver/hery"
	"os"
	"text/template"
)

// Service defines the template operations interface.
type Service interface {
	ListTemplates()
	Weave(templatePaths string, data any)
}

type service struct {
	Hery hery.Service
}

// ListTemplates
func (s *service) ListTemplates() {

}

// Weave
func (s *service) Weave(templatePaths string, data any) {

	funcMap := template.FuncMap{
		"hery": s.Hery.HeryFunc,
	}

	// Doc: https://pkg.go.dev/text/template#FuncMap

	tmpl, err := template.New("config").Funcs(funcMap).ParseFiles(templatePaths)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
	}
}

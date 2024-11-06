package template

import "github.com/AmadlaOrg/weaver/hery"

// NewTemplateService to set up the hery service
func NewTemplateService() ITemplate {
	return &STemplate{
		Hery: hery.NewHeryService(),
	}
}

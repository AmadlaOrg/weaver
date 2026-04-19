package template

import "github.com/AmadlaOrg/weaver/hery"

// New creates a new template service.
func New() Service {
	return &service{
		Hery: hery.New(),
	}
}

package hery

// NewHeryService to set up the hery service
func NewHeryService() IHery {
	return &SHery{}
}

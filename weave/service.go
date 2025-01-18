package weave

// NewWeaveService to set up the weave service
func NewWeaveService() IWeave {
	return &SWeave{}
}

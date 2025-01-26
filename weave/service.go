package weave

import "io"

// NewWeaveService to set up the weave service
func NewWeaveService(tmplFile string, input io.Reader, output io.Writer) IWeave {
	return &SWeave{
		tmplFile: tmplFile,
		input:    input,
		output:   &output,
	}
}

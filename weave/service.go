package weave

import "io"

// New creates a new Weaver instance.
func New(tmplFile string, input io.Reader, output io.Writer) Weaver {
	return &weaver{
		tmplFile: tmplFile,
		input:    input,
		output:   &output,
	}
}

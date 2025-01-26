package weave

import (
	"errors"
	"io"
	"text/template"
)

type IWeave interface {
	Do() error
}

type SWeave struct {
	tmplFile string
	input    io.Reader
	output   *io.Writer
}

var (
	templateParseFiles = template.ParseFiles
	ioReadAll          = io.ReadAll
)

// Do process the template with all data loaded into memory
func (s *SWeave) Do() error {
	// Open the template file
	tmpl, err := templateParseFiles(s.tmplFile)
	if err != nil {
		return err
	}

	// Read the input data
	dataBytes, err := ioReadAll(s.input)
	if err != nil {
		return err
	}

	// Parse the input data
	data, err := parseData(dataBytes)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return errors.New("no valid data found in input")
	}

	// Process the entire dataset
	for _, item := range data {
		if err := tmpl.Execute(*s.output, item); err != nil {
			return err
		}
	}

	return nil
}

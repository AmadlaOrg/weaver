package fs

import "os"

type IFs interface {
	OpenDataSource(filePath string) (*os.File, error)
	OpenOutput(filePath string, append bool) (*os.File, error)
}

type SFs struct {
	outputFilePerm os.FileMode
}

var (
	osOpen     = os.Open
	osCreate   = os.Create
	osOpenFile = os.OpenFile
)

// OpenDataSource
func (s *SFs) OpenDataSource(filePath string) (*os.File, error) {
	dataFile, err := osOpen(filePath)
	if err != nil {
		return nil, err
	}

	return dataFile, nil
}

// OpenOutput
func (s *SFs) OpenOutput(filePath string, append bool) (*os.File, error) {
	var (
		outputFile *os.File
		err        error
	)

	if append {
		// TODO: Validate that it is a file
		outputFile, err = osOpenFile(filePath, os.O_RDONLY|os.O_APPEND, 0666)
	} else {
		// TODO: Validate that the directory exists
		outputFile, err = osCreate("output.txt")
	}

	if err != nil {
		return nil, err
	}

	return outputFile, nil
}

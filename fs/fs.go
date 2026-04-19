package fs

import "os"

// FS defines the file system operations interface.
type FS interface {
	OpenDataSource(filePath string) (*os.File, error)
	OpenOutput(filePath string, append bool) (*os.File, error)
}

type fsImpl struct {
	outputFilePerm os.FileMode
}

var (
	osOpen     = os.Open
	osCreate   = os.Create
	osOpenFile = os.OpenFile
)

// OpenDataSource
func (s *fsImpl) OpenDataSource(filePath string) (*os.File, error) {
	dataFile, err := osOpen(filePath)
	if err != nil {
		return nil, err
	}

	return dataFile, nil
}

// OpenOutput
func (s *fsImpl) OpenOutput(filePath string, append bool) (*os.File, error) {
	var (
		outputFile *os.File
		err        error
	)

	if append {
		// TODO: Validate that it is a file
		outputFile, err = osOpenFile(filePath, os.O_RDONLY|os.O_APPEND, 0666)
	} else {
		// TODO: Validate that the directory exists
		outputFile, err = osCreate(filePath)
	}

	if err != nil {
		return nil, err
	}

	return outputFile, nil
}

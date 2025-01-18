package fs

import "os"

type IFs interface{}

type SFs struct{}

func (s *SFs) Mkdir(name string, perm os.FileMode) error {

}

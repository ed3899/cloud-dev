package aws

import "os"

type Vars struct {
	file *os.File
}

func NewVars() (vars *Vars, err error) {
	return
}

func (v *Vars) GetFile() (file *os.File) {
	return v.file
}

package hashicorp_vars

import "os"

type HashicorpVarsI interface {
	GetFile() *os.File
}

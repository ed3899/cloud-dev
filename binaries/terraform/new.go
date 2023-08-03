package terraform

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/tool_config"
)

type Binary struct {
	absPath string
}

func New(kumoExecAbsPath string) (binary *Binary) {
	return &Binary{
		absPath: filepath.Join(
			kumoExecAbsPath,
			dirs.DEPENDENCIES_DIR_NAME,
			tool_config.TERRAFORM_NAME,
			fmt.Sprintf("%s.exe", tool_config.TERRAFORM_NAME),
		),
	}
}

func (b *Binary) Init() (err error) {
	return
}

func (b *Binary) Up() (err error) {
	return
}

func (b *Binary) Destroy() (err error) {
	return
}

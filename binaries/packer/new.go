package packer

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/common/dirs"
	common_tool_constants "github.com/ed3899/kumo/common/tool/constants"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
)

type Binary struct {
	absPath string
}

func New(kumoExecAbsPath string) (binary *Binary) {
	return &Binary{
		absPath: filepath.Join(
			kumoExecAbsPath,
			dirs.DEPENDENCIES_DIR_NAME,
			common_tool_constants.PACKER_NAME,
			fmt.Sprintf("%s.exe", common_tool_constants.PACKER_NAME),
		),
	}
}

func (b *Binary) Init() (err error) {
	var (
		cmd         = exec.Command(b.absPath, "init", "-upgrade", ".")
		oopsBuilder = oops.
				Code("packer_init_failed")
	)

	if err = utils.RunCmdAndStream(cmd); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occured while running and streaming packer init command")
		return
	}

	return
}

func (b *Binary) Build() (err error) {
	var (
		cmd         = exec.Command(b.absPath, "build", ".")
		oopsBuilder = oops.
				Code("packer_build_failed")
	)

	if err = utils.RunCmdAndStream(cmd); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occured running and streaming packer build command")
		return
	}

	return
}

package binaries

import (
	"os/exec"

	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/manager"
	"github.com/ed3899/kumo/utils/cmd"
	"github.com/samber/oops"
)

// Returns a Packer instance. The Packer instance is used to run packer commands.
func NewPacker(_manager *manager.Manager) (*Packer, error) {
	oopsBuilder := oops.
		Code("NewPacker").
		With("manager", _manager)

	if _manager.Tool.Iota() != iota.Packer {
		err := oopsBuilder.
			Errorf("tool is not iota.Packer")
		return nil, err
	}

	return &Packer{
		Path: _manager.Path.Executable,
	}, nil
}

func (p *Packer) Init() error {
	oopsBuilder := oops.
		Code("Init").
		In("binaries").
		Tags("Packer")

	_cmd := exec.Command(p.Path, "init", "-upgrade", ".")

	err := cmd.RunCmdAndStream(_cmd)
	if err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occured while running and streaming packer init command")

		return err
	}

	return nil
}

func (p *Packer) Build() error {

	_cmd := exec.Command(p.Path, "build", ".")
	oopsBuilder := oops.
		Code("Build").
		In("binaries").
		Tags("Packer")

	err := cmd.RunCmdAndStream(_cmd)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Error occured running and streaming packer build command")
		return err
	}

	return nil
}

type Packer struct {
	Path string
}

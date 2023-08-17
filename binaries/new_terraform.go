package binaries

import (
	"os/exec"

	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/manager"
	"github.com/ed3899/kumo/utils/cmd"
	"github.com/samber/oops"
)

func NewTerraform(_manager *manager.Manager) (*Terraform, error) {
	oopsBuilder := oops.
		Code("NewTerraform").
		In("binaries").
		Tags("Terraform").
		With("manager", _manager)

	if _manager.Tool.Iota() != iota.Terraform {
		err := oopsBuilder.
			Errorf("tool is not iota.Terraform")
		return nil, err
	}

	return &Terraform{
		Path: _manager.Path.Executable,
	}, nil
}

func (t *Terraform) Init() error {
	oopsBuilder := oops.
		Code("Init").
		In("binaries").
		Tags("Terraform")

	_cmd := exec.Command(t.Path, "init")

	err := cmd.RunCmdAndStream(_cmd)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Error occured while running and streaming terraform init command")

		return err
	}

	return err
}

func (t *Terraform) Apply() error {
	_cmd := exec.Command(t.Path, "apply", "-auto-approve")
	oopsBuilder := oops.
		Code("Init").
		In("binaries").
		Tags("Terraform")

	err := cmd.RunCmdAndStream(_cmd)
	if err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occured while running and streaming terraform apply command")

		return err
	}

	return nil
}

func (t *Terraform) Destroy() error {
	oopsBuilder := oops.
		Code("Destroy").
		In("binaries").
		Tags("Terraform")

	_cmd := exec.Command(t.Path, "destroy", "-auto-approve")

	err := cmd.RunCmdAndStream(_cmd)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Error occured while running and streaming terraform destroy command")

		return err
	}

	return nil
}

type Terraform struct {
	Path string
}

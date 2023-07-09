package binz

import (
	"log"
	// "os/exec"
)

type PackerI interface {
	Build()
}

type Packer struct {
	ExecutablePath string
}

func (p *Packer) Build() {
	log.Println("Building AMI with default kumo config file")
	// cmd := exec.Command(p.ExecutablePath, "build", "kumo.config.yaml")
	// cmd.Run()
}
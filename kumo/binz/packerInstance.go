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

func (p *Packer) Build(kc KumoConfig) {
	log.Println("Building AMI...")
	log.Print(kc)

	// cmd := exec.Command(p.ExecutablePath, "build", "kumo.config.yaml")
	// cmd.Run()
}

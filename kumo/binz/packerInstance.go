package binz

import (
	"log"
	"os/exec"
)

type PackerI interface {
	Build()
}

type Packer struct {
	ExecutablePath string
}

func (p *Packer) Build(kc KumoConfig) {
	log.Println("Building AMI...")
	log.Print(p.ExecutablePath)
	cmd := exec.Command(p.ExecutablePath, "version")
	output,err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error occurred while running packer: %v", err)
	}

	log.Printf("Output: %s", string(output))
}

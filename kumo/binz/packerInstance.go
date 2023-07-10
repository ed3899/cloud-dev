package binz

import (
	"log"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type PackerI interface {
	Build()
}

type Packer struct {
	ExecutablePath string
}

func (p *Packer) Build2(ccmd *cobra.Command, args []string) {
	switch {
	case len(args) == 0:
		// Check if kumo config is present
		kc, err := GetKumoConfig()
		if err != nil {
			err = errors.Wrapf(err, "%s failed", ccmd.Name())
			log.Fatal(err)
		}
		log.Println("Building AMI...")
		cmd := exec.Command(p.ExecutablePath, "version")

		cmd.Env = append(cmd.Env, "")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Error occurred while running packer: %v", err)
		}
	case len(args) == 1:
	default:
		log.Fatalf("Invalid number of arguments: %v", args)
		log.Fatalf("Please see 'kumo build --help' for usage")
	}
}

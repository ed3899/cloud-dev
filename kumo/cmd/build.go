package cmd

import (
	"github.com/pkg/errors"
	"log"

	"github.com/ed3899/kumo/binaries"
	"github.com/spf13/cobra"
)

func GetBuildCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "build",
		Short: "Build an AMI with ready to use tools",
		Long: `Build an AMI with ready to use tools. Any AMI you build with Kumo will have a set of SSH keys downloaded
		to your root directory. Please keep these keys safe. If you lose them, you will not be able
		to SSH into your instance.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := binaries.PackerBuildWorkflow(); err != nil {
				err = errors.Wrap(err, "Error occurred running packer build workflow")
				log.Fatal(err)
			}
		},
	}

	return cc
}

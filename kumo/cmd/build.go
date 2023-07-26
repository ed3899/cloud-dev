package cmd

import (
	"log"

	"github.com/pkg/errors"

	"github.com/ed3899/kumo/binaries/workflows/packer"
	"github.com/spf13/cobra"
)

func GetBuildCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Build an AMI with ready to use tools",
		Long: `Build an AMI with ready to use tools. Any AMI you build with Kumo will have a set of SSH keys downloaded
		to your root directory. Please keep these keys safe. If you lose them, you will not be able
		to SSH into your instance.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := packer.BuildWorkflow(); err != nil {
				log.Fatal(errors.Wrap(err, "Error occurred running packer build workflow"))
			}
		},
	}
}

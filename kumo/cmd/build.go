package cmd

import (
	"log"

	"github.com/ed3899/kumo/binz"
	"github.com/spf13/cobra"
)

func GetBuildCommand(p *binz.Packer) *cobra.Command {
	cc := &cobra.Command{
		Use:   "build [ /path/to/kumo.config.yaml ]",
		Short: "Build an AMI with predefined tools",
		Long: `Build an AMI with predefined tools. Any AMI you build with Kumo will have a set of SSH keys downloaded
		to your root directory. Please keep these keys safe. If you lose them, you will not be able
		to SSH into your instance.`,
		Run: func(cmd *cobra.Command, args []string) {
			switch {
			case len(args) == 0:
				p.Build()
			case len(args) == 0:
			case len(args) == 1:
			default:
				log.Fatalf("Invalid number of arguments: %v", args)
				log.Fatalf("Please see 'kumo build --help' for usage")
			}
		},
	}

	return cc
}

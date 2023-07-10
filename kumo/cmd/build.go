package cmd

import (
	"github.com/ed3899/kumo/binz"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetBuildCommand(p *binz.Packer) *cobra.Command {
	cc := &cobra.Command{
		Use:   "build",
		Short: "Build an AMI with ready to use tools",
		Long: `Build an AMI with ready to use tools. Any AMI you build with Kumo will have a set of SSH keys downloaded
		to your root directory. Please keep these keys safe. If you lose them, you will not be able
		to SSH into your instance.`,
		Run: func(cmd *cobra.Command, args []string) {
			p.Build(viper.GetString("cloud"))
		},
	}

	return cc
}

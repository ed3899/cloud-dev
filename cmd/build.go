package cmd

import (
	"os"

	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/manager"
	"github.com/ed3899/kumo/utils/file"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Build() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Build an AMI with ready to use tools",
		Long: `Build an AMI with ready to use tools. Any AMI you build with Kumo will have a set of SSH keys downloaded
		to your root directory. Please keep these keys safe. If you lose them, you will not be able
		to SSH into your instance.`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				oopsBuilder = oops.
					Code("Build").
					In("cmd").
					Tags("cobra.Command").
					With("command", cmd.Name()).
					With("args", args)
			)

			toolManager := manager.NewManagerWith(
				os.Executable,
				iota.CloudIota(viper.GetString("Cloud")),
				iota.Packer,
			)

			isExecutablePresent := manager.IsExecutablePresentWith(
				file.IsFilePresentWith(os.Stat, os.IsNotExist),
			)(toolManager)

			if !isExecutablePresent {
				
			}
		},
	}
}

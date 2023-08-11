package cmd

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/tool"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
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
				oopsBuilder = oops.Code("Build").
					With("command", cmd.Name()).
					With("args", args)
			)

			currentExecutablePath, err := os.Executable()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to get current executable path")
				panic(err)
			}

			toolManager := new(tool.ToolManager)
			toolManager.
				SetDirInitial(
					filepath.Dir(currentExecutablePath),
				).
				SetDirRun(
					filepath.Join(
						currentExecutablePath,
						iota.Packer.Name(),
					),
				)

		},
	}
}

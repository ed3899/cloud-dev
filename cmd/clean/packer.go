package clean

import "github.com/spf13/cobra"

func Packer() *cobra.Command {
	return &cobra.Command{
		Use:   "packer",
		Short: "Cleans up packer files.",
		Long:  `Cleans up all the files created by packer under a specific cloud found in the kumo.config.yaml file.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

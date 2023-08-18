package clean

import "github.com/spf13/cobra"

func Terraform() *cobra.Command {
	return &cobra.Command{
		Use:   "terraform",
		Short: "Cleans up terraform files.",
		Long:  `Cleans up all the files created by terraform under a specific cloud found in the kumo.config.yaml file.`,
	}
}

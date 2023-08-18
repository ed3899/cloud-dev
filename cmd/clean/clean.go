package clean

import (
	"github.com/spf13/cobra"
)

func Clean() *cobra.Command {
	clean := &cobra.Command{
		Use:   "clean",
		Short: "Cleans up your environment. Doesn't delete AMIs or deployments.",
		Long:  `Cleans up all the files created by kumo under a specific cloud found in the kumo.config.yaml file. This command should be used only if you wish to start from scratch, without any state. Usually this shouldn't be necessary as the size of the files created by kumo is very small.`,
	}

	clean.AddCommand(
		Packer(),
		Terraform(),
	)

	return clean
}

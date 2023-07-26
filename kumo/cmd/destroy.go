package cmd

import (
	"log"

	"github.com/ed3899/kumo/binaries"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func GetDestroyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "Destroy your cloud environment",
		Long:  `Destroy your last deployed cloud environment. Doesn't destroy the AMI. It will also remove the SSH config file.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := binaries.TerraformDestroyWorkflow(); err != nil {
				log.Fatal(errors.Wrap(err, "Error occurred running terraform destroy workflow"))
			}
		},
	}
}

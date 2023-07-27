package cmd

import (
	"log"

	"github.com/ed3899/kumo/binaries/workflows/terraform"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
)

func GetDestroyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "Destroy your cloud environment",
		Long:  `Destroy your last deployed cloud environment. Doesn't destroy the AMI. It will also remove the SSH config file.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := terraform.DestroyWorkflow(); err != nil {
				log.Fatalf(
					"%+v",
					oops.Code("get_destroy_cmd_failed").
						Wrapf(err, "Error occurred running terraform destroy workflow"),
				)
			}
		},
	}
}

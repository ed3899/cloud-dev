package cmd

import (
	"log"

	"github.com/ed3899/kumo/binz"
	"github.com/spf13/cobra"
)

func GetDestroyCommand(p *binz.Terraform) *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "Destroy your cloud environment",
		Long:  `Destroy your last deployed cloud environment. Doesn't destroy the AMI. It will also remove the SSH config file.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			log.Println("Hello World")
		},
	}
}

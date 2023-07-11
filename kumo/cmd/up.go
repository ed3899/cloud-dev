package cmd

import (
	"log"

	"github.com/ed3899/kumo/binz"
	"github.com/spf13/cobra"
)

func GetUpCommand(p *binz.Terraform) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Deploy your AMI to the cloud",
		Long: `Deploy you cloud development environment. If no AMI is specified in the config file, Kumo will
		deploy the latest AMI built. It generates an SSH config file for you to easily SSH into your
		instances from VSCode.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Hello World")
		},
	}
}

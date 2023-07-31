package cmd

import (
	"log"

	"github.com/ed3899/kumo/workflows"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
)

func UpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Deploy your AMI to the cloud",
		Long: `Deploy you cloud development environment. If no AMI is specified in the config file, Kumo will
		deploy the latest AMI built. It generates an SSH config file for you to easily SSH into your
		instances from VSCode.`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				oopsBuilder = oops.Code("get_up_command_failed").
					With("command", cmd.Name()).
					With("args", args)
			)

			if err := workflows.Up(); err != nil {
				log.Fatalf(
					"%+v",
					oopsBuilder.
						Wrapf(err, "Error occurred while running terraform up workflow"),
				)
			}
		},
	}
}

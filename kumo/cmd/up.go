package cmd

import (
	"log"

	"github.com/ed3899/kumo/binz"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetUpCommand(t *binz.Terraform) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Deploy your AMI to the cloud",
		Long: `Deploy you cloud development environment. If no AMI is specified in the config file, Kumo will
		deploy the latest AMI built. It generates an SSH config file for you to easily SSH into your
		instances from VSCode.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get public IP
			ip, err := utils.GetPublicIp()
			if err != nil {
				err = errors.Wrap(err, "Error occurred while getting public IP, your instance will have default SSH from '0.0.0.0/0'")
				log.Println(err)
			}
			// Set allowed IP
			viper.Set("ALLOWED_ID", ip)
			log.Printf("Your public IP is %s, it will be used as the allowed IP for SSH", ip)

			t.Up(viper.GetString("cloud"))
		},
	}
}

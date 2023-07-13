package cmd

import (
	"log"

	"github.com/ed3899/kumo/binz"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetDestroyCommand(t *binz.Terraform) *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "Destroy your cloud environment",
		Long:  `Destroy your last deployed cloud environment. Doesn't destroy the AMI. It will also remove the SSH config file.`,
		Run: func(cmd *cobra.Command, args []string) {
			cloud := viper.GetString("Cloud")
			if cloud == "" {
				log.Fatal("Cloud is not set")
			}

			t.Destroy(cloud)
		},
	}
}

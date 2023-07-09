package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long: `🌩️ Kumo: Your quick and easy way to a cloud development environment.`,
}

func init() {
	var buildCmd = &cobra.Command{
		Use: "build [ /path/to/kumo.config.yaml ]",
    Short: "Build an AMI with predefined tools",
		Long: `Build an AMI with predefined tools. Any AMI you build with Kumo will have a set of SSH keys downloaded
      to your root directory. Please keep these keys safe. If you lose them, you will not be able
      to SSH into your instance.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			log.Println("Hello World")
		},
	}

	var upCmd = &cobra.Command{
		Use: "up [ /path/to/kumo.config.yaml ]",
    Short: "Deploy your AMI to the cloud",
		Long: `Deploy you cloud development environment. If no AMI is specified in the config file, Kumo will
    deploy the latest AMI built. It generates an SSH config file for you to easily SSH into your
    instances from VSCode.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			log.Println("Hello World")
		},
	}

	var destroyCmd = &cobra.Command{
		Use: "destroy [ /path/to/kumo.config.yaml ]",
    Short: "Destroy your cloud environment",
		Long: `Destroy your last deployed cloud environment. Doesn't destroy the AMI. It will also remove the SSH config file.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			log.Println("Hello World")
		},
	}

	rootCmd.AddCommand(buildCmd, upCmd, destroyCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

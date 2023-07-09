package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ed3899/kumo/binz"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var Packer *binz.Packer
var Pulumi *binz.Pulumi

var rootCmd = &cobra.Command{
	Long: `🌩️ Kumo: Your quick and easy cloud development environment.`,
}

func init() {
	// Get the binaries
	bins, err := utils.GetBinaries()
	if err != nil {
		log.Fatalf("Error occurred while getting binaries: %v", err)
	}
	packer, _ := binz.GetBinaryInstances(bins)

	var buildCmd = &cobra.Command{
		Use:   "build [ /path/to/kumo.config.yaml ]",
		Short: "Build an AMI with predefined tools",
		Long: `Build an AMI with predefined tools. Any AMI you build with Kumo will have a set of SSH keys downloaded
      to your root directory. Please keep these keys safe. If you lose them, you will not be able
      to SSH into your instance.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if kumo config is present
			kc, err := binz.GetKumoConfig()
			if err != nil {
				err = errors.Wrapf(err, "%s failed", cmd.Name())
				log.Fatal(err)
			}

			switch {
			case len(args) == 0:
				packer.Build(kc)
			case len(args) == 0:
			case len(args) == 1:
			default:
				log.Fatalf("Invalid number of arguments: %v", args)
				log.Fatalf("Please see 'kumo build --help' for usage")
			}
		},
	}

	var upCmd = &cobra.Command{
		Use:   "up [ /path/to/kumo.config.yaml ]",
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
		Use:   "destroy [ /path/to/kumo.config.yaml ]",
		Short: "Destroy your cloud environment",
		Long:  `Destroy your last deployed cloud environment. Doesn't destroy the AMI. It will also remove the SSH config file.`,
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

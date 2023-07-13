package cmd

import (
	"log"

	"github.com/ed3899/kumo/binz"
	"github.com/ed3899/kumo/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Long: `🌩️ Kumo: Your quick and easy cloud development environment.`,
}

func init() {
	// Get the binaries
	bins := binz.GetBinaries()
	// Read the config
	config.ReadKumoConfig(&config.KumoConfig{
		Name: "kumo.config",
		Type: "yaml",
		Path: ".",
	})
	// Check cloud compatibility
	cloud := viper.GetString("cloud")
	if cloud != "aws" {
		log.Fatalf("Cloud '%s' not yet supported", cloud)
	}
	// Assemble commands
	ccmds := GetAllCommands(bins)
	rootCmd.AddCommand(*ccmds...)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		err = errors.Wrapf(err, "Error occurred while running kumo")
		log.Fatal(err)
	}
}

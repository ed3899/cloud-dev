package cmd

import (
	"fmt"
	"os"

	"github.com/ed3899/kumo/binz"
	"github.com/ed3899/kumo/config"
	"github.com/ed3899/kumo/utils"
	"github.com/spf13/cobra"
)

var Packer *binz.Packer
var Pulumi *binz.Pulumi

var rootCmd = &cobra.Command{
	Long: `🌩️ Kumo: Your quick and easy cloud development environment.`,
}

func init() {
	// Get the binaries
	bins := utils.GetBinaries()
	// Read the config
	config.ReadKumoConfig()
	// Assemble commands
	ccmds := GetAllCommands(bins)
	rootCmd.AddCommand(*ccmds...)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

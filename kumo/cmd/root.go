package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ed3899/kumo/binz"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Packer *binz.Packer
var Pulumi *binz.Pulumi

var rootCmd = &cobra.Command{
	Long: `🌩️ Kumo: Your quick and easy cloud development environment.`,
}

func init() {
	// Get the binaries
	bins := utils.GetBinaries()
	packer, pulumi := binz.GetBinaryInstances(bins)

	viper.SetConfigName("kumo.config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		err = errors.Wrap(err, "Error reading config file")
		log.Fatal(err)
	}

	ccmds := []*cobra.Command{
		GetBuildCommand(packer),
		GetUpCommand(pulumi),
		GetDestroyCommand(pulumi),
	}

	rootCmd.AddCommand(ccmds...)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

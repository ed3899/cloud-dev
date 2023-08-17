package cmd

import (
	"log"

	"github.com/samber/oops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{}

func init() {
	viper.SetConfigName("kumo.config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf(
			"%+v",
			oops.
				Code("cmd-root.go-init").
				In("cmd").
				Wrapf(err, "Error occurred while reading kumo config"),
		)
	}

	// Assemble commands
	kumo := &cobra.Command{
		Use:  "kumo",
		Long: `🌩️ Your quick and easy cloud development environment.`,
	}
	kumo.AddCommand(GetCommands()...)
	rootCmd.AddCommand(kumo)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf(
			"%+v",
			oops.Code("root_cmd_execute_failed").
				Wrapf(err, "Error occurred while running kumo"),
		)
	}
}

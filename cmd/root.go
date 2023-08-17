package cmd

import (
	"log"
	"os"

	"github.com/samber/oops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{}

func init() {
	oopsBuilder := oops.
		Code("Init").
		In("cmd").
		Tags("Cobra")

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf(
			"%+v",
			oopsBuilder.
				Wrapf(err, "Error occurred while getting current working directory"),
		)
	}

	viper.SetConfigName("kumo.config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cwd)

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf(
			"%+v",
			oopsBuilder.
				Wrapf(err, "Error occurred while reading kumo config"),
		)
	}

	// Assemble commands
	kumo := &cobra.Command{
		Use:   "kumo",
		Short: "🌩️ Your quick and easy cloud development environment.",
		Long:  `🌩️ Your quick and easy cloud development environment.`,
	}
	kumo.AddCommand(GetCommands()...)
	rootCmd.AddCommand(kumo)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf(
			"%+v",
			oops.
				Code("Execute").
				In("cmd").
				Tags("Cobra").
				Wrapf(err, "Error occurred while running kumo"),
		)
	}
}

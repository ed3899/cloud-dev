package cmd

import (
	"log"
	"os"

	"github.com/samber/oops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var root = &cobra.Command{}

func init() {
	oopsBuilder := oops.
		Code("Kumo").
		In("cmd").
		Tags("Cobra", "PreRun")

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
				Wrapf(err, "Error occurred while reading config file. Make sure a kumo.config.yaml file exists in the current working directory"),
		)
	}

	root.AddCommand(Kumo())
}

func Execute() {
	err := root.Execute()
	if err != nil {
		log.Fatalf(
			"%+v",
			oops.
				Code("Execute").
				In("cmd").
				Tags("Cobra", "root").
				Wrapf(err, "Error occurred while running kumo"),
		)
	}
}

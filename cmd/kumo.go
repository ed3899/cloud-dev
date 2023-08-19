package cmd

import (
	"log"

	"github.com/samber/oops"
	"github.com/spf13/cobra"
)

func init() {
	kumo.AddCommand(*Commands()...)
}

var kumo = &cobra.Command{
	Use:     "kumo",
	Short:   "🌩️ Your quick and easy cloud development environment.",
	Long:    `🌩️ Your quick and easy cloud development environment.`,
	Version: "0.0.0",
}

func Commands() *CobraCmds {
	return &CobraCmds{
		Build(),
		Up(),
		Destroy(),
		Reset(),
	}
}

type CobraCmds []*cobra.Command

func Execute() {
	err := kumo.Execute()
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

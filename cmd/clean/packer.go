package clean

import (
	"log"
	"os"

	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/manager"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Packer() *cobra.Command {
	return &cobra.Command{
		Use:   "packer",
		Short: "Cleans up packer files.",
		Long:  `Cleans up all the files created by packer under a specific cloud found in the kumo.config.yaml file.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			oopsBuilder := oops.
				Code("Packer").
				In("cmd").
				In("clean").
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
		},
		Run: func(cmd *cobra.Command, args []string) {
			oopsBuilder := oops.
				Code("Packer").
				In("cmd").
				Tags("cobra.Command").
				With("command", *cmd).
				With("args", args)

			defer func() {
				if r := recover(); r != nil {
					err := oopsBuilder.Errorf("%v", r)
					log.Fatalf("panic: %+v", err)
				}
			}()

			_manager, err := manager.NewManager(iota.CloudIota(viper.GetString("cloud")), iota.Packer)
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to create new manager")

				panic(err)
			}

			err = _manager.Clean()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to clean up packer files")

				panic(err)
			}
		},
	}
}

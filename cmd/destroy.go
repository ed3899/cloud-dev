package cmd

import (
	"log"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/download"
	"github.com/ed3899/kumo/manager"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Destroy() *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "Destroy your cloud environment",
		Long:  `Destroy your last deployed cloud environment. Doesn't destroy the AMI. It will also remove the SSH config file.`,
		Run: func(cmd *cobra.Command, args []string) {
			oopsBuilder := oops.
				Code("Destroy").
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

			_manager, err := manager.NewManager(iota.CloudIota(viper.GetString("cloud")), iota.Terraform)
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to create new manager")

				panic(err)
			}

			if !_manager.ToolExecutableExists() {
				_download, err := download.NewDownload(_manager)
				if err != nil {
					err := oopsBuilder.
						Wrapf(err, "failed to create new download")

					panic(err)
				}

				defer _download.RemoveZip()
				defer _download.ProgressShutdown()

				err = _download.DownloadAndShowProgress()
				if err != nil {
					err := oopsBuilder.
						Wrapf(err, "failed to download")

					panic(err)
				}

				err = _download.ExtractAndShowProgress()
				if err != nil {
					err := oopsBuilder.
						Wrapf(err, "failed to extract")

					panic(err)
				}
			}

			terraform, err := binaries.NewTerraform(_manager)
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to create new terraform")

				panic(err)
			}

			err = _manager.ChDirToManagerDirRun()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to chdir to manager dir")

				panic(err)
			}
			defer _manager.ChdirToManagerDirInitial()

			err = terraform.Init()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to init")

				panic(err)
			}

			err = terraform.Destroy()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to apply")

				panic(err)
			}

			err = _manager.DeleteSshConfig()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to generate ssh config")

				panic(err)
			}
		},
	}
}

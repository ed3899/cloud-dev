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

func Build() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Build an AMI with ready to use tools",
		Long: `Build an AMI with ready to use tools.`,
		Run: func(cmd *cobra.Command, args []string) {
			oopsBuilder := oops.
				Code("Build").
				In("cmd").
				Tags("cobra.Command").
				With("command", cmd.Name()).
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

			err = _manager.SetManagerCloudCredentials()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to set manager cloud credentials")

				panic(err)
			}
			defer _manager.UnsetManagerCloudCredentials()

			err = _manager.SetPluginsEnvironmentVars()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to set plugins environment vars")

				panic(err)
			}
			defer _manager.UnsetPluginsEnvironmentVars()

			err = _manager.CreateTemplate()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to create template")

				panic(err)
			}

			template, err := _manager.ParseTemplate()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to parse template")

				panic(err)
			}
			defer _manager.DeleteTemplate()

			vars, err := _manager.CreateVars()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to create vars")

				panic(err)
			}
			defer vars.Close()

			err = template.Execute(vars, _manager.Environment)
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to execute template")

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

				_download.ProgressShutdown()
			}

			packer, err := binaries.NewPacker(_manager)
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to create new packer")

				panic(err)
			}

			err = _manager.ChDirToManagerDirRun()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to chdir to manager dir")

				panic(err)
			}
			defer _manager.ChdirToManagerDirInitial()

			err = packer.Init()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to init")

				panic(err)
			}

			err = packer.Build()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to build")

				panic(err)
			}
		},
	}
}

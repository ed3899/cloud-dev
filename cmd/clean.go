package cmd

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/manager"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Clean() *cobra.Command {
	clean := &cobra.Command{
		Use:   "clean",
		Short: "Cleans up your environment. Doesn't delete AMIs or deployments.",
		Long:  `Cleans up all the files created by kumo under all available clouds found in the kumo.config.yaml file. This command should be used only if you wish to start from scratch, without any state. Usually this shouldn't be necessary as the size of the files created by kumo is very small.`,
		Run: func(cmd *cobra.Command, args []string) {
			oopsBuilder := oops.
				Code("Clean").
				In("cmd").
				Tags("cobra.Command").
				With("command", *cmd).
				With("args", args)

			logger, _ := zap.NewProduction(
				zap.AddCaller(),
			)
			defer logger.Sync()

			defer func() {
				if r := recover(); r != nil {
					err := oopsBuilder.Errorf("%v", r)
					log.Fatalf("panic: %+v", err)
				}
			}()

			currentExecutablePath, err := os.Executable()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to get current executable path")

				panic(err)
			}

			unsuccesfulItems := make(chan *UnsuccesfulItem, 5)
			removedItems := make(chan string, 5)
			_ = &sync.WaitGroup{}

			clouds := []iota.Cloud{
				iota.Aws,
			}

			commonItems := []string{
				filepath.Join(
					currentExecutablePath,
					iota.Dependencies.Name(),
				),
			}

			// Append packer manifests
			for _, c := range clouds {
				commonItems = append(commonItems, filepath.Join(
					currentExecutablePath,
					iota.Packer.Name(),
					c.Name(),
				))
			}

			for _, c := range commonItems {
				go func(item string) {
					err := os.RemoveAll(item)
					if err != nil {
						unsuccesfulItems <- &UnsuccesfulItem{
							Item: item,
							Err:  err,
						}
					}

					removedItems <- item
				}(c)
			}

			for _, c := range clouds {
				go func(cloud iota.Cloud) {
					_manager, err := manager.NewManager(cloud, iota.Packer)
					if err != nil {
						err = oopsBuilder.
							Wrapf(err, "failed to create manager for cloud %s", cloud.Name())

						errChan <- err
					}

					err = _manager.Clean()
					if err != nil {
						err = oopsBuilder.
							Wrapf(err, "failed to clean up cloud %s", cloud.Name())

						errChan <- err
					}
				}(c)
			}

		},
	}

	return clean
}

type UnsuccesfulItem struct {
	Item string
	Err  error
}

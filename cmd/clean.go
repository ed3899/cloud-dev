package cmd

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ed3899/kumo/common/constants"
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
			errChan := make(chan error, 5)
			wg := &sync.WaitGroup{}

			clouds := []iota.Cloud{
				iota.Aws,
			}

			additionalItems := []string{
				filepath.Join(
					currentExecutablePath,
					iota.Dependencies.Name(),
				),
			}

			terraformFilePath := func(
				cloud iota.Cloud,
				filename string,
			) string {
				return filepath.Join(
					currentExecutablePath,
					iota.Terraform.Name(),
					cloud.Name(),
					filename,
				)
			}

			// Append packer manifests and terraform files to the list of items to be removed.
			for _, c := range clouds {
				additionalItems = append(
					additionalItems,
					filepath.Join(
						currentExecutablePath,
						iota.Packer.Name(),
						c.Name(),
					))

				additionalItems = append(additionalItems, terraformFilePath(c, constants.TERRAFORM_LOCK))
				additionalItems = append(additionalItems, terraformFilePath(c, constants.TERRAFORM_STATE))
				additionalItems = append(additionalItems, terraformFilePath(c, constants.TERRAFORM_BACKUP))
			}

			for _, a := range additionalItems {
				wg.Add(1)

				go func(item string) {
					defer wg.Done()
					err := os.RemoveAll(item)
					if err != nil {
						unsuccesfulItems <- &UnsuccesfulItem{
							Item: item,
							Err:  err,
						}
					}

					removedItems <- item
				}(a)
			}

			for _, c := range clouds {
				wg.Add(1)

				go func(cloud iota.Cloud) {
					defer wg.Done()
					_manager, err := manager.NewManager(cloud, iota.Packer)
					if err != nil {
						err = oopsBuilder.
							Wrapf(err, "failed to create manager for cloud %s and tool %s", cloud.Name(), iota.Packer.Name())

						errChan <- err
					}

					err = _manager.DeletePluginsDir()
					if err != nil {
						err = oopsBuilder.
							Wrapf(err, "failed to delete %s", _manager.Path.Dir.Plugins)

						unsuccesfulItems <- &UnsuccesfulItem{
							Item: _manager.Path.Dir.Plugins,
							Err:  err,
						}
					}

					removedItems <- _manager.Path.Dir.Plugins

					err = _manager.DeleteVars()
					if err != nil {
						err = oopsBuilder.
							Wrapf(err, "failed to delete %s", _manager.Path.Vars)

						unsuccesfulItems <- &UnsuccesfulItem{
							Item: _manager.Path.Vars,
							Err:  err,
						}
					}

					removedItems <- _manager.Path.Vars
				}(c)

				go func(cloud iota.Cloud) {
					defer wg.Done()
					_manager, err := manager.NewManager(cloud, iota.Terraform)
					if err != nil {
						err = oopsBuilder.
							Wrapf(err, "failed to create manager for cloud %s and tool %s", cloud.Name(), iota.Terraform.Name())

						errChan <- err
					}

					err = _manager.DeletePluginsDir()
					if err != nil {
						err = oopsBuilder.
							Wrapf(err, "failed to delete %s", _manager.Path.Dir.Plugins)

						unsuccesfulItems <- &UnsuccesfulItem{
							Item: _manager.Path.Dir.Plugins,
							Err:  err,
						}
					}

					removedItems <- _manager.Path.Dir.Plugins

					err = _manager.DeleteVars()
					if err != nil {
						err = oopsBuilder.
							Wrapf(err, "failed to delete %s", _manager.Path.Vars)

						unsuccesfulItems <- &UnsuccesfulItem{
							Item: _manager.Path.Vars,
							Err:  err,
						}
					}

					removedItems <- _manager.Path.Vars
				}(c)
			}

			wg.Wait()
			close(unsuccesfulItems)
			close(removedItems)
			close(errChan)

			for err := range errChan {
				logger.Error("failure", zap.Error(err))
			}

			for u := range unsuccesfulItems {
				logger.Error("failed to remove item", zap.String("item", u.Item), zap.Error(u.Err))
			}

			for r := range removedItems {
				logger.Info("removed item", zap.String("item", r))
			}
		},
	}

	return clean
}

type UnsuccesfulItem struct {
	Item string
	Err  error
}

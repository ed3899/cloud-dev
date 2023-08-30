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

// Returns a cobra command. The reset command is used to reset kumo state.
func Reset() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Resets kumo state. Doesn't delete AMIs or deployments. Be cautious!",
		Long:  `Deletes all the files created by kumo under all available clouds found in the kumo.config.yaml file. This command should be used only if you wish to start from scratch, without any state. Usually this shouldn't be necessary as the size of the files created by kumo is very small.`,
		Run: func(cmd *cobra.Command, args []string) {
			oopsBuilder := oops.
				Code("Reset").
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
			currentExecutableDir := filepath.Dir(currentExecutablePath)

			unsuccesfulItemsChan := make(chan *UnsuccesfulItem, 1000)
			wg := &sync.WaitGroup{}

			clouds := []iota.Cloud{
				iota.Aws,
			}

			// Delete plugins and vars directories for packer and terraform.
			for _, c := range clouds {
				wg.Add(2)

				go func(cloud iota.Cloud) {
					defer wg.Done()
					_manager, err := manager.NewManager(cloud, iota.Packer)
					if err != nil {
						// Return without handling error. We only care about deleting files. See explanation below in the
						// terraform section.
						return
					}

					err = _manager.DeletePluginsDir()
					if err != nil {
						err = oopsBuilder.
							Wrapf(err, "failed to delete %s", _manager.Path.Dir.Plugins)

						unsuccesfulItemsChan <- &UnsuccesfulItem{
							Item: _manager.Path.Dir.Plugins,
							Err:  err,
						}
					}

					err = _manager.DeleteVars()
					if err != nil {
						err = oopsBuilder.
							Wrapf(err, "failed to delete %s", _manager.Path.Vars)

						unsuccesfulItemsChan <- &UnsuccesfulItem{
							Item: _manager.Path.Vars,
							Err:  err,
						}
					}
				}(c)

				go func(cloud iota.Cloud) {
					defer wg.Done()
					_manager, err := manager.NewManager(cloud, iota.Terraform)
					if err != nil {
						// Return without handling error. We only care about deleting files. Under the hood
						// NewManager with terraform creates an environment which needs a packer manifest to
						// be present. Since we are deleting the packer manifest in the next step, we can't
						// instantiate the manager and therefore can't delete the files again.
						return
					}

					err = _manager.DeletePluginsDir()
					if err != nil {
						err = oopsBuilder.
							Wrapf(err, "failed to delete %s", _manager.Path.Dir.Plugins)

						unsuccesfulItemsChan <- &UnsuccesfulItem{
							Item: _manager.Path.Dir.Plugins,
							Err:  err,
						}
					}

					err = _manager.DeleteVars()
					if err != nil {
						err = oopsBuilder.
							Wrapf(err, "failed to delete %s", _manager.Path.Vars)

						unsuccesfulItemsChan <- &UnsuccesfulItem{
							Item: _manager.Path.Vars,
							Err:  err,
						}
					}
				}(c)
			}

			// Wait to delete additional items. A packer manifest
			// is needed to instantiate a manager with terraform.
			wg.Wait()

			additionalItems := []string{
				filepath.Join(
					currentExecutableDir,
					iota.Dependencies.Name(),
				),
			}

			packerManifestPath := func(
				cloud iota.Cloud,
				filename string,
			) string {
				return filepath.Join(
					currentExecutableDir,
					iota.Packer.Name(),
					cloud.Name(),
					filename,
				)
			}

			terraformFilePath := func(
				cloud iota.Cloud,
				filename string,
			) string {
				return filepath.Join(
					currentExecutableDir,
					iota.Terraform.Name(),
					cloud.Name(),
					filename,
				)
			}

			// Append packer manifests and terraform files to the list of items to be removed.
			for _, c := range clouds {
				additionalItems = append(
					additionalItems,
					packerManifestPath(c, constants.PACKER_MANIFEST),
					packerManifestPath(c, constants.PACKER_MANIFEST_LOCK),
					terraformFilePath(c, constants.TERRAFORM_LOCK),
					terraformFilePath(c, constants.TERRAFORM_STATE),
					terraformFilePath(c, constants.TERRAFORM_BACKUP),
				)
			}

			for _, a := range additionalItems {
				wg.Add(1)
				go func(item string) {
					defer wg.Done()

					err := os.RemoveAll(item)
					if err != nil {
						unsuccesfulItemsChan <- &UnsuccesfulItem{
							Item: item,
							Err:  err,
						}

						return
					}
				}(a)
			}

			go func() {
				defer close(unsuccesfulItemsChan)
				wg.Wait()
			}()

			for u := range unsuccesfulItemsChan {
				logger.Error("failed to remove item", zap.String("item", u.Item), zap.Error(u.Err))
			}

			logger.Info("reset completed")
		},
	}
}

type UnsuccesfulItem struct {
	Item string
	Err  error
}

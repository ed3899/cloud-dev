package manager

import (
	"fmt"
	"os"
	"sync"

	"github.com/ed3899/kumo/utils/file"
	"github.com/samber/oops"
	"go.uber.org/zap"
)

func (m *Manager) Clean() error {
	oopsBuilder := oops.
		Code("Clean").
		In("manager").
		Tags("Manager")

	logger, _ := zap.NewProduction(
		zap.AddCaller(),
	)
	defer logger.Sync()

	unsuccesfulItems := make(chan *UnsuccesfulItem, 5)
	removedItems := make(chan string, 5)
	itemsGroup := &sync.WaitGroup{}

	items := []string{
		m.Path.Vars,
		m.Path.Dir.Plugins,
		m.Path.PackerManifest,
		m.Path.Terraform.Lock,
		m.Path.Terraform.State,
		m.Path.Terraform.Backup,
	}

	for _, c := range items {
		itemsGroup.Add(1)
		go func(item string) {
			defer itemsGroup.Done()
			if file.IsFilePresent(item) {
				err := os.Remove(item)
				if err != nil {
					err = fmt.Errorf("failed to remove %s", item)

					unsuccesfulItems <- &UnsuccesfulItem{
						Item: item,
						Err:  err,
					}

					return
				}

				removedItems <- item
			}
		}(c)
	}

	go func() {
		defer close(unsuccesfulItems)
		defer close(removedItems)
		itemsGroup.Wait()
	}()

	for u := range unsuccesfulItems {
		if u != nil {
			logger.Error("failed to remove item",
				zap.String("tool", m.Tool.Name()),
				zap.String("cloud", m.Cloud.Name()),
				zap.String("item", u.Item),
				zap.Error(u.Err),
			)
		}
	}

	for r := range removedItems {
		if r != "" {
			logger.Info("removed item",
				zap.String("tool", m.Tool.Name()),
				zap.String("cloud", m.Cloud.Name()),
				zap.String("item", r),
			)
		}
	}

	return nil
}

type UnsuccesfulItem struct {
	Item string
	Err  error
}

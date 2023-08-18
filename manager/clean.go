package manager

import (
	"os"

	"github.com/ed3899/kumo/common/iota"
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
		zap.AddStacktrace(
			zap.ErrorLevel,
		),
	)
	defer logger.Sync()

	removedItems := []string{}

	if file.IsFilePresent(m.Path.Executable) {
		err := os.Remove(m.Path.Executable)
		if err != nil {
			err = oopsBuilder.
				Wrapf(err, "failed to remove %s", m.Path.Executable)

			return err
		}

		removedItems = append(removedItems, m.Path.Executable)
	}

	if file.IsFilePresent(m.Path.Vars) {
		err := os.Remove(m.Path.Vars)
		if err != nil {
			err = oopsBuilder.
				Wrapf(err, "failed to remove %s", m.Path.Vars)

			return err
		}

		removedItems = append(removedItems, m.Path.Vars)
	}

	if file.IsFilePresent(m.Path.Dir.Plugins) {
		err := os.RemoveAll(m.Path.Dir.Plugins)
		if err != nil {
			err = oopsBuilder.
				Wrapf(err, "failed to remove %s", m.Path.Dir.Plugins)

			return err
		}

		removedItems = append(removedItems, m.Path.Dir.Plugins)
	}

	switch m.Tool.Iota() {
	case iota.Packer:
		if file.IsFilePresent(m.Path.PackerManifest) {
			err := os.Remove(m.Path.PackerManifest)
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "failed to remove %s", m.Path.PackerManifest)

				return err
			}

			removedItems = append(removedItems, m.Path.PackerManifest)
		}

	case iota.Terraform:
		if file.IsFilePresent(m.Path.Terraform.Lock) {
			err := os.Remove(m.Path.Terraform.Lock)
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "failed to remove %s", m.Path.Terraform.Lock)

				return err
			}

			removedItems = append(removedItems, m.Path.Terraform.Lock)
		}

		if file.IsFilePresent(m.Path.Terraform.State) {
			err := os.Remove(m.Path.Terraform.State)
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "failed to remove %s", m.Path.Terraform.State)

				return err
			}

			removedItems = append(removedItems, m.Path.Terraform.State)
		}

		if file.IsFilePresent(m.Path.Terraform.Backup) {
			err := os.Remove(m.Path.Terraform.Backup)
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "failed to remove %s", m.Path.Terraform.Backup)

				return err
			}

			removedItems = append(removedItems, m.Path.Terraform.Backup)
		}

	default:
		err := oopsBuilder.
			Errorf("unknown tool: %#v", m.Tool)

		return err
	}

	logger.Info("cleaned up!",
		zap.String("cloud", m.Cloud.Name()),
		zap.String("tool", m.Tool.Name()),
		zap.Strings("removed items", removedItems),
	)

	return nil
}

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
	)
	defer logger.Sync()

	commonItems := []string{
		m.Path.Executable,
		m.Path.Vars,
		m.Path.Dir.Plugins,
	}
	removedItems := []string{}

	for _, i := range commonItems {
		if file.IsFilePresent(i) {
			err := os.Remove(i)
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "failed to remove %s", i)

				return err
			}

			removedItems = append(removedItems, i)
		}
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
		terraformItems := []string{
			m.Path.Terraform.Lock,
			m.Path.Terraform.State,
			m.Path.Terraform.Backup,
		}

		for _, t := range terraformItems {
			if file.IsFilePresent(t) {
				err := os.Remove(t)
				if err != nil {
					err = oopsBuilder.
						Wrapf(err, "failed to remove %s", t)

					return err
				}

				removedItems = append(removedItems, t)
			}
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

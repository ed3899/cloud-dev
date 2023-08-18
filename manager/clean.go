package manager

import (
	"os"

	"github.com/ed3899/kumo/common/iota"
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

	err := os.Remove(m.Path.Executable)
	if err != nil {
		logger.Error("failed to remove executable", zap.Error(err))
	}
	logger.Info("removed executable", zap.String("path", m.Path.Executable))

	err = os.Remove(m.Path.Vars)
	if err != nil {
		logger.Error("failed to remove vars file", zap.Error(err))
	}
	logger.Info("removed vars file", zap.String("path", m.Path.Vars))

	err = os.RemoveAll(m.Path.Dir.Plugins)
	if err != nil {
		logger.Error("failed to remove plugins directory", zap.Error(err))
	}
	logger.Info("removed plugins directory", zap.String("path", m.Path.Dir.Plugins))

	switch m.Tool.Iota() {
	case iota.Packer:
		err = os.Remove(m.Path.PackerManifest)
		if err != nil {
			logger.Error("failed to remove packer manifest file", zap.Error(err))
		}
		logger.Info("removed packer manifest file", zap.String("path", m.Path.PackerManifest))

	case iota.Terraform:
		err = os.Remove(m.Path.Terraform.Lock)
		if err != nil {
			logger.Error("failed to remove terraform lock file", zap.Error(err))
		}
		logger.Info("removed terraform lock file", zap.String("path", m.Path.Terraform.Lock))

		err = os.Remove(m.Path.Terraform.State)
		if err != nil {
			logger.Error("failed to remove terraform state file", zap.Error(err))
		}
		logger.Info("removed terraform state file", zap.String("path", m.Path.Terraform.State))

		err = os.Remove(m.Path.Terraform.Backup)
		if err != nil {
			logger.Error("failed to remove terraform backup file", zap.Error(err))
		}
		logger.Info("removed terraform backup file", zap.String("path", m.Path.Terraform.Backup))

	default:
		err = oopsBuilder.
			Errorf("unknown tool: %#v", m.Tool)

		return err
	}

	return nil
}

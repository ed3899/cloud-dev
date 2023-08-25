package manager

import (
	"os"

	"github.com/samber/oops"
	"go.uber.org/zap"
)

// Deletes the ssh config file.
func (m *Manager) DeleteSshConfig() error {
	logger, _ := zap.NewProduction(
		zap.AddCaller(),
	)
	defer logger.Sync()

	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("DeleteSshConfig")

	err := os.Remove(m.Path.Terraform.SshConfig)
	if err != nil {
		return oopsBuilder.
			Wrapf(err, "failed to delete ssh config file: %s", m.Path.Terraform.SshConfig)
	}

	logger.Info("Successfully deleted ssh config file",
		zap.String("path", m.Path.Terraform.SshConfig),
	)

	return nil
}

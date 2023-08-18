package manager

import (
	"os"

	"github.com/samber/oops"
	"go.uber.org/zap"
)

func (m *Manager) DeleteSshConfig() error {
	logger, _ := zap.NewProduction(
		zap.AddCaller(),
	)
	defer logger.Sync()

	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("DeleteSshConfig")

	err := os.Remove(m.Path.SshConfig)
	if err != nil {
		return oopsBuilder.
			Wrapf(err, "failed to delete ssh config file: %s", m.Path.SshConfig)
	}

	logger.Info("Successfully deleted ssh config file",
		zap.String("path", m.Path.SshConfig),
	)

	return nil
}

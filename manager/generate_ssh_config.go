package manager

import (
	"fmt"
	"os"

	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/utils/ip"
	"github.com/samber/oops"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Generates a ssh config file at the current working directory
func (m *Manager) GenerateSshConfig() error {
	logger, _ := zap.NewProduction(
		zap.AddCaller(),
	)
	defer logger.Sync()

	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("GenerateSshConfig")

	ip, err := ip.ReadIpFromFile(m.Path.Terraform.IpFile)
	if err != nil {
		return oopsBuilder.Wrapf(err, "failed to read ip from file")
	}

	content := fmt.Sprintf(`Host %s
	HostName %s
	IdentityFile "%s"
	User %s
	Port %d
	StrictHostKeyChecking %s
	PasswordAuthentication %s
	IdentitiesOnly %s
	LogLevel %s`,
		constants.HOST,
		ip,
		m.Path.Terraform.IdentityFile,
		viper.GetString("AMI.User"),
		constants.SSH_PORT,
		"no",
		"no",
		"yes",
		"error",
	)

	file, err := os.Create(m.Path.Terraform.SshConfig)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Error occurred while creating file %s", m.Path.Terraform.SshConfig)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Error occurred while writing to file %s", m.Path.Terraform.SshConfig)
		return err
	}

	logger.Info("Successfully generated ssh config file",
		zap.String("path", m.Path.Terraform.SshConfig),
	)
	logger.Info("Run `ssh -F ./kumossh kumo` to SSH into your instances")

	return nil
}

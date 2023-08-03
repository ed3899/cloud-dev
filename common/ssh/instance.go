package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/tool"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Contents struct {
	Host                   string
	HostName               string
	IdentityFile           string
	User                   string
	Port                   string
	StrictHostKeyChecking  string
	PasswordAuthentication string
	IdentitiesOnly         string
	LogLevel               string
}

type SshConfig struct {
	absPath  string
	contents *Contents
}

func NewSshConfig(toolSetup tool.ToolSetupI, cloudSetup cloud.CloudSetupI) (sshConfig *SshConfig, err error) {
	var (
		oopsBuilder = oops.
				Code("new_ssh_config_failed").
				With("toolSetup", toolSetup).
				With("cloudSetup", cloudSetup.GetCloudName())
		cloudDir                = filepath.Join(toolSetup.GetInitialDir(), tool.TERRAFORM_NAME, cloudSetup.GetCloudName())
		absPathToInstanceIpFile = filepath.Join(cloudDir, IP_FILE_NAME)

		cwd        string
		instanceIp string
	)

	if cwd, err = os.Getwd(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while getting current working directory")
		return
	}

	if instanceIp, err = utils.ReadIpFromFile(absPathToInstanceIpFile); err != nil {
		err = oopsBuilder.
			With("absPathToInstanceIpFile", absPathToInstanceIpFile).
			Wrapf(err, "Error occurred while reading IP address from file")
		return
	}

	sshConfig = &SshConfig{
		absPath: filepath.Join(cwd, CONFIG_NAME),
		contents: &Contents{
			Host:                   HOST,
			HostName:               instanceIp,
			IdentityFile:           filepath.Join(cloudDir, KEY_NAME),
			User:                   viper.GetString("AMI.User"),
			Port:                   strconv.Itoa(SSH_PORT),
			StrictHostKeyChecking:  "no",
			PasswordAuthentication: "no",
			IdentitiesOnly:         "yes",
			LogLevel:               "ERROR",
		},
	}

	return
}

func (sc *SshConfig) Create() (err error) {
	var (
		oopsBuilder = oops.
				Code("create_ssh_config_failed")
		content = fmt.Sprintf(`Host %s
    HostName %s
    IdentityFile "%s"
    User %s
    Port %s
    StrictHostKeyChecking %s
    PasswordAuthentication %s
    IdentitiesOnly %s
    LogLevel %s`,
			sc.contents.Host,
			sc.contents.HostName,
			sc.contents.IdentityFile,
			sc.contents.User,
			sc.contents.Port,
			sc.contents.StrictHostKeyChecking,
			sc.contents.PasswordAuthentication,
			sc.contents.IdentitiesOnly,
			sc.contents.LogLevel,
		)

		file *os.File
	)

	if file, err = os.Create(sc.absPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while creating file %s", sc.absPath)
		return
	}
	defer file.Close()

	if _, err = file.WriteString(content); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while writing to file %s", sc.absPath)
		return
	}

	return
}

func (sc *SshConfig) Remove() (err error) {
	var (
		oopsBuilder = oops.
			Code("remove_ssh_config_failed")
	)

	if err = os.Remove(sc.absPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing file %s", sc.absPath)
		return
	}

	return
}

func (sc *SshConfig) GetAbsPath() string {
	return sc.absPath
}

package templates

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type AWS_PackerEnvironment struct {
	AWS_ACCESS_KEY                     string
	AWS_SECRET_KEY                     string
	AWS_IAM_PROFILE                    string
	AWS_USER_IDS                       []string
	AWS_AMI_NAME                       string
	AWS_INSTANCE_TYPE                  string
	AWS_REGION                         string
	AWS_EC2_AMI_NAME_FILTER            string
	AWS_EC2_AMI_ROOT_DEVICE_TYPE       string
	AWS_EC2_AMI_VIRTUALIZATION_TYPE    string
	AWS_EC2_AMI_OWNERS                 []string
	AWS_EC2_SSH_USERNAME               string
	AWS_EC2_INSTANCE_USERNAME          string
	AWS_EC2_INSTANCE_USERNAME_HOME     string
	AWS_EC2_INSTANCE_USERNAME_PASSWORD string
}

func CraftAWSPackerVarsFile() (awsPackerHclPath string, err error) {
	awsEnv := &AWS_PackerEnvironment{
		AWS_ACCESS_KEY:                     viper.GetString("AWS.AccessKeyId"),
		AWS_SECRET_KEY:                     viper.GetString("AWS.SecretAccessKey"),
		AWS_IAM_PROFILE:                    viper.GetString("AWS.IamProfile"),
		AWS_USER_IDS:                       viper.GetStringSlice("AWS.UserIds"),
		AWS_AMI_NAME:                       viper.GetString("AMI.Name"),
		AWS_INSTANCE_TYPE:                  viper.GetString("AWS.EC2.Instance.Type"),
		AWS_REGION:                         viper.GetString("AWS.Region"),
		AWS_EC2_AMI_NAME_FILTER:            viper.GetString("AMI.Base.Filter"),
		AWS_EC2_AMI_ROOT_DEVICE_TYPE:       viper.GetString("AMI.Base.RootDeviceType"),
		AWS_EC2_AMI_VIRTUALIZATION_TYPE:    viper.GetString("AMI.Base.VirtualizationType"),
		AWS_EC2_AMI_OWNERS:                 viper.GetStringSlice("AMI.Base.Owners"),
		AWS_EC2_SSH_USERNAME:               viper.GetString("AMI.Base.User"),
		AWS_EC2_INSTANCE_USERNAME:          viper.GetString("AMI.User"),
		AWS_EC2_INSTANCE_USERNAME_HOME:     viper.GetString("AMI.Home"),
		AWS_EC2_INSTANCE_USERNAME_PASSWORD: viper.GetString("AMI.Password"),
	}

	genEnv := &GeneralPackerEnvironment{
		AWS_PackerEnvironment: awsEnv,
		GIT_USERNAME: 												viper.GetString("Git.Username"),
		GIT_EMAIL: 													viper.GetString("Git.Email"),
		ANSIBLE_TAGS: 												viper.GetStringSlice("AMI.Tools"),
		GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
	}

	cwd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return "", err
	}

	awsTemplatePath := filepath.Join(cwd, "templates", "AWS_PackerVarsTemplate.tmpl")

	tmpl, err := template.ParseFiles(awsTemplatePath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while parsing Packer AWS Vars template file")
		return "", err
	}

	phcldir, err := utils.GetPackerHclDirPath()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Packer HCL directory path")
		return "", err
	}

	awsPackerHclPath = filepath.Join(phcldir, "aws_ami.pkrvars.hcl")

	file, err := os.Create(awsPackerHclPath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while creating Packer AWS Vars file")
		return "", err
	}
	defer file.Close()

	err = tmpl.Execute(file, genEnv)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while executing Packer AWS Vars template file")
		return "", err
	}


	return awsPackerHclPath, nil
}

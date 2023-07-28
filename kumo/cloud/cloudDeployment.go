package cloud

import (
	"os"
	"path/filepath"

	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Credentials interface {
	Set() error
	Unset() error
}

type CloudDeployment struct {
	Credentials Credentials
	RunDir      *RunDir
}

func NewCloudDeployment(cloud Cloud, tool Tool) (cloudDeployment *CloudDeployment, err error) {
	const (
		PACKER_RUN_DIR_NAME    = "packer"
		TERRAFORM_RUN_DIR_NAME = "terraform"
		AWS_RUN_DIR_NAME       = "aws"
	)

	var (
		oopsBuilder = oops.
				Code("new_cloud_deployment_failed").
				With("cloud", cloud)

		credentials           Credentials
		pickedToolRunDirName  string
		pickedCloudRunDirName string
		initialRunDir         string
		targetRunDir          string
	)

	switch tool {
	case Packer:
		pickedToolRunDirName = PACKER_RUN_DIR_NAME

	case Terraform:
		pickedToolRunDirName = TERRAFORM_RUN_DIR_NAME

	default:
		err = oopsBuilder.
			Errorf("Tool '%v' not supported", tool)
		return
	}

	switch cloud {
	case AWS:
		pickedCloudRunDirName = AWS_RUN_DIR_NAME

		credentials = &AwsCredentials{
			AccessKeyId:     viper.GetString("AWS.AccessKeyId"),
			SecretAccessKey: viper.GetString("AWS.SecretAccessKey"),
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cloud)
		return
	}

	if initialRunDir, err = os.Getwd(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while getting current working directory")
		return
	}

	if targetRunDir, err = filepath.Abs(filepath.Join(pickedToolRunDirName, pickedCloudRunDirName)); err != nil {
		err = oopsBuilder.
			With("pickedToolRunDirName", pickedToolRunDirName).
			Wrapf(err, "Error occurred while crafting absolute path to %s", pickedCloudRunDirName)
		return
	}

	cloudDeployment = &CloudDeployment{
		Credentials: credentials,
		RunDir: &RunDir{
			initial: initialRunDir,
			target:  targetRunDir,
		},
	}

	return
}

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

type AwsCredentials struct {
	AccessKeyId     string
	SecretAccessKey string
}

func (ac *AwsCredentials) Set() (err error) {
	return
}

func (ac *AwsCredentials) Unset() (err error) {
	return
}

type RunDirs struct {
	Initial string
	Target  string
}

type CloudDeployment struct {
	Credentials Credentials
	RunDirs     *RunDirs
}

func newCloudDeployment(cloud, tool string) (cloudDeployment *CloudDeployment, err error) {
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
		initialRunDir            string
		targetRunDir             string
	)

	switch tool {
	case "packer":
		pickedToolRunDirName = PACKER_RUN_DIR_NAME

	case "terraform":
		pickedToolRunDirName = TERRAFORM_RUN_DIR_NAME

	default:
		err = oopsBuilder.
			Errorf("Tool '%v' not supported", tool)
		return
	}

	switch cloud {
	case "aws":
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
		RunDirs: &RunDirs{
			Initial: initialRunDir,
			Target:  targetRunDir,
		},
	}

	return
}

func (cd *CloudDeployment) SetRunDir() (err error) {
	var (
		oopsBuilder = oops.Code(
			"cloud_deployment_set_run_dir_failed",
		)
	)

	var (
		absPathToRunDir string
	)

	switch cd.Kind {
	case AWS:
		if absPathToRunDir, err = filepath.Abs(filepath.Join(PACKER_RUN_DIR_NAME, AWS_RUN_DIR_NAME)); err != nil {
			err = oopsBuilder.
				With("PACKER_RUN_DIR_NAME", PACKER_RUN_DIR_NAME).
				Wrapf(err, "Error occurred while crafting absolute path to %s", AWS_RUN_DIR_NAME)
			return
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cd.Kind)
		return
	}
	return
}

func (cd *CloudDeployment) UnsetRunDir() (err error) {
	return
}

package aws

import (
	"os"

	common_cloud_interfaces "github.com/ed3899/kumo/common/cloud/interfaces"
	"github.com/ed3899/kumo/common/cloud_credentials/interfaces"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Credentials struct {
	accessKeyId     string
	secretAccessKey string
}

func New(cloud common_cloud_interfaces.Cloud) (credential interfaces.Credentials) {
	return &Credentials{
		accessKeyId:     viper.GetString("AWS.AccessKeyId"),
		secretAccessKey: viper.GetString("AWS.SecretAccessKey"),
	}
}

func (c *Credentials) Set() (err error) {
	var (
		oopsBuilder = oops.
			Code("aws_credentials_set_failed")
	)

	if err = os.Setenv(AWS_ACCESS_KEY_ID_NAME, c.accessKeyId); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while setting %s environment variable", AWS_ACCESS_KEY_ID_NAME)
		return
	}

	if err = os.Setenv(AWS_SECRET_ACCESS_KEY_NAME, c.secretAccessKey); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while setting %s environment variable", AWS_SECRET_ACCESS_KEY_NAME)
		return
	}

	return
}

func (c *Credentials) Unset() (err error) {
	var (
		oopsBuilder = oops.
			Code("aws_credentials_unset_failed")
	)

	if err = os.Unsetenv(AWS_ACCESS_KEY_ID_NAME); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while unsetting %s environment variable", AWS_ACCESS_KEY_ID_NAME)
		return
	}

	if err = os.Unsetenv(AWS_SECRET_ACCESS_KEY_NAME); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while unsetting %s environment variable", AWS_SECRET_ACCESS_KEY_NAME)
		return
	}

	return
}

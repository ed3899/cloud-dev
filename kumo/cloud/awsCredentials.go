package cloud

import (
	"os"

	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type AwsCredentials struct {
	AccessKeyId     string
	SecretAccessKey string
}

const (
	AWS_ACCESS_KEY_ID     = "AWS_ACCESS_KEY_ID"
	AWS_SECRET_ACCESS_KEY = "AWS_SECRET_ACCESS_KEY"
)

func (ac *AwsCredentials) Set() (err error) {
	var (
		oopsBuilder = oops.
			Code("aws_credentials_set_failed")
	)

	if err = os.Setenv(AWS_ACCESS_KEY_ID, viper.GetString("AWS.AccessKeyId")); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while setting %s environment variable", AWS_ACCESS_KEY_ID)
		return
	}

	if err = os.Setenv(AWS_SECRET_ACCESS_KEY, viper.GetString("AWS.SecretAccessKey")); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while setting %s environment variable", AWS_SECRET_ACCESS_KEY)
		return
	}

	return
}

func (ac *AwsCredentials) Unset() (err error) {
	var (
		oopsBuilder = oops.
			Code("aws_credentials_unset_failed")
	)

	if err = os.Unsetenv(AWS_ACCESS_KEY_ID); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while unsetting %s environment variable", AWS_ACCESS_KEY_ID)
		return
	}

	if err = os.Unsetenv(AWS_SECRET_ACCESS_KEY); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while unsetting %s environment variable", AWS_SECRET_ACCESS_KEY)
		return
	}

	return
}

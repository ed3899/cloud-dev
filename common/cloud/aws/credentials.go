package aws

import (
	"os"

	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Credentials struct {
	AccessKeyId     string
	SecretAccessKey string
}

func (ac *Credentials) Set() (err error) {
	var (
		oopsBuilder = oops.
			Code("aws_credentials_set_failed")
	)

	if err = os.Setenv(AWS_ACCESS_KEY_ID_NAME, viper.GetString("AWS.AccessKeyId")); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while setting %s environment variable", AWS_ACCESS_KEY_ID_NAME)
		return
	}

	if err = os.Setenv(AWS_SECRET_ACCESS_KEY_NAME, viper.GetString("AWS.SecretAccessKey")); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while setting %s environment variable", AWS_SECRET_ACCESS_KEY_NAME)
		return
	}

	return
}

func (ac *Credentials) Unset() (err error) {
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

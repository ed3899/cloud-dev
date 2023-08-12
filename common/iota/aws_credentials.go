package iota

import (
	"os"

	"github.com/samber/oops"
)

const (
	AWS_ACCESS_KEY_ID     = "AWS_ACCESS_KEY_ID"
	AWS_SECRET_ACCESS_KEY = "AWS_SECRET_ACCESS_KEY"
)

type awsCredentials struct {
	accessKeyId     string
	secretAccessKey string
}

func (a awsCredentials) Set() (err error) {
	var (
		oopsBuilder = oops.
			Code("aws_credentials_set_failed")
	)

	if err = os.Setenv(AWS_ACCESS_KEY_ID, a.accessKeyId); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while setting %s environment variable", AWS_ACCESS_KEY_ID)
		return
	}

	if err = os.Setenv(AWS_SECRET_ACCESS_KEY, a.secretAccessKey); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while setting %s environment variable", AWS_SECRET_ACCESS_KEY)
		return
	}

	return
}

func (a awsCredentials) Unset() (err error) {
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

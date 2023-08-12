package manager

import "github.com/spf13/viper"

var (
	awsCredentials = map[string]string{
		"AWS_ACCESS_KEY_ID":     viper.GetString("AWS.AccessKeyId"),
		"AWS_SECRET_ACCESS_KEY": viper.GetString("AWS.SecretAccessKey"),
	}
)

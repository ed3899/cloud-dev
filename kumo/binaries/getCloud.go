package binaries

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func GetCloud() (cloud Cloud, err error) {
	var (
		cloudFromConfig = viper.GetString("Cloud")
	)

	switch cloudFromConfig {
	case "aws":
		cloud = AWS
	default:
		err = errors.Errorf("Cloud %s is not supported", cloudFromConfig)
	}
	return
}

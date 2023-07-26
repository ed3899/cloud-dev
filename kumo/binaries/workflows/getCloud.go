package workflows

import (
	"github.com/ed3899/kumo/binaries"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func GetCloud() (cloud binaries.Cloud, err error) {
	var (
		cloudFromConfig = viper.GetString("Cloud")
	)

	switch cloudFromConfig {
	case "aws":
		cloud = binaries.AWS
	default:
		err = errors.Errorf("Cloud %s is not supported", cloudFromConfig)
	}
	return
}

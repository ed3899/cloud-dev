package utils

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type GetCurrentCloudF = func() (cloud string, err error)

// Retrieves the current cloud from the config file, error if not set.
//
// Example with a *kumo.config* file containing:
//
//	```
//	Cloud: aws
//	```
// 	() -> ("aws", nil)
func GetCurrentCloud() (cloud string, err error) {
	cloud = viper.GetString("Cloud")
	if cloud == "" {
		err = errors.New("Cloud is not set in config file")
		return
	}
	return
}
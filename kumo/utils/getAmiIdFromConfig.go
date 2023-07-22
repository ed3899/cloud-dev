package utils

import "github.com/spf13/viper"

type GetAmiIdFromConfigF = func() (amiIdFromConfig string)

// Retrieves the AMI ID from the config file. If not set, returns an empty string.
//
// Example with a *kumo.config* file containing:
//
//	```
//	Up:
//	 AMI_Id: ami-1234567890
//	```
// 	() -> ("ami-1234567890")
func GetAmiIdFromConfig() (amiIdFromConfig string) {
	amiIdFromConfig = viper.GetString("Up.AMI_Id")
	return
}
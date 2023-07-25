package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type KumoConfig struct {
	Name string
	Type string
	Path string
}

func ReadKumoConfig(kc *KumoConfig) (err error) {
	viper.SetConfigName(kc.Name)
	viper.SetConfigType(kc.Type)
	viper.AddConfigPath(kc.Path)

	if err = viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "Error reading config file")
	}
	return
}

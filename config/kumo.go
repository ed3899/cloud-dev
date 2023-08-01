package config

import (
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type KumoConfig struct {
	Name string
	Type string
	Path string
}

func ReadKumoConfig(kc *KumoConfig) (err error) {
	var (
		oopsBuilder = oops.Code("read_kumo_config_failed").
			With("kc", kc)
	)
	viper.SetConfigName(kc.Name)
	viper.SetConfigType(kc.Type)
	viper.AddConfigPath(kc.Path)

	if err = viper.ReadInConfig(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error reading config file")
		return
	}

	return
}

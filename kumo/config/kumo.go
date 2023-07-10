package config

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type KumoConfig struct {
	Name string
	Type string
	Path string
}

func ReadKumoConfig(kc *KumoConfig) {
	viper.SetConfigName(kc.Name)
	viper.SetConfigType(kc.Type)
	viper.AddConfigPath(kc.Path)
	err := viper.ReadInConfig()
	if err != nil {
		err = errors.Wrap(err, "Error reading config file")
		log.Fatal(err)
	}
}

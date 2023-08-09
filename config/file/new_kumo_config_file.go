package file

import (
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type ConfigFile struct {
	Name string
	Type string
	Path string
}

func NewConfigFile(
	options ...Option,
) (
	kc *ConfigFile,
) {
	var (
		option Option
	)

	kc = &ConfigFile{}
	for _, option = range options {
		option(kc)
	}

	return
}

func WithName(
	name string,
) (
	option Option,
) {
	option = func(kc *ConfigFile) {
		kc.Name = name
	}

	return
}

func WithType(
	_type string,
) (
	option Option,
) {
	option = func(kc *ConfigFile) {
		kc.Type = _type
	}

	return
}

func WithPath(
	path string,
) (
	option Option,
) {
	option = func(kc *ConfigFile) {
		kc.Path = path
	}

	return
}

func (kc *ConfigFile) CallSetConfigName(setConfigName func(string)) (self *ConfigFile) {
	setConfigName(kc.Name)
	return
}

func (kc *ConfigFile) CallSetConfigType(setConfigType func(string)) (self *ConfigFile) {
	setConfigType(kc.Type)
	return
}

func (kc *ConfigFile) CallAddConfigPath(addConfigPath func(string)) (self *ConfigFile) {
	addConfigPath(kc.Path)
	return
}

func (kc *ConfigFile) CallReadInConfig(readInConfig func() error) (err error) {
	var (
		oopsBuilder = oops.
			Code("ReadInConfig").
			With("configReader", readInConfig)
	)

	if err = readInConfig(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error while calling configReader")
		return
	}

	return

}

func ReadKumoConfig(kc *ConfigFile) (err error) {
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

type Option func(*ConfigFile)

package file

import (
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type KumoConfigFile struct {
	Name string
	Type string
	Path string
}

func NewKumoConfigFile(
	options ...Option,
) (
	kc *KumoConfigFile,
) {
	var (
		option Option
	)

	kc = &KumoConfigFile{}
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
	option = func(kc *KumoConfigFile) {
		kc.Name = name
	}

	return
}

func WithType(
	_type string,
) (
	option Option,
) {
	option = func(kc *KumoConfigFile) {
		kc.Type = _type
	}

	return
}

func WithPath(
	path string,
) (
	option Option,
) {
	option = func(kc *KumoConfigFile) {
		kc.Path = path
	}

	return
}

func (kc *KumoConfigFile) CallSetConfigName(setConfigName func(string)) (self *KumoConfigFile) {
	setConfigName(kc.Name)
	return
}

func (kc *KumoConfigFile) CallSetConfigType(setConfigType func(string)) (self *KumoConfigFile) {
	setConfigType(kc.Type)
	return
}

func (kc *KumoConfigFile) CallAddConfigPath(addConfigPath func(string)) (self *KumoConfigFile) {
	addConfigPath(kc.Path)
	return
}

func (kc *KumoConfigFile) CallReadInConfig(readInConfig func() error) (err error) {
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

func ReadKumoConfig(kc *KumoConfigFile) (err error) {
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

type Option func(*KumoConfigFile)

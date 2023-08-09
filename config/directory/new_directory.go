package dir

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/config/cloud"
	"github.com/ed3899/kumo/config/tool"
	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
)

func NewDirectory(
	options ...Option,
) (
	directory *Directory,
) {

	var (
		option Option
	)

	directory = &Directory{}
	for _, option = range options {
		option(directory)
	}

	return
}

func WithPlugins[
	T interfaces.NameGetter[tool.ToolName],
	C interfaces.NameGetter[cloud.CloudName],
](
	tool T,
	cloud C,
	kumoExecAbsPath string,
) (
	option Option,
) {

	option = func(directory *Directory) {
		directory.plugins = Plugins(
			filepath.Join(
				kumoExecAbsPath,
				tool.Name().String(),
				cloud.Name().String(),
				constants.PLUGINS_DIR_NAME,
			),
		)
	}

	return
}

func WithRun[
	T interfaces.NameGetter[tool.ToolName],
	C interfaces.NameGetter[cloud.CloudName],
](
	tool T,
	cloud C,
	kumoExecAbsPath string,
) (
	option Option,
) {

	option = func(directory *Directory) {
		directory.run = Run(
			filepath.Join(
				kumoExecAbsPath,
				tool.Name().String(),
				cloud.Name().String(),
			),
		)
	}

	return
}

func WithInitial(
	kumoExecAbsPath string,
) (
	option Option,
) {

	option = func(directory *Directory) {
		directory.initial = Initial(
			kumoExecAbsPath,
		)
	}

	return
}

func (d *Directory) SetPluginsPath(
	os_Setenv func(key string, value string) error,
) (
	err error,
) {
	var (
		oopsBuilder = oops.
			Code("SetPluginsPath")
	)

	if err = os_Setenv(constants.PACKER_PLUGIN_PATH, d.plugins.String()); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to set environment variable '%s'", constants.PACKER_PLUGIN_PATH)
		return
	}

	return
}

func (d *Directory) UnsetPluginsPath(
	os_Unset func(key string) error,
) (
	err error,
) {
	var (
		oopsBuilder = oops.
			Code("UnsetPluginsPath")
	)

	if err = os_Unset(constants.PACKER_PLUGIN_PATH); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to unset environment variable '%s'", constants.PACKER_PLUGIN_PATH)
		return
	}

	return
}

func (d *Directory) GoToRunDirectory(
	os_Chdir func(dir string) error,
) (
	err error,
) {
	var (
		oopsBuilder = oops.
			Code("GoRun")
	)

	if err = os_Chdir(d.run.String()); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to change directory to '%s'", d.run.String())
		return
	}

	return
}

func (d *Directory) GoToInitialDirectory(
	os_Chdir func(dir string) error,
) (
	err error,
) {
	var (
		oopsBuilder = oops.
			Code("GoInitial")
	)

	if err = os_Chdir(d.initial.String()); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to change directory to '%s'", d.initial.String())
		return
	}

	return
}

func (p Plugins) String() (plugins string) {
	plugins = string(p)

	return
}

func (r Run) String() (run string) {
	run = string(r)

	return
}

func (i Initial) String() (initial string) {
	initial = string(i)

	return
}

type Directory struct {
	plugins Plugins
	run     Run
	initial Initial
}

type Plugins string

type Run string

type Initial string

type Option func(*Directory)

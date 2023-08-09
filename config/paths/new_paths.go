package path

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/config/cloud"
	"github.com/ed3899/kumo/config/tool"
	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
)

func NewPaths(
	options ...Option,
) (
	directory *Paths,
) {

	var (
		option Option
	)

	directory = &Paths{}
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

	option = func(directory *Paths) {
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

	option = func(directory *Paths) {
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

	option = func(directory *Paths) {
		directory.initial = Initial(
			kumoExecAbsPath,
		)
	}

	return
}

func (d *Paths) SetPluginsPath(
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

func (d *Paths) UnsetPluginsPath(
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

func (d *Paths) GoRun(
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

func (d *Paths) GoInitial(
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

type Paths struct {
	plugins Plugins
	run     Run
	initial Initial
}

type Plugins string

type Run string

type Initial string

type Option func(*Paths)

type PluginSetter interface {
	SetPluginsPath(os_Setenv func(key string, value string) error) (err error)
}

type PluginUnsetter interface {
	UnsetPluginsPath(os_Unset func(key string) error) (err error)
}

type RunChanger interface {
	GoRun(os_Chdir DirChangerF) (err error)
}

type InitialChanger interface {
	GoInitial(os_Chdir DirChangerF) (err error)
}

type DirChangerF func(dir string) error

type PathsI interface {
	PluginSetter
	PluginUnsetter
	RunChanger
	InitialChanger
}

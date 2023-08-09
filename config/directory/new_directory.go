package dir

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/config/cloud"
	"github.com/ed3899/kumo/config/tool"
	"github.com/ed3899/kumo/constants"
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
		directory.Plugins = Plugins(
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

type Plugins string
type Run string
type Initial string

type Directory struct {
	Plugins Plugins
	Run     Run
	Initial Initial
}

type Option func(*Directory)

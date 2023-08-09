package dir

import (
	"github.com/ed3899/kumo/config/tool"
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

func WithPlugins(
	tool tool.ToolI,
) (
	option Option,
) {
	option = func(directory *Directory) {
		directory.Plugins = plugins
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

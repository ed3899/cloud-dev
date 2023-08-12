package iota

import "github.com/samber/oops"

type Dirs int

const (
	Dependencies Dirs = iota
	Plugins
	Templates
)

func (d Dirs) Name() (string, error) {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Dirs").
		Code("Name")

	switch d {
	case Dependencies:
		return "dependencies", nil
	case Plugins:
		return "plugins", nil
	case Templates:
		return "templates", nil
	default:
		err := oopsBuilder.
			Errorf("unknown dir: %#v", d)

		return "", err
	}
}

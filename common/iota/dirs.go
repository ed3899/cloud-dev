package iota

import (
	"log"

	"github.com/samber/oops"
)

type Dirs int

const (
	Dependencies Dirs = iota
	Plugins
	Templates
)

func (d Dirs) Name() string {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Dirs").
		Code("Name")

	switch d {
	case Dependencies:
		return "dependencies"

	case Plugins:
		return "plugins"

	case Templates:
		return "templates"

	default:
		err := oopsBuilder.
			Errorf("unknown dir: %#v", d)

		log.Fatalf("%+v", err)

		return ""
	}
}

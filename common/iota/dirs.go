package iota

type Dirs int

const (
	Dependencies Dirs = iota
	Plugins
	Templates
)

func (d Dirs) Name() string {
	switch d {
	case Dependencies:
		return "dependencies"
	case Plugins:
		return "plugins"
	case Templates:
		return "templates"
	default:
		panic("Unknown dir")
	}
}
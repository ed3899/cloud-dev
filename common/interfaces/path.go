package interfaces

type Plugins string

func (p Plugins) String() (plugins string) {
	plugins = string(p)

	return
}

type Run string

func (r Run) String() (run string) {
	run = string(r)

	return
}

type Initial string

func (i Initial) String() (initial string) {
	initial = string(i)

	return
}

type PluginPathSetter interface {
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

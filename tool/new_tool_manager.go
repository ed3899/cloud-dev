package tool

import "github.com/ed3899/kumo/common/iota"

func (t ToolManager) SetPathTemplate(path string) ToolManager {
	t.path.template = path
	return t
}

func (t ToolManager) GetPathTemplate() string {
	return t.path.template
}

func (t ToolManager) SetPathVars(path string) ToolManager {
	t.path.vars = path
	return t
}

func (t ToolManager) GetPathVars() string {
	return t.path.vars
}

func (t ToolManager) SetDirInitial(dir string) ToolManager {
	t.dir.initial = dir
	return t
}

func (t ToolManager) GetDirInitial() string {
	return t.dir.initial
}

func (t ToolManager) SetDirRun(dir string) ToolManager {
	t.dir.run = dir
	return t
}

func (t ToolManager) GetDirRun() string {
	return t.dir.run
}

func (t ToolManager) Clone() ToolManager {
	return ToolManager{
		path: struct {
			template string
			vars     string
		}{template: t.path.template, vars: t.path.vars},
		dir: struct {
			initial string
			run     string
		}{initial: t.dir.initial, run: t.dir.run},
	}
}

type ToolManager struct {
	tool iota.Tool
	path struct {
		template string
		vars     string
	}
	dir struct {
		initial string
		run     string
	}
}

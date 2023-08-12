package manager

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/common/iota"
)

func NewManager(
	osExecutablePath string,
	cloud iota.Cloud,
	tool iota.Tool,
) Manager {
	osExecutableDir := filepath.Dir(osExecutablePath)
	cloudTemplateName, baseTemplateName := cloud.Templates()

	return Manager{
		cloud: cloud,
		tool:  tool,
		path: Path{
			packerManifest: filepath.Join(
				osExecutableDir,
				iota.Packer.Name(),
				cloud.Name(),
				constants.PACKER_MANIFEST,
			),
			template: Template{
				cloud: filepath.Join(
					osExecutableDir,
					iota.Templates.Name(),
					tool.Name(),
					cloudTemplateName,
				),
				base: filepath.Join(
					osExecutableDir,
					iota.Templates.Name(),
					tool.Name(),
					baseTemplateName,
				),
			},
			vars: filepath.Join(
				osExecutableDir,
				tool.Name(),
				cloud.Name(),
				tool.VarsName(),
			),
		},
		dir: Dir{
			initial: osExecutableDir,
			run: filepath.Join(
				osExecutableDir,
				tool.Name(),
				cloud.Name(),
			),
		},
	}
}

func (m Manager) Cloud() iota.Cloud {
	return m.cloud
}

func (m Manager) Tool() iota.Tool {
	return m.tool
}

func (m Manager) Path() Path {
	return m.path
}

func (m Manager) Dir() Dir {
	return m.dir
}

type Manager struct {
	cloud iota.Cloud
	tool  iota.Tool
	path  Path
	dir   Dir
}

type Path struct {
	packerManifest string
	vars           string
	template       Template
}

func (p Path) Template() Template {
	return p.template
}

func (p Path) Vars() string {
	return p.vars
}

type Template struct {
	cloud string
	base  string
}

func (t Template) Cloud() string {
	return t.cloud
}

func (t Template) Base() string {
	return t.base
}

type Dir struct {
	initial string
	run     string
}

func (d Dir) Initial() string {
	return d.initial
}

func (d Dir) Run() string {
	return d.run
}

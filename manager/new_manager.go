package manager

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/iota"
)

type Manager struct {
	cloud iota.Cloud
	tool  iota.Tool
	path  Path
	dirs  Dir
}

type Path struct {
	template Template
	vars     string
}

type Template struct {
	cloud string
	base  string
}

type Dir struct {
	initial string
	run     string
}

func NewManager(
	osExecutablePath string,
	cloud iota.Cloud,
	tool iota.Tool,
) Manager {
	cloudTemplateName, baseTemplateName := cloud.Templates()

	return Manager{
		cloud: cloud,
		tool:  tool,
		path: Path{
			template: Template{
				cloud: filepath.Join(
					osExecutablePath,
					tool.Name(),
					cloud.Name(),
					cloudTemplateName,
				),
				base: filepath.Join(
					osExecutablePath,
					tool.Name(),
					cloud.Name(),
					baseTemplateName,
				),
			},
			vars: filepath.Join(
				osExecutablePath,
				tool.Name(),
				cloud.Name(),
				tool.VarsName(),
			),
		},
		dirs: Dir{
			initial: osExecutablePath,
			run: filepath.Join(
				osExecutablePath,
				tool.Name(),
				cloud.Name(),
			),
		},
	}
}

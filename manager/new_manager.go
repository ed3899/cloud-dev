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
	osExecutableDir := filepath.Dir(osExecutablePath)
	cloudTemplateName, baseTemplateName := cloud.Templates()

	return Manager{
		cloud: cloud,
		tool:  tool,
		path: Path{
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
		dirs: Dir{
			initial: osExecutableDir,
			run: filepath.Join(
				osExecutableDir,
				tool.Name(),
				cloud.Name(),
			),
		},
	}
}

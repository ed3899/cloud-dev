package manager

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/common/iota"
)

type IIotaGetter[I any] interface {
	Iota() I
}

type INameGetter interface {
	Name() string
}

type ITemplateGetter[T any] interface {
	Template() T
}

type ICloud interface {
	IIotaGetter[iota.Cloud]
	INameGetter
	ITemplateGetter[iota.Template]
}

type IPluginPathEnvironmentVariableGetter interface {
	PluginPathEnvironmentVariable() string
}

type IVarsNameGetter interface {
	VarsName() string
}

type IVersionGetter interface {
	Version() string
}

type ITool interface {
	IIotaGetter[iota.Tool]
	INameGetter
	IPluginPathEnvironmentVariableGetter
	IVarsNameGetter
	IVersionGetter
}

func NewManager(
	currentExecutableDir string,
	cloud ICloud,
	tool ITool,
) *Manager {
	templatePath := func(templateName string) string {
		return filepath.Join(
			currentExecutableDir,
			iota.Templates.Name(),
			tool.Name(),
			templateName,
		)
	}

	return &Manager{
		Cloud: cloud.Iota(),
		Tool:  tool.Iota(),
		Path: &Path{
			Executable: filepath.Join(
				currentExecutableDir,
				iota.Dependencies.Name(),
				tool.Name(),
				fmt.Sprintf("%s.exe", tool.Name()),
			),
			PackerManifest: filepath.Join(
				currentExecutableDir,
				iota.Packer.Name(),
				cloud.Name(),
				constants.PACKER_MANIFEST,
			),
			Template: &Template{
				Cloud: templatePath(cloud.Template().Cloud()),
				Base:  templatePath(cloud.Template().Base()),
			},
			Vars: filepath.Join(
				currentExecutableDir,
				tool.Name(),
				cloud.Name(),
				tool.VarsName(),
			),
		},
		Dir: &Dir{
			Initial: currentExecutableDir,
			Run: filepath.Join(
				currentExecutableDir,
				tool.Name(),
				cloud.Name(),
			),
		},
	}
}

func (m *Manager) Clone() *Manager {
	return &Manager{
		Cloud: m.Cloud,
		Tool:  m.Tool,
		Path:  m.Path.Clone(),
		Dir:   m.Dir.Clone(),
	}
}

type Manager struct {
	Cloud iota.Cloud
	Tool  iota.Tool
	Path  *Path
	Dir   *Dir
}

func (p *Path) Clone() *Path {
	return &Path{
		Executable:     p.Executable,
		PackerManifest: p.PackerManifest,
		Vars:           p.Vars,
		Template:       p.Template.Clone(),
	}
}

type Path struct {
	Executable     string
	PackerManifest string
	Vars           string
	Template       *Template
}

func (t *Template) Clone() *Template {
	return &Template{
		Cloud: t.Cloud,
		Base:  t.Base,
	}
}

type Template struct {
	Cloud string
	Base  string
}

func (d *Dir) Clone() *Dir {
	return &Dir{
		Initial: d.Initial,
		Run:     d.Run,
	}
}

type Dir struct {
	Initial string
	Run     string
}

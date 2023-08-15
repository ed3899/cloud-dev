package manager

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
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

func NewManagerWith(
	osExecutable func() (string, error),
) NewManager {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("NewManager")

	newManager := func(cloud ICloud, tool ITool) (*Manager, error) {
		osExecutablePath, err := osExecutable()
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to get executable path")
			return nil, err
		}

		osExecutableDir := filepath.Dir(osExecutablePath)

		templatePath := func(templateName string) string {
			return filepath.Join(
				osExecutableDir,
				iota.Templates.Name(),
				tool.Name(),
				templateName,
			)
		}

		return &Manager{
			cloud: cloud.Iota(),
			tool:  tool.Iota(),
			path: Path{
				executable: filepath.Join(
					osExecutableDir,
					iota.Dependencies.Name(),
					tool.Name(),
					fmt.Sprintf("%s.exe", tool.Name()),
				),
				packerManifest: filepath.Join(
					osExecutableDir,
					iota.Packer.Name(),
					cloud.Name(),
					constants.PACKER_MANIFEST,
				),
				template: Template{
					cloud: templatePath(cloud.Template().Cloud()),
					base:  templatePath(cloud.Template().Base()),
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
		}, nil
	}

	return newManager
}

type NewManager func(cloud ICloud, tool ITool) (*Manager, error)

func (m *Manager) Cloud() iota.Cloud {
	return m.cloud
}

type ICloudGetter[C any] interface {
	Cloud() C
}

func (m *Manager) Tool() iota.Tool {
	return m.tool
}

type IToolGetter interface {
	Tool() iota.Tool
}

func (m *Manager) Path() Path {
	return m.path
}

type IPathGetter interface {
	Path() Path
}

func (m *Manager) Dir() Dir {
	return m.dir
}

type IDirGetter interface {
	Dir() Dir
}

func (m *Manager) Clone() *Manager {
	return &Manager{
		cloud: m.cloud,
		tool:  m.tool,
		path:  m.path.Clone(),
		dir:   m.dir.Clone(),
	}
}

type IManager interface {
	ICloudGetter[iota.Cloud]
	IToolGetter
	IPathGetter
	IDirGetter
	interfaces.IClone[Manager]
}

type Manager struct {
	cloud iota.Cloud
	tool  iota.Tool
	path  *Path
	dir   *Dir
}

func (p *Path) Executable() string {
	return p.executable
}

type IExecutableGetter interface {
	Executable() string
}

func (p *Path) PackerManifest() string {
	return p.packerManifest
}

type IPackerManifestGetter interface {
	PackerManifest() string
}

func (p *Path) Template() Template {
	return p.template
}

func (p *Path) Vars() string {
	return p.vars
}

type IVarsGetter interface {
	Vars() string
}

func (p *Path) Clone() *Path {
	return &Path{
		executable:     p.executable,
		packerManifest: p.packerManifest,
		vars:           p.vars,
		template:       p.template.Clone(),
	}
}

type IPath interface {
	IExecutableGetter
	IPackerManifestGetter
	ITemplateGetter[Template]
	IVarsGetter
	interfaces.IClone[Path]
}

type Path struct {
	executable     string
	packerManifest string
	vars           string
	template       *Template
}

func (t *Template) Cloud() string {
	return t.cloud
}

func (t *Template) Base() string {
	return t.base
}

type IBaseGetter interface {
	Base() string
}

func (t *Template) Clone() *Template {
	return &Template{
		cloud: t.cloud,
		base:  t.base,
	}
}

type ITemplate interface {
	ICloudGetter[string]
	IBaseGetter
	interfaces.IClone[Template]
}

type Template struct {
	cloud string
	base  string
}

func (d *Dir) Initial() string {
	return d.initial
}

type IInitialGetter interface {
	Initial() string
}

func (d *Dir) Run() string {
	return d.run
}

type IRunGetter interface {
	Run() string
}

func (d *Dir) Clone() *Dir {
	return &Dir{
		initial: d.initial,
		run:     d.run,
	}
}

type IDir interface {
	IInitialGetter
	IRunGetter
	interfaces.IClone[Dir]
}

type Dir struct {
	initial string
	run     string
}

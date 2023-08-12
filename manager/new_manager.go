package manager

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/ed3899/kumo/common/constants"
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
	cloud ICloud,
	tool ITool,
) Manager {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("NewManager").
		With("cloud", cloud).
		With("tool", tool)

	osExecutablePath, err := osExecutable()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get executable path")
		log.Fatalf("%+v", err)
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

	return Manager{
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
	}
}

func (m Manager) Cloud() iota.Cloud {
	return m.cloud
}

type ICloudGetter[C any] interface {
	Cloud() C
}

type ForCloudGetter func(cloudGetter ICloudGetter[iota.Cloud]) error

func (m Manager) Tool() iota.Tool {
	return m.tool
}

type IToolGetter interface {
	Tool() iota.Tool
}

func (m Manager) Path() Path {
	return m.path.(Path)
}

type IPathGetter interface {
	Path() Path
}

func (m Manager) Dir() Dir {
	return m.dir.(Dir)
}

type IDirGetter interface {
	Dir() Dir
}

type ForDirGetter func(manager IDirGetter) error

type IManager interface {
	ICloudGetter[iota.Cloud]
	IToolGetter
	IPathGetter
	IDirGetter
}

type Manager struct {
	cloud iota.Cloud
	tool  iota.Tool
	path  IPath
	dir   IDir
}

func (p Path) Executable() string {
	return p.executable
}

type IExecutableGetter interface {
	Executable() string
}

func (p Path) PackerManifest() string {
	return p.packerManifest
}

type IPackerManifestGetter interface {
	PackerManifest() string
}

func (p Path) Template() Template {
	return p.template.(Template)
}

func (p Path) Vars() string {
	return p.vars
}

type IVarsGetter interface {
	Vars() string
}

type IPath interface {
	IExecutableGetter
	IPackerManifestGetter
	ITemplateGetter[Template]
	IVarsGetter
}

type Path struct {
	executable     string
	packerManifest string
	vars           string
	template       ITemplate
}

func (t Template) Cloud() string {
	return t.cloud
}

func (t Template) Base() string {
	return t.base
}

type IBaseGetter interface {
	Base() string
}

type ITemplate interface {
	ICloudGetter[string]
	IBaseGetter
}

type Template struct {
	cloud string
	base  string
}

func (d Dir) Initial() string {
	return d.initial
}

type IInitialGetter interface {
	Initial() string
}

func (d Dir) Run() string {
	return d.run
}

type IRunGetter interface {
	Run() string
}

type IDir interface {
	IInitialGetter
	IRunGetter
}

type Dir struct {
	initial string
	run     string
}

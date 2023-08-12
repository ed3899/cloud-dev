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

type ITemplateGetter interface {
	Template() iota.Template
}

type ICloud interface {
	IIotaGetter[iota.Cloud]
	INameGetter
	ITemplateGetter
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

type ForCloudGetter func(cloudGetter ICloudGetter) error

type ICloudGetter interface {
	Cloud() iota.Cloud
}

func (m Manager) Tool() iota.Tool {
	return m.tool
}

type IToolGetter interface {
	Tool() iota.Tool
}

func (m Manager) Path() Path {
	return m.path
}

type IPathGetter interface {
	Path() Path
}

func (m Manager) Dir() Dir {
	return m.dir
}

type ForDirGetter func(manager IDirGetter) error

type IDirGetter interface {
	Dir() Dir
}

type IManager interface {
	ICloudGetter
	IToolGetter
	IPathGetter
	IDirGetter
}

type Manager struct {
	cloud iota.Cloud
	tool  iota.Tool
	path  Path
	dir   Dir
}

type Path struct {
	executable     string
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

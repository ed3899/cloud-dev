package manager

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

func NewManagerWith(
	osExecutable func() (string, error),
	cloud iota.Cloud,
	tool iota.Tool,
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
		cloud: cloud,
		tool:  tool,
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

type IDirGetter interface {
	Dir() Dir
}

type IManager interface {
	ICloudGetter
	IToolGetter
	IPathGetter
	IDirGetter
}

var (
	awsCredentials = map[string]string{
		"AWS_ACCESS_KEY_ID":     viper.GetString("AWS.AccessKeyId"),
		"AWS_SECRET_ACCESS_KEY": viper.GetString("AWS.SecretAccessKey"),
	}
)

func SetCredentialsWith(
	osSetenv func(string, string) error,
) ForSomeCloudGetterMaybe {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("SetCredentialsWith")

	forManager := func(manager ICloudGetter) error {
		managerCloudName := manager.Cloud().Name()

		switch manager.Cloud() {
		case iota.Aws:
			for key, value := range awsCredentials {
				if err := osSetenv(key, value); err != nil {
					return oopsBuilder.
						With("cloudName", managerCloudName).
						Wrapf(err, "failed to set environment variable %s to %s", key, value)
				}
			}

		default:
			return oopsBuilder.
				With("cloudName", managerCloudName).
				Errorf("unknown cloud: %#v", manager.Cloud())
		}

		return nil
	}

	return forManager
}

func UnsetCredentialsWith(
	osUnsetenv func(string) error,
) ForSomeCloudGetterMaybe {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("UnsetCredentialsWith")

	forManager := func(manager ICloudGetter) error {
		managerCloudName := manager.Cloud().Name()

		switch manager.Cloud() {
		case iota.Aws:
			for key := range awsCredentials {
				if err := osUnsetenv(key); err != nil {
					return oopsBuilder.
						With("cloudName", managerCloudName).
						Wrapf(err, "failed to unset environment variable %s", key)
				}
			}

		default:
			return oopsBuilder.
				With("cloudName", managerCloudName).
				Errorf("unknown cloud: %#v", manager.Cloud())
		}

		return nil
	}

	return forManager
}

type ForSomeCloudGetterMaybe func(cloudGetter ICloudGetter) error

func ChangeToRunDirWith(
	osChdir func(string) error,
) ForSomeDirGetterMaybe {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToRunDirWith")

	forManager := func(manager IDirGetter) error {
		if err := osChdir(manager.Dir().Run()); err != nil {
			return oopsBuilder.
				With("runDir", manager.Dir().Run()).
				Wrapf(err, "failed to change to run dir")
		}

		return nil
	}

	return forManager
}

func ChangeToInitialDirWith(
	osChdir func(string) error,
) ForSomeDirGetterMaybe {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToInitialDirWith")

	forManager := func(manager IDirGetter) error {
		if err := osChdir(manager.Dir().Initial()); err != nil {
			return oopsBuilder.
				With("initialDir", manager.Dir().Initial()).
				Wrapf(err, "failed to change to initial dir")
		}

		return nil
	}

	return forManager
}

type ForSomeDirGetterMaybe func(manager IDirGetter) error

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

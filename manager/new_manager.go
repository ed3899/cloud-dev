package manager

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

func NewManager(
	osExecutablePath string,
	cloud iota.Cloud,
	tool iota.Tool,
) (Manager, error) {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("NewManager").
		With("cloud", cloud).
		With("tool", tool).
		With("osExecutablePath", osExecutablePath)

	osExecutableDir := filepath.Dir(osExecutablePath)

	cloudTemplate, err := cloud.Template()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get cloud template")

		return Manager{}, err
	}

	cloudName, err := cloud.Name()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get cloud name")

		return Manager{}, err
	}

	packerName, err := iota.Packer.Name()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get packer name")

		return Manager{}, err
	}

	toolName, err := tool.Name()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get tool name")

		return Manager{}, err
	}

	toolVarsName, err := tool.VarsName()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get tool vars name")

		return Manager{}, err
	}

	iotaTemplatesName, err := iota.Templates.Name()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get iota templates name")

		return Manager{}, err
	}

	return Manager{
		cloud: cloud,
		tool:  tool,
		path: Path{
			packerManifest: filepath.Join(
				osExecutableDir,
				packerName,
				cloudName,
				constants.PACKER_MANIFEST,
			),
			template: Template{
				cloud: filepath.Join(
					osExecutableDir,
					iotaTemplatesName,
					toolName,
					cloudTemplate.Cloud(),
				),
				base: filepath.Join(
					osExecutableDir,
					iotaTemplatesName,
					toolName,
					cloudTemplate.Base(),
				),
			},
			vars: filepath.Join(
				osExecutableDir,
				toolName,
				cloudName,
				toolVarsName,
			),
		},
		dir: Dir{
			initial: osExecutableDir,
			run: filepath.Join(
				osExecutableDir,
				toolName,
				cloudName,
			),
		},
	}, nil
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

var (
	awsCredentials = map[string]string{
		"AWS_ACCESS_KEY_ID":     viper.GetString("AWS.AccessKeyId"),
		"AWS_SECRET_ACCESS_KEY": viper.GetString("AWS.SecretAccessKey"),
	}
)

func SetCredentialsWith(
	osSetenv func(string, string) error,
) ForManager {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("SetCredentialsWith")

	forManager := func(manager Manager) error {
		managerCloudName, err := manager.Cloud().Name()
		if err != nil {
			return oopsBuilder.
				Wrapf(err, "failed to get manager cloud name")
		}

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
) ForManager {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("UnsetCredentialsWith")

	forManager := func(manager Manager) error {
		managerCloudName, err := manager.Cloud().Name()
		if err != nil {
			return oopsBuilder.
				Wrapf(err, "failed to get manager cloud name")
		}

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

func ChangeToRunDirWith(
	osChdir func(string) error,
) ForManager {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToRunDirWith")

	forManager := func(manager Manager) error {
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
) ForManager {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToInitialDirWith")

	forManager := func(manager Manager) error {
		if err := osChdir(manager.Dir().Initial()); err != nil {
			return oopsBuilder.
				With("initialDir", manager.Dir().Initial()).
				Wrapf(err, "failed to change to initial dir")
		}

		return nil
	}

	return forManager
}

type ForManager func(manager Manager) error

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

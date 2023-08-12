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

var (
	awsCredentials = map[string]string{
		"AWS_ACCESS_KEY_ID":     viper.GetString("AWS.AccessKeyId"),
		"AWS_SECRET_ACCESS_KEY": viper.GetString("AWS.SecretAccessKey"),
	}
)

func SetCredentialsWith(
	osSetenv func(string, string) error,
	viperGetString func(string) string,
) ForManager {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("SetCredentialsWith")

	forManager := func(manager Manager) error {
		switch manager.Cloud() {
		case iota.Aws:
			for key, value := range awsCredentials {
				if err := osSetenv(key, value); err != nil {
					return oopsBuilder.
						With("cloudName", manager.Cloud().Name()).
						Wrapf(err, "failed to set environment variable %s to %s", key, value)
				}
			}

		default:
			return oopsBuilder.
				With("cloudName", manager.Cloud().Name()).
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
		switch manager.Cloud() {
		case iota.Aws:
			for key := range awsCredentials {
				if err := osUnsetenv(key); err != nil {
					return oopsBuilder.
						With("cloudName", manager.Cloud().Name()).
						Wrapf(err, "failed to unset environment variable %s", key)
				}
			}

		default:
			return oopsBuilder.
				With("cloudName", manager.Cloud().Name()).
				Errorf("unknown cloud: %#v", manager.Cloud())
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

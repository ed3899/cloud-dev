package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/manager/environment"
	"github.com/samber/oops"
)

func NewManager(
	cloud iota.Cloud,
	tool iota.Tool,
) (*Manager, error) {
	oopsBuilder := oops.
		Code("NewManager").
		In("manager").
		Tags("Manager").
		With("cloud", cloud).
		With("tool", tool)

	currentExecutablePath, err := os.Executable()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get current executable path")

		return nil, err
	}
	currentExecutableDir := filepath.Dir(currentExecutablePath)

	currentWorkingDir, err := os.Getwd()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get current working directory")

		return nil, err
	}

	templatePath := func(templateName string) string {
		return filepath.Join(
			currentExecutableDir,
			iota.Templates.Name(),
			tool.Name(),
			templateName,
		)
	}

	pathToPackerManifest := filepath.Join(
		currentExecutableDir,
		iota.Packer.Name(),
		cloud.Name(),
		constants.PACKER_MANIFEST,
	)

	_environment, err := environment.NewEnvironment(
		tool,
		cloud,
		pathToPackerManifest,
	)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to create environment")

		return nil, err
	}

	return &Manager{
		Cloud: cloud.Iota(),
		Tool:  tool.Iota(),
		Path: &Path{
			PackerManifest: pathToPackerManifest,
			Plugins: filepath.Join(
				currentExecutableDir,
				tool.Name(),
				cloud.Name(),
				tool.PluginDirs(),
			),
			Executable: filepath.Join(
				currentExecutableDir,
				iota.Dependencies.Name(),
				tool.Name(),
				fmt.Sprintf("%s.exe", tool.Name()),
			),
			Template: &Template{
				Merged: templatePath("merged"),
				Cloud:  templatePath(cloud.Template().Cloud),
				Base:   templatePath(cloud.Template().Base),
			},
			Vars: filepath.Join(
				currentExecutableDir,
				tool.Name(),
				cloud.Name(),
				tool.VarsName(),
			),
			IpFile: filepath.Join(
				currentExecutableDir,
				iota.Terraform.Name(),
				cloud.Name(),
				constants.IP_FILE_NAME,
			),
			IdentityFile: filepath.Join(
				currentExecutableDir,
				iota.Terraform.Name(),
				cloud.Name(),
				constants.KEY_NAME,
			),
			SshConfig: filepath.Join(
				currentWorkingDir,
				constants.CONFIG_NAME,
			),
			Terraform: &Terraform{
				Lock: filepath.Join(
					currentExecutableDir,
					iota.Terraform.Name(),
					cloud.Name(),
					constants.TERRAFORM_LOCK,
				),
				State: filepath.Join(
					currentExecutableDir,
					iota.Terraform.Name(),
					cloud.Name(),
					constants.TERRAFORM_STATE,
				),
				Backup: filepath.Join(
					currentExecutableDir,
					iota.Terraform.Name(),
					cloud.Name(),
					constants.TERRAFORM_BACKUP,
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
		},
		Environment: _environment,
	}, nil
}

type Manager struct {
	Cloud       iota.Cloud
	Tool        iota.Tool
	Path        *Path
	Environment any
}

type Path struct {
	PackerManifest string
	Plugins        string
	Executable     string
	Vars           string
	SshConfig      string
	IpFile         string
	IdentityFile   string
	Terraform      *Terraform
	Template       *Template
	Dir            *Dir
}

type Terraform struct {
	Lock   string
	State  string
	Backup string
}

type Template struct {
	Merged string
	Cloud  string
	Base   string
}

type Dir struct {
	Initial string
	Run     string
}

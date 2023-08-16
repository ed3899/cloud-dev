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
		In("manager")

	currentExecutablePath, err := os.Executable()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get current executable path")

		return nil, err
	}

	currentExecutableDir := filepath.Dir(currentExecutablePath)

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

	var _environment any

	switch tool {
	case iota.Packer:
		_environment, err = environment.NewPackerEnvironment(cloud)
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to create packer environment")
			return nil, err
		}

	case iota.Terraform:
		_environment, err = environment.NewTerraformEnvironment(pathToPackerManifest, cloud)
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to create terraform environment")
			return nil, err
		}

	default:
		err := oopsBuilder.
			Wrapf(err, "invalid tool: %s", tool.Name())

		return nil, err
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
			PackerManifest: pathToPackerManifest,
			Template: &Template{
				Cloud: templatePath(cloud.Template().Cloud),
				Base:  templatePath(cloud.Template().Base),
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
		Environment: _environment,
	}, nil
}

type Manager struct {
	Cloud       iota.Cloud
	Tool        iota.Tool
	Path        *Path
	Dir         *Dir
	Environment any
}

type Path struct {
	Executable     string
	PackerManifest string
	Vars           string
	Template       *Template
}

type Template struct {
	Cloud string
	Base  string
}

type Dir struct {
	Initial string
	Run     string
}

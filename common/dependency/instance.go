package dependency

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/download"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
)

type Dependency struct {
	absPathToExecutable string
	absPathToRunDir     string
	zip                 *download.Zip
}

func New(config ConfigI) (dependency *Dependency, err error) {
	var (
		oopsBuilder = oops.
				Code("new_packer_failed")
		dependenciesDirName    = config.GetDependenciesDirName()
		toolName               = config.GetToolName()
		toolVersion            = config.GetToolVersion()
		toolExecutableName     = config.GetToolExecutableName()
		toolZipName            = config.GetToolZipName()
		currentOs, currentArch = config.GetCurrentHostSpecs()
		toolUrl                = config.CreateHashicorpURL(toolName, toolVersion, currentOs, currentArch)

		absPathToKumoExecutable    string
		absPathToKumoExecutableDir string
		absPathToToolExecutable    string
		absPathToToolRunDir        string
		absPathToToolZip           string
		toolZipContentLength       int64
	)

	if absPathToKumoExecutable, err = config.GetKumoExecutableAbsPath(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to create absolute path to kumo executable")
		return
	}

	absPathToKumoExecutableDir = filepath.Dir(absPathToKumoExecutable)
	absPathToToolExecutable = filepath.Join(absPathToKumoExecutableDir, dependenciesDirName, toolName, toolExecutableName)
	absPathToToolRunDir = filepath.Join(absPathToKumoExecutableDir, toolName)
	absPathToToolZip = filepath.Join(absPathToKumoExecutableDir, dependenciesDirName, toolName, toolZipName)

	if toolZipContentLength, err = config.GetContentLength(toolUrl); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get content length for: %s", toolUrl)
		return
	}

	dependency = &Dependency{
		absPathToExecutable: absPathToToolExecutable,
		absPathToRunDir:     absPathToToolRunDir,
		zip: &download.Zip{
			Name:          toolZipName,
			AbsPath:       absPathToToolZip,
			URL:           toolUrl,
			ContentLength: toolZipContentLength,
		},
	}

	return
}

func (i *Dependency) IsInstalled() (isInstalled bool) {
	return utils.FilePresent(i.absPathToExecutable)
}

func (i *Dependency) IsNotInstalled() (isNotInstalled bool) {
	return utils.FileNotPresent(i.absPathToExecutable)
}

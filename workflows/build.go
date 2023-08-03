package workflows

import (
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/binaries/packer"
	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/download"
	common_hashicorp_vars "github.com/ed3899/kumo/common/hashicorp_vars"
	"github.com/ed3899/kumo/common/tool"
	"github.com/ed3899/kumo/hashicorp_vars"
	"github.com/ed3899/kumo/templates"
	"github.com/samber/oops"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Build() (err error) {
	var (
		oopsBuilder = oops.
				Code("build_failed")
		logger, _ = zap.NewProduction()

		packerConfig             binaries.ConfigI
		packerInstance           *packer.Instance
		cloudSetup               *cloud.Config
		toolSetup                *tool.Config
		pickedTemplate           *templates.MergedTemplate
		pickedHashicorpVars      common_hashicorp_vars.HashicorpVarsI
		uncheckedCloudFromConfig string
	)

	defer logger.Sync()

	// 0. Instantiate config
	if packerConfig, err = binaries.NewConfig(tool.Packer); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while instantiating for tool %#v", tool.Packer)
		return
	}

	// 1. Instantiate Packer
	if packerInstance, err = packer.NewInstance(packerConfig); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while instantiating Packer")
		return
	}

	// 2. Download and install if needed
	if packerInstance.IsNotInstalled() {
		if err = download.Initiate(packerInstance.Zip, filepath.Dir(packerInstance.AbsPathToExecutable)); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while downloading %s", packerInstance.Zip.GetName())
			return
		}
	}

	// 3. Cloud setup
	uncheckedCloudFromConfig = viper.GetString("Cloud")
	if cloudSetup, err = cloud.NewConfig(uncheckedCloudFromConfig); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while instantiating CloudSetup for %s", uncheckedCloudFromConfig)
		return
	}
	// a. Set cloud credentials and defer unset
	if err = cloudSetup.Credentials.Set(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while setting credentials for %s", cloudSetup.GetCloudName())
		return
	}
	defer func() {
		if err := cloudSetup.Credentials.Unset(); err != nil {
			logger.Warn(
				"Failed to unset cloud credentials",
				zap.String("error", err.Error()),
				zap.String("cloud", cloudSetup.GetCloudName()),
			)
		}
	}()

	// b. Set packer plugin paths and defer unset
	if err = packerInstance.SetPluginPath(cloudSetup); err != nil {
		err = oopsBuilder.
			With("cloudSetup.GetCloudName()", cloudSetup.GetCloudName()).
			Wrapf(err, "Error occurred while setting plugin path for packer")
		return
	}
	defer func() {
		if err := packerInstance.UnsetPluginPath(); err != nil {
			logger.Warn(
				"Failed to unset plugin path for packer",
				zap.String("error", err.Error()),
			)
		}
	}()

	// 4. Tool setup
	if toolSetup, err = tool.NewConfig(tool.Packer, cloudSetup); err != nil {
		err = oopsBuilder.
			With("tool.Packer", tool.Packer).
			With("cloudSetup.GetCloudName()", cloudSetup.GetCloudName()).
			Wrapf(err, "Error occurred while instantiating ToolSetup for packer")
		return
	}

	// 5. Pick template and defer deletion
	if pickedTemplate, err = templates.PickTemplate(toolSetup, cloudSetup); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", toolSetup.GetToolType()).
			With("cloudSetup.GetCloudType()", cloudSetup.GetCloudType()).
			Wrapf(err, "Error occurred while picking template")
		return
	}
	defer func() {
		if err := pickedTemplate.Remove(); err != nil {
			logger.Warn(
				"Failed to remove temporary template",
				zap.String("error", err.Error()),
				zap.String("template", pickedTemplate.GetName()),
			)
		}
	}()

	// 6. Pick hashicorp vars
	if pickedHashicorpVars, err = hashicorp_vars.PickHashicorpVars(toolSetup, cloudSetup); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", toolSetup.GetToolType()).
			With("cloudSetup.GetCloudType()", cloudSetup.GetCloudType()).
			Wrapf(err, "Error occurred while picking hashicorp vars")
		return
	}

	// 7. Execute template on hashicorp vars
	if err = pickedTemplate.ExecuteOn(pickedHashicorpVars); err != nil {
		err = oopsBuilder.
			With("pickedTemplate.GetName()", pickedTemplate.GetName()).
			With("pickedHashicorpVars.GetFile().Name()", pickedHashicorpVars.GetFile().Name()).
			Wrapf(err, "Error occurred while executing template on hashicorp vars")
		return
	}

	// 8. Change to right directory and defer changing back
	if err = toolSetup.GoTargetDir(); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", toolSetup.GetToolType()).
			Wrapf(err, "Error occurred while changing to target directory")
	}
	defer func() {
		if err := toolSetup.GoInitialDir(); err != nil {
			logger.Warn(
				"Failed to change back to initial directory",
				zap.String("error", err.Error()),
			)
		}
	}()

	// 9. Initialize
	if err = packerInstance.Init(); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", toolSetup.GetToolType()).
			Wrapf(err, "Error occurred while initializing packer")
		return
	}

	// 10. Build
	if err = packerInstance.Build(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while building packer")
		return
	}

	return
}

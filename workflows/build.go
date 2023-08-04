package workflows

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/binaries/packer"
	"github.com/ed3899/kumo/common/cloud_config"
	"github.com/ed3899/kumo/common/cloud_credentials"
	cloud_credentials_interfaces "github.com/ed3899/kumo/common/cloud_credentials/interfaces"
	"github.com/ed3899/kumo/common/download"
	common_hashicorp_vars "github.com/ed3899/kumo/common/hashicorp_vars"
	"github.com/ed3899/kumo/common/tool_config"
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

		cloud                    cloud_config.CloudI
		cloudCredentials         cloud_credentials_interfaces.Credentials
		kumoExecAbsPath          string
		packerConfig             *packer.Binary
		packerInstance           *packer.Instance
		tool                     tool_config.ToolI
		pickedTemplate           *templates.MergedTemplate
		pickedHashicorpVars      common_hashicorp_vars.HashicorpVarsI
		uncheckedCloudFromConfig string
	)

	defer logger.Sync()

	// Set cloud config
	uncheckedCloudFromConfig = viper.GetString("Cloud")
	if cloud, err = cloud_config.New(uncheckedCloudFromConfig); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while instantiating cloud for %s", uncheckedCloudFromConfig)
		return
	}

	// Set cloud credentials and defer unset
	if cloudCredentials, err = cloud_credentials.New(cloud); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while instantiating cloud credentials for %s", cloud.Name())
		return
	}

	if err = cloudCredentials.Set(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while setting cloud credentials for %s", cloud.Name())
		return
	}
	defer func() {
		if err := cloudCredentials.Unset(); err != nil {
			logger.Warn(
				"Failed to unset cloud credentials",
				zap.String("error", err.Error()),
				zap.String("cloud", cloud.Name()),
			)
		}
	}()

	// Get kumo executable abs path
	if kumoExecAbsPath, err = os.Executable(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while getting kumo executable abs path")
		return
	}

	// Set tool config
	if tool, err = tool_config.New(tool_config.Packer, cloud, kumoExecAbsPath); err != nil {
		err = oopsBuilder.
			With("toolKind", tool_config.Packer).
			With("cloud", cloud.Name()).
			Wrapf(err, "Error occurred while instantiating ToolSetup for packer")
		return
	}

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

	// b. Set packer plugin paths and defer unset
	if err = packerInstance.SetPluginPath(cloud); err != nil {
		err = oopsBuilder.
			With("cloudSetup.GetCloudName()", cloud.GetCloudName()).
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

	// 5. Pick template and defer deletion
	if pickedTemplate, err = templates.PickTemplate(tool, cloud); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", tool.GetToolType()).
			With("cloudSetup.GetCloudType()", cloud.GetCloudType()).
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
	if pickedHashicorpVars, err = hashicorp_vars.PickHashicorpVars(tool, cloud); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", tool.GetToolType()).
			With("cloudSetup.GetCloudType()", cloud.GetCloudType()).
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
	if err = tool.GoTargetDir(); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", tool.GetToolType()).
			Wrapf(err, "Error occurred while changing to target directory")
	}
	defer func() {
		if err := tool.GoInitialDir(); err != nil {
			logger.Warn(
				"Failed to change back to initial directory",
				zap.String("error", err.Error()),
			)
		}
	}()

	// 9. Initialize
	if err = packerInstance.Init(); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", tool.GetToolType()).
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

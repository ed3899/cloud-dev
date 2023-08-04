package workflows

import (
	"os"
	"path/filepath"

	binaries_packer "github.com/ed3899/kumo/binaries/packer"
	binaries_packer_interfaces "github.com/ed3899/kumo/binaries/packer/interfaces"
	"github.com/ed3899/kumo/common/cloud_config"
	"github.com/ed3899/kumo/common/cloud_credentials"
	cloud_credentials_interfaces "github.com/ed3899/kumo/common/cloud_credentials/interfaces"
	"github.com/ed3899/kumo/common/download"
	common_hashicorp_vars "github.com/ed3899/kumo/common/hashicorp_vars"
	common_tool "github.com/ed3899/kumo/common/tool"
	common_tool_interfaces "github.com/ed3899/kumo/common/tool/interfaces"
	common_tool_constants "github.com/ed3899/kumo/common/tool/constants"
	"github.com/ed3899/kumo/common/utils"
	common_zip "github.com/ed3899/kumo/common/zip"
	common_zip_interfaces "github.com/ed3899/kumo/common/zip/interfaces"
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

		cloud            cloud_config.CloudI
		cloudCredentials cloud_credentials_interfaces.Credentials
		kumoExecAbsPath  string

		packer 			 binaries_packer_interfaces.Packer
		tool           common_tool_interfaces.Tool
		zip            common_zip_interfaces.Zip

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
	if tool, err = common_tool.New(common_tool_constants.Packer, cloud, kumoExecAbsPath); err != nil {
		err = oopsBuilder.
			With("toolKind", common_tool_constants.Packer).
			With("cloud", cloud.Name()).
			Wrapf(err, "Error occurred while instantiating ToolSetup for packer")
		return
	}

	// Set plugin path and defer unset


	// Verify presence of tool
	if utils.FileNotPresent(tool.ExecutableName()); err != nil {

		// Instantiate zip
		if zip, err = common_zip.New(tool); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while instantiating zip for %s", tool.Name())
			return
		}

		// Download zip
		if err = download.New(zip, filepath.Dir(zip.AbsPath())); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while downloading %s", zip.Name())
			return
		}

	}

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

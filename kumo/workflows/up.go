package workflows

import (
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
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

func Up() (err error) {
	var (
		oopsBuilder = oops.
				Code("up_failed")
		logger, _ = zap.NewProduction()

		terraform                *binaries.Terraform
		cloudSetup               *cloud.CloudSetup
		toolSetup                *tool.ToolSetup
		pickedTemplate           *templates.MergedTemplate
		pickedHashicorpVars      common_hashicorp_vars.HashicorpVarsI
		uncheckedCloudFromConfig string
	)

	// 1. Instantiate Terraform
	if terraform, err = binaries.NewTerraform(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while instantiating Terraform")
		return
	}

	// 2. Download and install if needed
	if terraform.IsNotInstalled() {
		if err = download.Initiate(terraform.Zip, filepath.Dir(terraform.AbsPathToExecutable)); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while downloading %s", terraform.Zip.GetName())
			return
		}
	}

	// 3. Cloud setup
	uncheckedCloudFromConfig = viper.GetString("Cloud")
	if cloudSetup, err = cloud.NewCloudSetup(uncheckedCloudFromConfig); err != nil {
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
				"Failed to unset credentials",
				zap.String("error", err.Error()),
				zap.String("cloud", cloudSetup.GetCloudName()),
			)
		}
	}()

	// 4. Tool setup
	if toolSetup, err = tool.NewToolSetup(tool.Terraform, cloudSetup); err != nil {
		err = oopsBuilder.
			With("tool.Terraform", tool.Terraform).
			With("cloudSetup.GetCloudName()", cloudSetup.GetCloudName()).
			Wrapf(err, "Error occurred while instantiating ToolSetup for terraform")
		return
	}

	// 5. Pick template and defer deletion
	if pickedTemplate, err = templates.PickTemplate(toolSetup, cloudSetup); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", toolSetup.GetToolType()).
			With("cloudSetup.GetCloudName()", cloudSetup.GetCloudName()).
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
			With("cloudSetup.GetCloudName()", cloudSetup.GetCloudName()).
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

	// 8. Change to the right directory and defer changing back
	if err = toolSetup.GoTargetDir(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while changing to target directory")
		return
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
	if err = terraform.Init(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while initializing terraform")
		return
	}

	// 10. Apply deploy
	if err = terraform.Up(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while deploying terraform resources")
		return
	}

	return
}

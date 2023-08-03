package workflows

import (
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/binaries/terraform"
	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/download"
	common_hashicorp_vars "github.com/ed3899/kumo/common/hashicorp_vars"
	"github.com/ed3899/kumo/common/ssh"
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

		terraformConfig          binaries.ConfigI
		terraformInstance        *terraform.Instance
		cloudSetup               *cloud.Cloud
		toolSetup                tool.ConfigI
		sshConfig                ssh.SshConfigI
		pickedTemplate           *templates.MergedTemplate
		pickedHashicorpVars      common_hashicorp_vars.HashicorpVarsI
		uncheckedCloudFromConfig string
	)

	defer logger.Sync()

	// 0. Instantiate config
	if terraformConfig, err = binaries.NewConfig(tool.Terraform); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while instantiating for tool %#v", tool.Terraform)
		return
	}

	// 1. Instantiate Terraform
	if terraformInstance, err = terraform.NewInstance(terraformConfig); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while instantiating Terraform")
		return
	}

	// 2. Download and install if needed
	if terraformInstance.IsNotInstalled() {
		if err = download.Initiate(terraformInstance.Zip, filepath.Dir(terraformInstance.AbsPathToExecutable)); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while downloading %s", terraformInstance.Zip.GetName())
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
				"Failed to unset credentials",
				zap.String("error", err.Error()),
				zap.String("cloud", cloudSetup.GetCloudName()),
			)
		}
	}()

	// 4. Tool setup
	if toolSetup, err = tool.New(tool.Terraform, cloudSetup); err != nil {
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

	// 8. Change to the target directory
	if err = toolSetup.GoTargetDir(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while changing to target directory")
		return
	}

	// 9. Initialize
	if err = terraformInstance.Init(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while initializing terraform")
		return
	}

	// 10. Apply deploy
	if err = terraformInstance.Up(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while deploying terraform resources")
		return
	}

	// 11. Change back to the initial directory
	if err = toolSetup.GoInitialDir(); err != nil {
		logger.Warn(
			"Failed to change back to initial directory",
			zap.String("error", err.Error()),
		)
		err = nil
		return
	}

	// 12. Instantiate ssh config
	if sshConfig, err = ssh.NewSshConfig(toolSetup, cloudSetup); err != nil {
		logger.Warn(
			"Failed to instantiate ssh config",
			zap.String("error", err.Error()),
			zap.Any("toolSetup.GetToolType()", toolSetup.GetToolType()),
			zap.String("cloudSetup.GetCloudName()", cloudSetup.GetCloudName()),
		)
		err = nil
		return
	}

	// 12. Write ssh config
	if err = sshConfig.Create(); err != nil {
		logger.Warn(
			"Failed to write ssh config",
			zap.String("error", err.Error()),
			zap.String("sshConfig.GetAbsPath()", sshConfig.GetAbsPath()),
		)
		err = nil
		return
	}

	return
}

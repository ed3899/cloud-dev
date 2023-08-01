package workflows

import (
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/download"
	"github.com/ed3899/kumo/common/tool"
	"github.com/samber/oops"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Destroy() (err error) {
	var (
		oopsBuilder = oops.
				Code("destroy_failed")
		logger, _ = zap.NewProduction()

		terraform                *binaries.Terraform
		cloudSetup               *cloud.CloudSetup
		toolSetup                *tool.ToolSetup
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

	// 5. Change to the right directory and defer changing back
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

	// 6. Initialize
	if err = terraform.Init(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while initializing terraform")
		return
	}

	// 7. Destroy
	if err = terraform.Destroy(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while destroying terraform resources")
		return
	}

	return
}

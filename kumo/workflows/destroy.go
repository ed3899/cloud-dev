package workflows

import (
	"log"
	"path/filepath"

	"github.com/ed3899/kumo/binaries/instances"
	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/download"
	"github.com/ed3899/kumo/common/tool"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

func Destroy() (err error) {
	var (
		oopsBuilder = oops.
				Code("destroy_failed")
		terraform                *instances.Terraform
		cloudSetup               *cloud.CloudSetup
		toolSetup                *tool.ToolSetup
		uncheckedCloudFromConfig string
	)

	// 1. Instantiate Terraform
	if terraform, err = instances.NewTerraform(); err != nil {
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
	defer cloudSetup.Credentials.Unset()

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
			log.Fatalf(
				"%+v",
				oopsBuilder.
					With("toolSetup.GetToolType()", toolSetup.GetToolType()).
					Wrapf(err, "Error occurred while changing back to initial directory"),
			)
		}
	}()

	// 6. Initialize
	if err = terraform.Init(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while initializing terraform")
		return
	}

	// 7. Apply destroy
	if err = terraform.Destroy(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while destroying terraform resources")
		return
	}

	return
}

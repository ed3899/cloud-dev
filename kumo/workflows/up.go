package workflows

import (
	"log"
	"path/filepath"

	"github.com/ed3899/kumo/binaries/instances"
	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/download"
	common_hashicorp_vars "github.com/ed3899/kumo/common/hashicorp_vars"
	"github.com/ed3899/kumo/common/tool"
	"github.com/ed3899/kumo/hashicorp_vars"
	"github.com/ed3899/kumo/templates"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

func Up() (err error) {
	var (
		oopsBuilder = oops.
				Code("up_failed")
		terraform                *instances.Terraform
		cloudSetup               *cloud.CloudSetup
		toolSetup                *tool.ToolSetup
		pickedTemplate           *templates.MergedTemplate
		pickedHashicorpVars      common_hashicorp_vars.HashicorpVarsI
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

	// 5. Pick template
	if pickedTemplate, err = templates.PickTemplate(toolSetup, cloudSetup); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", toolSetup.GetToolType()).
			With("cloudSetup.GetCloudName()", cloudSetup.GetCloudName()).
			Wrapf(err, "Error occurred while picking template")
		return
	}

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
			log.Fatalf(
				"%+v",
				oopsBuilder.
					With("toolSetup.GetToolType()", toolSetup.GetToolType()).
					Wrapf(err, "Error occurred while changing back to initial directory"),
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

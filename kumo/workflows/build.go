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

func Build() (err error) {
	var (
		oopsBuilder = oops.
				Code("build_failed")

		packer                   *instances.Packer
		cloudSetup               *cloud.CloudSetup
		toolSetup                *tool.ToolSetup
		pickedTemplate           *templates.MergedTemplate
		pickedHashicorpVars      common_hashicorp_vars.HashicorpVarsI
		uncheckedCloudFromConfig string
	)

	// 1. Instantiate Packer
	if packer, err = instances.NewPacker(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while instantiating Packer")
		return
	}
	// 2. Download and install if needed
	if packer.IsNotInstalled() {
		if err = download.Initiate(packer.Zip, filepath.Dir(packer.AbsPathToExecutable)); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while downloading %s", packer.Zip.GetName())
			return
		}
	}
	// 3. CloudSetup
	uncheckedCloudFromConfig = viper.GetString("Cloud")
	if cloudSetup, err = cloud.NewCloudSetup(uncheckedCloudFromConfig); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while instantiating CloudSetup for %s", uncheckedCloudFromConfig)
		return
	}
	// a. Set credentials and defer unset
	if err = cloudSetup.Credentials.Set(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while setting credentials for %s", cloudSetup.GetCloudName())
		return
	}
	defer cloudSetup.Credentials.Unset()
	// b. Set plugin paths and defer unset
	if err = packer.SetPluginPath(cloudSetup); err != nil {
		err = oopsBuilder.
			With("cloudSetup.GetCloudName()", cloudSetup.GetCloudName()).
			Wrapf(err, "Error occurred while setting plugin path for packer")
		return
	}
	defer packer.UnsetPluginPath()
	// 4. ToolSetup
	if toolSetup, err = tool.NewToolSetup(tool.Packer, cloudSetup); err != nil {
		err = oopsBuilder.
			With("tool.Packer", tool.Packer).
			With("cloudSetup.GetCloudName()", cloudSetup.GetCloudName()).
			Wrapf(err, "Error occurred while instantiating ToolSetup for packer")
		return
	}
	// 5. Create template
	if pickedTemplate, err = templates.PickTemplate(toolSetup, cloudSetup); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", toolSetup.GetToolType()).
			With("cloudSetup.GetCloudType()", cloudSetup.GetCloudType()).
			Wrapf(err, "Error occurred while picking template")
		return
	}
	// 6. Create hashicorp vars
	if pickedHashicorpVars, err = hashicorp_vars.PickHashicorpVars(toolSetup, cloudSetup); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", toolSetup.GetToolType()).
			With("cloudSetup.GetCloudType()", cloudSetup.GetCloudType()).
			Wrapf(err, "Error occurred while picking hashicorp vars")
		return
	}
	// 7. Execute template on hashicorp vars
	if pickedTemplate.ExecuteOn(pickedHashicorpVars); err != nil {
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
			log.Fatalf(
				"%+v",
				oopsBuilder.
					With("toolSetup.GetToolType()", toolSetup.GetToolType()).
					Wrapf(err, "Error occurred while changing back to initial directory"),
			)
		}
	}()
	// 9. Initialize
	if err = packer.Init(); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", toolSetup.GetToolType()).
			Wrapf(err, "Error occurred while initializing packer")
		return
	}
	// 10. Build
	if err = packer.Build(); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", toolSetup.GetToolType()).
			Wrapf(err, "Error occurred while building packer")
		return
	}

	return
}

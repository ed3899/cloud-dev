package workflows

import (
	"path/filepath"

	"github.com/ed3899/kumo/binaries/instances"
	"github.com/ed3899/kumo/hashicorp_vars/packer"
	"github.com/ed3899/kumo/templates"
	"github.com/ed3899/kumo/common/tool"
	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/download"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

func Build() (err error) {
	var (
		oopsBuilder = oops.
			Code("build_failed")

		packer *instances.Packer
		cloudSetup *cloud.CloudSetup
		toolSetup *tool.ToolSetup
		pickedTemplate *templates.MergedTemplate
		pickedHashicorpVar any
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
	// 4. ToolSetup
	if toolSetup, err = tool.NewToolSetup(tool.Packer, cloudSetup); err != nil {
		err = oopsBuilder.
			With("tool.Packer", tool.Packer).
			With("cloudSetup", cloudSetup.GetCloudName()).
			Wrapf(err, "Error occurred while instantiating ToolSetup for packer")
		return
	}
	// 5. Create template
	if pickedTemplate ,err = templates.PickTemplate(toolSetup.GetToolType(), cloudSetup.GetCloudType()); err != nil {
		err = oopsBuilder.
			With("toolSetup.GetToolType()", toolSetup.GetToolType()).
			With("cloudSetup.GetCloudType()", cloudSetup.GetCloudType()).
			Wrapf(err, "Error occurred while picking template")
		return
	}
	// 6. Create hashicorp vars

	// 7. Change to right directory and defer change back

	// 8. Set plugin path

	// 9. Initialize

	// 10. Build

	return
}

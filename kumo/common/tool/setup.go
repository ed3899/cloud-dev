package tool

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/common/cloud"
	"github.com/samber/oops"
)

type ToolSetup struct {
	toolType   ToolType
	initialDir string
	toolDir  string
}

func (ts *ToolSetup) GetToolType() (toolType ToolType) {
	return ts.toolType
}

func (ts *ToolSetup) GoInitialDir() (err error) {
	log.Printf("Changing directory to %s", ts.initialDir)
	var (
		oopsBuilder = oops.
			Code("go_initial_dir_failed")
	)

	if err = os.Chdir(ts.initialDir); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while changing directory to %s", ts.initialDir)
		return
	}

	return
}

func (ts *ToolSetup) GoTargetDir() (err error) {
	log.Printf("Changing directory to %s", ts.toolDir)
	var (
		oopsBuilder = oops.
			Code("go_target_dir_failed")
	)

	if err = os.Chdir(ts.toolDir); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while changing directory to %s", ts.toolDir)
		return
	}

	return
}

func NewToolSetup(tool ToolType, cloud cloud.CloudSetupI) (toolSetup *ToolSetup, err error) {
	const (
		PACKER_RUN_DIR_NAME    = "packer"
		TERRAFORM_RUN_DIR_NAME = "terraform"
	)

	var (
		oopsBuilder = oops.
				Code("new_tool_setup_failed").
				With("tool", tool)
		cwd string
	)

	if cwd, err = os.Getwd(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while getting current working directory")
		return
	}

	switch tool {
	case Packer:
		toolSetup = &ToolSetup{
			toolType:   Packer,
			initialDir: cwd,
			toolDir:  filepath.Join(PACKER_RUN_DIR_NAME, cloud.GetCloudName()),
		}

	case Terraform:
		toolSetup = &ToolSetup{
			toolType:   Terraform,
			initialDir: cwd,
			toolDir:  filepath.Join(TERRAFORM_RUN_DIR_NAME, cloud.GetCloudName()),
		}

	default:
		err = oopsBuilder.
			Errorf("Tool '%v' not supported", tool)
		return
	}

	return
}

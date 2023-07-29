package aws

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/templates"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Template struct {
	instance      *template.Template
	parentDirName string
	environment   templates.EnvironmentI
}

func NewTemplate(packerManifest templates.PackerManifestI) (newTemplate *Template, err error) {
	const (
		NAME = "AWS_TerraformVars.tmpl"
	)

	var (
		oopsBuilder = oops.
				Code("new_template_failed")

		absPath  string
		instance *template.Template
	)

	if absPath, err = filepath.Abs(filepath.Join(dirs.TERRAFORM_DIR_NAME, dirs.AWS_DIR_NAME, NAME)); err != nil {
		err = oopsBuilder.
			With("dirs.PACKER_DIR_NAME", dirs.PACKER_DIR_NAME).
			With("dirs.AWS_DIR_NAME", dirs.AWS_DIR_NAME).
			Wrapf(err, "Error occurred while crafting absolute path to %s", NAME)
		return
	}

	if instance, err = template.ParseFiles(absPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPath)
		return
	}

	newTemplate = &Template{
		instance:      instance,
		parentDirName: dirs.TERRAFORM_DIR_NAME,
		environment: &Environment{
			AWS_REGION:                   viper.GetString("AWS.Region"),
			AWS_INSTANCE_TYPE:            viper.GetString("AWS.EC2.Instance.Type"),
			AWS_EC2_INSTANCE_VOLUME_TYPE: viper.GetString("AWS.EC2.Volume.Type"),
			AWS_EC2_INSTANCE_VOLUME_SIZE: viper.GetInt("AWS.EC2.Volume.Size"),
			AMI_ID:                       packerManifest.GetLastBuiltAmiId(),
		},
	}

	return
}

func (t *Template) GetParentDirName() (dir string) {
	return t.parentDirName
}

func (t *Template) GetName() (name string) {
	return t.instance.Name()
}

func (t *Template) GetInstance() (instance *template.Template) {
	return t.instance
}

func (t *Template) GetEnvironment() (environment templates.EnvironmentI) {
	return t.environment
}

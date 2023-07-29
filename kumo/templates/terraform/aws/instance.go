package aws

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/templates"
	"github.com/ed3899/kumo/templates/packer"
	"github.com/ed3899/kumo/templates/terraform"
	"github.com/samber/oops"
)

type Environment struct {
	AWS_REGION                   string
	AWS_INSTANCE_TYPE            string
	AWS_EC2_INSTANCE_VOLUME_TYPE string
	AWS_EC2_INSTANCE_VOLUME_SIZE int
	AMI_ID                       string
}

func (e *Environment) IsEnvironment() (isEnvironment bool) {
	return true
}

type Template struct {
	instance      *template.Template
	parentDirName string
	environment   *Environment
}

func NewTemplate(t *Template) (newTemplate *Template, err error) {
	const (
		NAME = "AWS_TerraformVars.tmpl"
	)

	var (
		oopsBuilder = oops.
				Code("new_template_failed")

		instance *template.Template
		absPath  string
		pickedIp string
	)

	if absPath, err = filepath.Abs(filepath.Join(terraform.TERRAFORM_TEMPLATES_DIR_NAME, templates.AWS_TEMPLATES_DIR_NAME, NAME)); err != nil {
		err = oopsBuilder.
			With("packer.PACKER_TEMPLATES_DIR_NAME", packer.PACKER_TEMPLATES_DIR_NAME).
			With("templates.AWS_TEMPLATES_DIR_NAME", templates.AWS_TEMPLATES_DIR_NAME).
			Wrapf(err, "Error occurred while crafting absolute path to %s", NAME)
		return
	}

	if instance, err = template.ParseFiles(absPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPath)
		return
	}

	return
}

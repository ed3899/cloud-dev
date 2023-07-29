package general

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/templates"
	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
	"go.uber.org/zap"
)

type Template struct {
	instance      *template.Template
	parentDirName string
	environment   *Environment
}

func NewTemplate() (newTemplate *Template, err error) {
	const (
		NAME       = "GeneralTerraformTfVars.tmpl"
		DEFAULT_IP = "0.0.0.0"
	)

	var (
		oopsBuilder = oops.
				Code("new_template_failed")
		logger, _ = zap.NewProduction()

		absPath  string
		instance *template.Template
		publicIp string
		pickedIp string
	)

	defer logger.Sync()

	if absPath, err = filepath.Abs(filepath.Join(templates.TERRAFORM_DIR_NAME, templates.GENERAL_DIR_NAME, NAME)); err != nil {
		err = oopsBuilder.
			With("templates.TERRAFORM_DIR_NAME", templates.TERRAFORM_DIR_NAME).
			With("templates.GENERAL_DIR_NAME", templates.GENERAL_DIR_NAME).
			Wrapf(err, "Error occurred while crafting absolute path to %s", NAME)
		return
	}

	if instance, err = template.ParseFiles(absPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPath)
		return
	}

	if publicIp, err = utils.GetPublicIp(); err != nil {
		logger.Sugar().Warnf("Failed to get public IP address. Defaulting to %s", DEFAULT_IP)
		pickedIp = DEFAULT_IP
	} else {
		pickedIp = publicIp
	}

	newTemplate = &Template{
		instance:      instance,
		parentDirName: templates.TERRAFORM_DIR_NAME,
		environment: &Environment{
			ALLOWED_IP: utils.MaskIp(pickedIp, 32),
		},
	}

	return

}

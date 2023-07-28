package templates

import (
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
)

type Environment interface {
	IsEnvironment() bool
}

type PackerMergedEnvironment[E Environment] struct {
	EnvironmentOne E
	EnvironmentTwo   E
}

type MergedTemplate struct {
	instance    *template.Template
	environment *PackerMergedEnvironment[Environment]
}

type TemplateI interface {
	GetName() string
	GetInstance() *template.Template
	GetEnvironment() Environment
}

func NewMergedTemplate(templateOne, templateTwo TemplateI) (packerMergedTemplate *MergedTemplate, err error) {
	const (
		TEMPLATE_DIR_NAME = "templates"
		PACKER_MERGED_TEMPLATE_NAME = "temp_merged_packer_template"
	)

	var (
		oopsBuilder = oops.
				Code("new_packer_merged_template_failed").
				With("templateOne", templateOne.GetName()).
				With("templateTwo", templateTwo.GetName())

		packerMergedTemplateInstance      *template.Template
		absPathToTemplatesDir             string
		absPathToTempPackerMergedTemplate string
	)

	if absPathToTemplatesDir, err = filepath.Abs(TEMPLATE_DIR_NAME); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while crafting absolute path to %s", binaries.TEMPLATE_DIR_NAME)
		return
	}

	absPathToTempPackerMergedTemplate = filepath.Join(absPathToTemplatesDir, PACKER_SUBDIR_NAME, PACKER_MERGED_TEMPLATE_NAME)

	if err = utils.MergeFilesTo(
		absPathToTempPackerMergedTemplate,
		generalPackerTemplate.GetName(),
		packerCloudTemplate.GetName(),
	); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while merging %s and %s to %s", generalPackerTemplate.GetName(), packerCloudTemplate.GetName(), absPathToTempPackerMergedTemplate)
	}

	if packerMergedTemplateInstance, err = template.ParseFiles(absPathToTempPackerMergedTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPathToTempPackerMergedTemplate)
		return
	}

	packerMergedTemplate = &MergedTemplate{
		instance: packerMergedTemplateInstance,
		environment: &PackerMergedEnvironment[PackerCloudEnvironment]{
			PackerGeneralEnvironment: generalPackerTemplate.GetEnvironment(),
			PackerCloudEnvironment:   packerCloudTemplate.GetEnvironment(),
		},
	}

	return
}

func (pmt *MergedTemplate) GetName() (name string) {
	return pmt.instance.Name()
}

func (pmt *MergedTemplate) GetInstance() (instance *template.Template) {
	return pmt.instance
}

func (pmt *MergedTemplate) GetEnvironment() (environment *PackerMergedEnvironment[PackerCloudEnvironment]) {
	return pmt.environment
}

func (pmt *MergedTemplate) Remove() (err error) {
	var (
		oopsBuilder = oops.
			Code("template_remove_failed")
	)

	if os.RemoveAll(pmt.instance.Name()); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", pmt.instance.Name())
		return
	}

	return
}

func (pmt *MergedTemplate) Execute(writer io.Writer) (err error) {
	var (
		oopsBuilder = oops.
			Code("template_execute_failed").
			With("writer", writer)
	)

	if err = pmt.instance.Execute(writer, pmt.environment); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while executing template: %s", pmt.instance.Name())
		return
	}

	return
}

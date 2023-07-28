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

type GeneralEnvironment interface {
	IsGeneralEnvironment() bool
}

type PackerMergedEnvironment[PCE GeneralEnvironment] struct {
	PackerGeneralEnvironment *PackerGeneralEnvironment
	PackerCloudEnvironment   PCE
}

type PackerMergedTemplate struct {
	instance    *template.Template
	environment *PackerMergedEnvironment[PackerCloudEnvironment]
}

type PackerCloudTemplate interface {
	GetName() string
	GetInstance() *template.Template
	GetEnvironment() *PackerAWSEnvironment
}

func NewPackerMergedTemplate(generalPackerTemplate *PackerGeneralTemplate, packerCloudTemplate PackerCloudTemplate) (packerMergedTemplate *PackerMergedTemplate, err error) {
	const (
		PACKER_MERGED_TEMPLATE_NAME = "temp_merged_packer_template"
	)

	var (
		oopsBuilder = oops.
				Code("new_packer_merged_template_failed").
				With("generalPackerTemplateName", generalPackerTemplate.GetName()).
				With("packerCloudTemplateName", packerCloudTemplate.GetName())

		packerMergedTemplateInstance      *template.Template
		absPathToTemplatesDir             string
		absPathToTempPackerMergedTemplate string
	)

	if absPathToTemplatesDir, err = filepath.Abs(binaries.TEMPLATE_DIR_NAME); err != nil {
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

	packerMergedTemplate = &PackerMergedTemplate{
		instance: packerMergedTemplateInstance,
		environment: &PackerMergedEnvironment[PackerCloudEnvironment]{
			PackerGeneralEnvironment: generalPackerTemplate.GetEnvironment(),
			PackerCloudEnvironment:   packerCloudTemplate.GetEnvironment(),
		},
	}

	return
}

func (pmt *PackerMergedTemplate) GetName() (name string) {
	return pmt.instance.Name()
}

func (pmt *PackerMergedTemplate) GetInstance() (instance *template.Template) {
	return pmt.instance
}

func (pmt *PackerMergedTemplate) GetEnvironment() (environment *PackerMergedEnvironment[PackerCloudEnvironment]) {
	return pmt.environment
}

func (pmt *PackerMergedTemplate) Remove() (err error) {
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

func (pmt *PackerMergedTemplate) Execute(writer io.Writer) (err error) {
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

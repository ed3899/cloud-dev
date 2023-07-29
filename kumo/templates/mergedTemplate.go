package templates

import (
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
)

type MergedEnvironment[E EnvironmentI] struct {
	general E
	cloud   E
}

type MergedTemplate struct {
	instance    *template.Template
	environment *MergedEnvironment[EnvironmentI]
}

func NewMergedTemplate(generalTemplate, cloudTemplate TemplateSingle) (packerMergedTemplate *MergedTemplate, err error) {
	const (
		TEMPLATE_DIR_NAME    = "templates"
		MERGED_TEMPLATE_NAME = "temp_merged_template"
	)

	var (
		oopsBuilder = oops.
				Code("new_packer_merged_template_failed").
				With("generalTemplate", generalTemplate.GetName()).
				With("cloudTemplate", cloudTemplate.GetName())

		mergedTemplateInstance            *template.Template
		absPathToTemplatesDir             string
		absPathToTempPackerMergedTemplate string
	)

	if generalTemplate.GetParentDirName() != cloudTemplate.GetParentDirName() {
		err = oopsBuilder.
			With("generalTemplate.GetParentDirName()", generalTemplate.GetParentDirName()).
			With("cloudTemplate.GetParentDirName()", cloudTemplate.GetParentDirName()).
			Errorf("generalTemplate and cloudTemplate must be in the same directory")
		return
	}

	if generalTemplate.GetEnvironment().IsNotValidEnvironment() || cloudTemplate.GetEnvironment().IsNotValidEnvironment() {
		err = oopsBuilder.
			Errorf("generalTemplate and cloudTemplate must have valid environments")
		return
	}

	if absPathToTemplatesDir, err = filepath.Abs(TEMPLATE_DIR_NAME); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while crafting absolute path to %s", TEMPLATE_DIR_NAME)
		return
	}

	absPathToTempPackerMergedTemplate = filepath.Join(absPathToTemplatesDir, generalTemplate.GetParentDirName(), MERGED_TEMPLATE_NAME)

	if err = utils.MergeFilesTo(
		absPathToTempPackerMergedTemplate,
		generalTemplate.GetName(),
		cloudTemplate.GetName(),
	); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while merging %s and %s to %s", generalTemplate.GetName(), cloudTemplate.GetName(), absPathToTempPackerMergedTemplate)
		return
	}

	if mergedTemplateInstance, err = template.ParseFiles(absPathToTempPackerMergedTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPathToTempPackerMergedTemplate)
		return
	}

	packerMergedTemplate = &MergedTemplate{
		instance: mergedTemplateInstance,
		environment: &MergedEnvironment[EnvironmentI]{
			general: generalTemplate.GetEnvironment(),
			cloud:   cloudTemplate.GetEnvironment(),
		},
	}

	return
}

func (mt *MergedTemplate) GetName() (name string) {
	return mt.instance.Name()
}

func (mt *MergedTemplate) GetInstance() (instance *template.Template) {
	return mt.instance
}

func (mt *MergedTemplate) GetEnvironment() (environment *MergedEnvironment[EnvironmentI]) {
	return mt.environment
}

func (mt *MergedTemplate) Remove() (err error) {
	var (
		oopsBuilder = oops.
			Code("merged_template_remove_failed")
	)

	if os.RemoveAll(mt.instance.Name()); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", mt.instance.Name())
		return
	}

	return
}

func (mt *MergedTemplate) Execute(writer io.Writer) (err error) {
	var (
		oopsBuilder = oops.
			Code("merged_template_execute_failed").
			With("writer", writer)
	)

	if err = mt.instance.Execute(writer, mt.environment); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while executing template: %s", mt.instance.Name())
		return
	}

	return
}

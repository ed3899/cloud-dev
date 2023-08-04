package merged

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/hashicorp_vars"
	common_templates_interfaces "github.com/ed3899/kumo/common/templates/interfaces"
	common_templates_structs "github.com/ed3899/kumo/common/templates/structs"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
)

type Template struct {
	instance    *template.Template
	absPath     string
	environment *common_templates_structs.Environment[common_templates_interfaces.Environment]
}

func New(generalTemplate, cloudTemplate common_templates_interfaces.Template) (packerMergedTemplate *Template, err error) {
	var (
		oopsBuilder = oops.
				Code("new_packer_merged_template_failed").
				With("generalTemplate", generalTemplate.AbsPath()).
				With("cloudTemplate", cloudTemplate.AbsPath())

		mergedTemplateInstance     *template.Template
		absPathToTemplatesDir      string
		absPathToMergedTemplateDir string
		absPathToMergedTemplate    string
	)

	if generalTemplate.ParentDirName() != cloudTemplate.ParentDirName() {
		err = oopsBuilder.
			With("generalTemplate.GetParentDirName()", generalTemplate.ParentDirName()).
			With("cloudTemplate.GetParentDirName()", cloudTemplate.ParentDirName()).
			Errorf("generalTemplate and cloudTemplate must be in the same directory")
		return
	}

	if generalTemplate.Environment().IsNotValidEnvironment() || cloudTemplate.Environment().IsNotValidEnvironment() {
		err = oopsBuilder.
			Errorf("generalTemplate and cloudTemplate must have valid environments")
		return
	}

	if absPathToTemplatesDir, err = filepath.Abs(dirs.TEMPLATES_DIR_NAME); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while crafting absolute path to %s", dirs.TEMPLATES_DIR_NAME)
		return
	}

	absPathToMergedTemplateDir = filepath.Join(absPathToTemplatesDir, generalTemplate.ParentDirName())

	if absPathToMergedTemplate, err = utils.MergeFilesTo(
		absPathToMergedTemplateDir,
		generalTemplate.AbsPath(),
		cloudTemplate.AbsPath(),
	); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while merging %s and %s to %s", generalTemplate.AbsPath(), cloudTemplate.AbsPath(), absPathToMergedTemplateDir)
		return
	}

	if mergedTemplateInstance, err = template.ParseFiles(absPathToMergedTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPathToMergedTemplate)
		return
	}

	packerMergedTemplate = &Template{
		instance: mergedTemplateInstance,
		absPath:  absPathToMergedTemplate,
		environment: &common_templates_structs.Environment[common_templates_interfaces.Environment]{
			General: generalTemplate.Environment(),
			Cloud:   cloudTemplate.Environment(),
		},
	}

	return
}

func (mt *Template) AbsPath() (path string) {
	return mt.absPath
}

func (mt *Template) Name() (name string) {
	return mt.instance.Name()
}

func (mt *Template) Instance() (instance *template.Template) {
	return mt.instance
}

func (mt *Template) Environment() (
	environment *common_templates_structs.Environment[common_templates_interfaces.Environment],
) {
	return mt.environment
}

func (mt *Template) Remove() (err error) {
	var (
		oopsBuilder = oops.
			Code("merged_template_remove_failed")
	)

	if err = os.Remove(mt.absPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", mt.absPath)
		return
	}

	return
}

func (mt *Template) ExecuteOn(hashicorpVars hashicorp_vars.HashicorpVarsI) (err error) {
	var (
		oopsBuilder = oops.
			Code("merged_template_execute_failed").
			With("hashicorpVars.GetFile().Name()", hashicorpVars.GetFile().Name())
	)

	if err = mt.instance.Execute(hashicorpVars.GetFile(), mt.environment); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while executing template: %s", mt.instance.Name())
		return
	}

	return
}

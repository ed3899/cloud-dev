package interfaces

import (
	"text/template"
)

type TemplateFileRemover interface {
	Remove(
		os_Remove func(string) error,
	) (
		err error,
	)
}

type TemplateFileCreator interface {
	Create(
		file_MergedFilesTo func(
			outputDirAbsPath string,
			inputAbsFilePaths ...string,
		) (
			err error,
		),
		toolConfig ToolConfig,
	) (
		err error,
	)
}

type TemplateFileTextTemplater interface {
	TextTemplate(
		template_ParseFiles func(...string) (*template.Template, error),
	) (
		textTemplate *template.Template,
		err error,
	)
}

type TemplateFile interface {
	TemplateFileCreator
	TemplateFileTextTemplater
	TemplateFileRemover
}

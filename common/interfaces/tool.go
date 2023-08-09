package interfaces

type GeneralTemplateFileGetter interface {
	GeneralTemplateFile() string
}

type CloudTemplateFileGetter interface {
	CloudTemplateFile() string
}

type ToolConfig interface {
	GeneralTemplateFileGetter
	CloudTemplateFileGetter
}

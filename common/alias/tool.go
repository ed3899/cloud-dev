package alias

type ToolName string

func (tn ToolName) String() (toolName string) {
	toolName = string(tn)
	return
}

type ToolVersion string

func (tv ToolVersion) String() (toolVersion string) {
	toolVersion = string(tv)
	return
}

type ToolUrl string

func (tu ToolUrl) String() (toolUrl string) {
	toolUrl = string(tu)
	return
}


package tool

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/common/alias"
	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/config/cloud"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/utils/host"
	"github.com/ed3899/kumo/utils/url"
	"github.com/samber/oops"
)

func NewTool(
	opts ...Option,
) *Tool {
	tool := &Tool{}

	for _, option := range opts {
		option(tool)
	}

	return tool
}

type Tool struct {
	ExecutableName string
	Path 				 string
}

type ToolAbsPath struct {
	Executable   string
	Dir          *ToolDir
	TemplateFile *TemplateFileCombo
}

type DirChangerF func(dir string) error

type Option func(*Tool)

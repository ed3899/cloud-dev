package config

import (
	"fmt"
)

func ToolUrl(
	args *ToolUrlArgs,
) string {
	return fmt.Sprintf(
		"https://releases.hashicorp.com/%s/%s/%s_%s_%s_%s.zip",
		args.ToolName(),
		args.ToolVersion(),
		args.ToolName(),
		args.ToolVersion(),
		args.CurrentOs(),
		args.CurrentArch(),
	)
}

type ToolUrlArgs struct {
	ToolName    func() string
	ToolVersion func() string
	CurrentOs   func() string
	CurrentArch func() string
}

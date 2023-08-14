package url

import "fmt"

type ITool interface {
	Name() string
	Version() string
}

func BuildHashicorpUrl(tool ITool, os, arch string) string {
	return fmt.Sprintf("https://releases.hashicorp.com/%s/%s/%s_%s_%s_%s.zip", tool.Name(), tool.Version(), tool.Name(), tool.Version(), os, arch)
}

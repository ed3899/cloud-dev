package manager

import (
	"os"

	"github.com/ed3899/kumo/common/iota"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PluginsEnvironmentVars", Ordered, func() {
	var (
		envValue = "some-value"

		manager *Manager
	)

	BeforeEach(func() {
		manager = &Manager{
			Tool: iota.Packer,
			Path: &Path{
				Dir: &Dir{
					Plugins: envValue,
				},
			},
		}
	})

	It("should set the plugin path environment variable", Label("unit"), func() {
		err := manager.SetPluginsEnvironmentVars()
		Expect(err).ToNot(HaveOccurred())

		value, ok := os.LookupEnv(manager.Tool.PluginPathEnvironmentVariable())
		Expect(ok).To(BeTrue())
		Expect(value).To(Equal(envValue))
	})

	It("should unset the plugin path environment variable", Label("unit"), func() {
		err := manager.UnsetPluginsEnvironmentVars()
		Expect(err).ToNot(HaveOccurred())

		_, ok := os.LookupEnv(manager.Tool.PluginPathEnvironmentVariable())
		Expect(ok).To(BeFalse())
	})
})

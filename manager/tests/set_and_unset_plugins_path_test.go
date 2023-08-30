package tests

import (
	"os"

	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PluginsEnvironmentVars", Ordered, func() {
	var (
		envValue = "some-value"

		_manager *manager.Manager
	)

	BeforeEach(func() {
		_manager = &manager.Manager{
			Tool: iota.Packer,
			Path: &manager.Path{
				Dir: &manager.Dir{
					Plugins: envValue,
				},
			},
		}
	})

	It("should set the plugin path environment variable", Label("unit"), func() {
		err := _manager.SetPluginsPath()
		Expect(err).ToNot(HaveOccurred())

		value, ok := os.LookupEnv(_manager.Tool.PluginPathEnvironmentVariable())
		Expect(ok).To(BeTrue())
		Expect(value).To(Equal(envValue))
	})

	It("should unset the plugin path environment variable", Label("unit"), func() {
		err := _manager.UnsetPluginsEnvironmentVars()
		Expect(err).ToNot(HaveOccurred())

		_, ok := os.LookupEnv(_manager.Tool.PluginPathEnvironmentVariable())
		Expect(ok).To(BeFalse())
	})
})

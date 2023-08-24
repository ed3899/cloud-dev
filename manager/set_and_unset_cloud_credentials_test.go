package manager

import (
	"os"

	"github.com/ed3899/kumo/common/iota"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("SetAndUnsetCloudCredentials", Ordered, func() {
	var (
		mockAccessKeyId     = "some-access-key-id"
		mockSecretAccessKey = "some-secret-access-key"

		manager *Manager
	)

	Context("when the cloud is AWS", func() {
		BeforeEach(func() {
			manager = &Manager{
				Cloud: iota.Aws,
			}

			viper.Set("AWS.AccessKeyId", mockAccessKeyId)
			viper.Set("AWS.SecretAccessKey", mockSecretAccessKey)
		})

		AfterEach(func() {
			viper.Reset()
		})

		It("sets the the right environment variables", func() {
			err := manager.SetCloudCredentials()
			Expect(err).NotTo(HaveOccurred())

			value1, ok := os.LookupEnv("AWS_ACCESS_KEY_ID")
			Expect(ok).To(BeTrue())
			Expect(value1).To(Equal(mockAccessKeyId))

			value2, ok := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
			Expect(ok).To(BeTrue())
			Expect(value2).To(Equal(mockSecretAccessKey))
		})

		It("unsets the right environment variables", func() {
			err := manager.UnsetCloudCredentials()
			Expect(err).NotTo(HaveOccurred())

			_, ok := os.LookupEnv("AWS_ACCESS_KEY_ID")
			Expect(ok).To(BeFalse())

			_, ok = os.LookupEnv("AWS_SECRET_ACCESS_KEY")
			Expect(ok).To(BeFalse())
		})
	})

})

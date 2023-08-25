package tests

import (
	"os"

	"github.com/ed3899/kumo/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoToDirRunAndInitial", Ordered, func() {
	Context("with a valid dir path", Label("unit"), func() {
		var (
			initialDirPath string
			tempDirPath    string
			_manager       *manager.Manager
			err            error
		)

		BeforeEach(func() {
			initialDirPath, err = os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			tempDirPath, err = os.MkdirTemp(initialDirPath, "test")
			Expect(err).ToNot(HaveOccurred())

			_manager = &manager.Manager{
				Path: &manager.Path{
					Dir: &manager.Dir{
						Initial: initialDirPath,
						Run:     tempDirPath,
					},
				},
			}
		})

		AfterEach(func() {
			Expect(os.RemoveAll(tempDirPath)).To(Succeed())
		})

		It("should complete a round trip", func() {
			By("changing to the run dir", func() {
				Expect(_manager.GoToDirRun()).To(Succeed())
				Expect(os.Getwd()).To(Equal(tempDirPath))
			})

			By("changing to the initial dir", func() {
				Expect(_manager.GoToDirInitial()).To(Succeed())
				Expect(os.Getwd()).To(Equal(initialDirPath))
			})
		})
	})

	Context("with an invalid dir path", Label("unit"), func() {
		var (
			_manager *manager.Manager
		)

		BeforeEach(func() {
			_manager = &manager.Manager{
				Path: &manager.Path{
					Dir: &manager.Dir{
						Initial: "",
						Run:     "",
					},
				},
			}
		})

		It("should return an error", func() {
			By("changing to the run dir", func() {
				Expect(_manager.GoToDirRun()).ToNot(Succeed())
			})

			By("changing to the initial dir", func() {
				Expect(_manager.GoToDirInitial()).ToNot(Succeed())
			})
		})
	})
})

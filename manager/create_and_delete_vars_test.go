package manager

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateAndDeleteVars", Ordered, func() {
	var (
		manager *Manager
	)

	Context("with a valid path", func() {
		BeforeEach(func() {
			cwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			manager = &Manager{
				Path: &Path{
					Vars: filepath.Join(cwd, "vars.yml"),
				},
			}
		})

		It("creates a vars file", func() {
			file, err := manager.CreateVars()
			defer file.Close()
			Expect(err).NotTo(HaveOccurred())
			Expect(file).NotTo(BeNil())
		})

		It("deletes a vars file", func() {
			err := manager.DeleteVars()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("with an invalid path", func() {
		BeforeEach(func() {
			manager = &Manager{
				Path: &Path{
					Vars: "",
				},
			}
		})

		It("fails to create a vars file", func() {
			file, err := manager.CreateVars()
			Expect(err).To(HaveOccurred())
			Expect(file).To(BeNil())
		})
	})
})

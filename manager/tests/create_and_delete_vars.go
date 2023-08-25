package tests

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateAndDeleteVars", Ordered, func() {
	var (
		_manager *manager.Manager
	)

	Context("with a valid path", func() {
		BeforeEach(func() {
			cwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			_manager = &manager.Manager{
				Path: &manager.Path{
					Vars: filepath.Join(cwd, "vars.yml"),
				},
			}
		})

		It("creates a vars file", func() {
			file, err := _manager.CreateVars()
			defer file.Close()
			Expect(err).NotTo(HaveOccurred())
			Expect(file).NotTo(BeNil())
		})

		It("deletes a vars file", func() {
			err := _manager.DeleteVars()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("with an invalid path", func() {
		BeforeEach(func() {
			_manager = &manager.Manager{
				Path: &manager.Path{
					Vars: "",
				},
			}
		})

		It("fails to create a vars file", func() {
			file, err := _manager.CreateVars()
			Expect(err).To(HaveOccurred())
			Expect(file).To(BeNil())
		})
	})
})

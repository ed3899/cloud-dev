package tests

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateAndDeleteSSHConfig", Ordered, func() {
	Context("with a valid ip file", func() {
		var (
			ipFilePath string
		)

		BeforeEach(func() {
			cwd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			ipFilePath = filepath.Join(cwd, "ipfile")

			err = os.WriteFile(ipFilePath, []byte("127.0.0.1"), 0644)
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err := os.Remove(ipFilePath)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("with a valid ssh config path", Label("unit"), func() {
			var (
				_manager      *manager.Manager
				sshConfigPath string
			)

			BeforeEach(func() {
				cwd, err := os.Getwd()
				Expect(err).ToNot(HaveOccurred())

				sshConfigPath = filepath.Join(cwd, "sshconfig")

				_manager = &manager.Manager{
					Path: &manager.Path{
						Terraform: &manager.Terraform{
							IpFile:    ipFilePath,
							SshConfig: sshConfigPath,
						},
					},
				}
			})

			It("should generate a ssh config file", func() {
				err := _manager.CreateSshConfig()
				Expect(err).ToNot(HaveOccurred())

				_, err = os.Stat(sshConfigPath)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should delete the ssh config file", func() {
				err := _manager.DeleteSshConfig()
				Expect(err).ToNot(HaveOccurred())

				_, err = os.Stat(sshConfigPath)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with an invalid ssh config path", Label("unit"), func() {
			var (
				_manager *manager.Manager
			)

			BeforeEach(func() {
				_manager = &manager.Manager{
					Path: &manager.Path{
						Terraform: &manager.Terraform{
							IpFile:    ipFilePath,
							SshConfig: "",
						},
					},
				}
			})

			It("should return an error", func() {
				err := _manager.CreateSshConfig()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("with an invalid ip file", func() {
		Context("with a valid ssh config path", Label("unit"), func() {
			var (
				_manager      *manager.Manager
				sshConfigPath string
			)

			BeforeEach(func() {
				cwd, err := os.Getwd()
				Expect(err).ToNot(HaveOccurred())

				sshConfigPath = filepath.Join(cwd, "sshconfig")

				_manager = &manager.Manager{
					Path: &manager.Path{
						Terraform: &manager.Terraform{
							IpFile:    "invalid-ip-file",
							SshConfig: sshConfigPath,
						},
					},
				}
			})

			It("should return an error when generating a sshconfig", func() {
				err := _manager.CreateSshConfig()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with an invalid ssh config path", Label("unit"), func() {
			var (
				_manager *manager.Manager
			)

			BeforeEach(func() {
				_manager = &manager.Manager{
					Path: &manager.Path{
						Terraform: &manager.Terraform{
							IpFile:    "invalid-ip-file",
							SshConfig: "",
						},
					},
				}
			})

			It("should return an error when generating a sshconfig", func() {
				err := _manager.CreateSshConfig()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

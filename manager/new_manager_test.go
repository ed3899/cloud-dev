package manager

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/common/iota"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manager", func() {
	var (
		manager                      *Manager
		err                          error
		pathExeSubstring             string
		pathTemplateMergedSubstring  string
		pathTemplateCloudSubstring   string
		pathTemplateBaseSubstring    string
		pathVarsSubtring             string
		pathTerraformLockSubstring   string
		pathTerraformStateSubstring  string
		pathTerraformBackupSubstring string
		pathTerraformIpFile          string
		pathTerraformIdentityFile    string
		pathSshConfigSubstring       string
		pathDirPlugins               string
		pathDirRun                   string
	)

	BeforeEach(func() {
		manager, err = NewManager(iota.Aws, iota.Packer)
		Expect(err).ToNot(HaveOccurred())
		Expect(manager).ToNot(BeNil())

		pathExeSubstring = filepath.Join(
			iota.Dependencies.Name(),
			iota.Packer.Name(),
			fmt.Sprintf("%s.exe", iota.Packer.Name()),
		)

		templateSubstring := func(templateName string) string {
			return filepath.Join(
				iota.Templates.Name(),
				iota.Packer.Name(),
				templateName,
			)
		}
		pathTemplateMergedSubstring = templateSubstring(constants.MERGED_TEMPLATE_NAME)
		pathTemplateCloudSubstring = templateSubstring(iota.Aws.TemplateFiles().Cloud)
		pathTemplateBaseSubstring = templateSubstring(iota.Aws.TemplateFiles().Base)
		pathVarsSubtring = filepath.Join(
			iota.Packer.Name(),
			iota.Aws.Name(),
			iota.Packer.VarsName(),
		)

		terraformPathSubstring := func(fileName string) string {
			return filepath.Join(
				iota.Terraform.Name(),
				iota.Aws.Name(),
				fileName,
			)
		}

		pathTerraformLockSubstring = terraformPathSubstring(constants.TERRAFORM_LOCK)
		pathTerraformStateSubstring = terraformPathSubstring(constants.TERRAFORM_STATE)
		pathTerraformBackupSubstring = terraformPathSubstring(constants.TERRAFORM_BACKUP)
		pathTerraformIpFile = terraformPathSubstring(constants.IP_FILE_NAME)
		pathTerraformIdentityFile = terraformPathSubstring(constants.KEY_NAME)
		pathSshConfigSubstring = constants.CONFIG_NAME

		pathDirPlugins = filepath.Join(
			iota.Packer.Name(),
			iota.Aws.Name(),
			iota.Packer.PluginDir(),
		)
		pathDirRun = filepath.Join(
			iota.Packer.Name(),
			iota.Aws.Name(),
		)

	})

	It("should create a new manager instance", Label("unit"), func() {
		Expect(manager.Cloud).To(Equal(iota.Aws))
		Expect(manager.Tool).To(Equal(iota.Packer))
		Expect(manager.Path.Executable).To(ContainSubstring(pathExeSubstring))
		Expect(manager.Path.Template.Merged).To(ContainSubstring(pathTemplateMergedSubstring))
		Expect(manager.Path.Template.Cloud).To(ContainSubstring(pathTemplateCloudSubstring))
		Expect(manager.Path.Template.Base).To(ContainSubstring(pathTemplateBaseSubstring))
		Expect(manager.Path.Vars).To(ContainSubstring(pathVarsSubtring))
		Expect(manager.Path.Terraform.Lock).To(ContainSubstring(pathTerraformLockSubstring))
		Expect(manager.Path.Terraform.State).To(ContainSubstring(pathTerraformStateSubstring))
		Expect(manager.Path.Terraform.Backup).To(ContainSubstring(pathTerraformBackupSubstring))
		Expect(manager.Path.Terraform.IpFile).To(ContainSubstring(pathTerraformIpFile))
		Expect(manager.Path.Terraform.IdentityFile).To(ContainSubstring(pathTerraformIdentityFile))
		Expect(manager.Path.Terraform.SshConfig).To(ContainSubstring(pathSshConfigSubstring))
		Expect(manager.Path.Dir.Plugins).To(ContainSubstring(pathDirPlugins))
		Expect(manager.Path.Dir.Initial).ToNot(BeNil())
		Expect(manager.Path.Dir.Run).To(ContainSubstring(pathDirRun))
		Expect(manager.Environment).ToNot(BeNil())
	})
})

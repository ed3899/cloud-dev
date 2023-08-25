package tests

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/manager"
	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/common/iota"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manager", func() {
	var (
		_manager                      *manager.Manager
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
		_manager, err = manager.NewManager(iota.Aws, iota.Packer)
		Expect(err).ToNot(HaveOccurred())
		Expect(_manager).ToNot(BeNil())

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
		Expect(_manager.Cloud).To(Equal(iota.Aws))
		Expect(_manager.Tool).To(Equal(iota.Packer))
		Expect(_manager.Path.Executable).To(ContainSubstring(pathExeSubstring))
		Expect(_manager.Path.Template.Merged).To(ContainSubstring(pathTemplateMergedSubstring))
		Expect(_manager.Path.Template.Cloud).To(ContainSubstring(pathTemplateCloudSubstring))
		Expect(_manager.Path.Template.Base).To(ContainSubstring(pathTemplateBaseSubstring))
		Expect(_manager.Path.Vars).To(ContainSubstring(pathVarsSubtring))
		Expect(_manager.Path.Terraform.Lock).To(ContainSubstring(pathTerraformLockSubstring))
		Expect(_manager.Path.Terraform.State).To(ContainSubstring(pathTerraformStateSubstring))
		Expect(_manager.Path.Terraform.Backup).To(ContainSubstring(pathTerraformBackupSubstring))
		Expect(_manager.Path.Terraform.IpFile).To(ContainSubstring(pathTerraformIpFile))
		Expect(_manager.Path.Terraform.IdentityFile).To(ContainSubstring(pathTerraformIdentityFile))
		Expect(_manager.Path.Terraform.SshConfig).To(ContainSubstring(pathSshConfigSubstring))
		Expect(_manager.Path.Dir.Plugins).To(ContainSubstring(pathDirPlugins))
		Expect(_manager.Path.Dir.Initial).ToNot(BeNil())
		Expect(_manager.Path.Dir.Run).To(ContainSubstring(pathDirRun))
		Expect(_manager.Environment).ToNot(BeNil())
	})
})

package templates

import (
	"os"
	"testing"

	templates_terraform_aws "github.com/ed3899/kumo/templates/terraform/aws"
	templates_terraform_generic "github.com/ed3899/kumo/templates/terraform/generic"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTemplates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Templates Suite")
}

var _ = Describe("CraftGenericCloudTerraformTfVarsFile", func() {

	var (
		resultingTerraformTfVarsPath string
		err                          error
		awsEnv                       = &templates_terraform_aws.AWS_TerraformEnvironment{
			AWS_REGION:                   "us-east-1",
			AWS_INSTANCE_TYPE:            "t2.micro",
			AWS_EC2_INSTANCE_VOLUME_TYPE: "gp2",
			AWS_EC2_INSTANCE_VOLUME_SIZE: 8,
		}
	)

	When("the template file exists", func() {
		const (
			cloud                 = "aws"
			existentTemplateName  = "AWS_TerraformTfVars.tmpl"
			terraformVarsFileName = "aws.auto.tfvars"
		)

		AfterEach(func() {
			os.Remove(resultingTerraformTfVarsPath)
		})

		It("should create a file with the specified name", func() {
			resultingTerraformTfVarsPath, err = templates_terraform_generic.CraftGenericCloudTerraformTfVarsFile(cloud, existentTemplateName, terraformVarsFileName, *awsEnv)
			Expect(err).NotTo(HaveOccurred())
			Expect(resultingTerraformTfVarsPath).To(Equal("terraform/aws/example_vars.tfvars"))
		})
	})
})

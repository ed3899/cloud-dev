package templates

import (
	"os"
	"path/filepath"
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
		initialLocation              string
		runLocation                  string
		expectedPath                 string
		awsEnv = &templates_terraform_aws.AWS_TerraformEnvironment{
			AWS_REGION:                   "us-east-1",
			AWS_INSTANCE_TYPE:            "t2.micro",
			AWS_EC2_INSTANCE_VOLUME_TYPE: "gp2",
			AWS_EC2_INSTANCE_VOLUME_SIZE: 8,
		}
	)
	const (
		cloud                   = "aws"
		existentTemplateName    = "AWS_TerraformTfVars.tmpl"
		nonExistentTemplateName = "AWS_TerraformTfVars_non_existent.tmpl"
		terraformVarsFileName   = "aws.auto.tfvars"
	)

	// Change to root directory
	BeforeEach(func() {
		initialLocation, err = os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		runLocation = "../"
		err = os.Chdir(runLocation)
		Expect(err).NotTo(HaveOccurred())

		cwd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		expectedPath = filepath.Join(cwd, "terraform", "aws", "aws.auto.tfvars")
	})

	// Change back to initial directory
	AfterEach(func() {
		err = os.Chdir(initialLocation)
		Expect(err).NotTo(HaveOccurred())
	})

	When("the template file exists", func() {
		It("should create a file with the specified name and return the path", func() {
			resultingTerraformTfVarsPath, err = templates_terraform_generic.CraftGenericCloudTerraformTfVarsFile(cloud, existentTemplateName, terraformVarsFileName, *awsEnv)
			Expect(err).NotTo(HaveOccurred())
			// Remove file after test
			defer func() {
				err = os.Remove(resultingTerraformTfVarsPath)
				Expect(err).NotTo(HaveOccurred())
			}()
			Expect(resultingTerraformTfVarsPath).To(Equal(expectedPath))
			_, err = os.Stat(resultingTerraformTfVarsPath)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("the template file does not exist", func() {
		It("should return an error", func() {
			resultingTerraformTfVarsPath, err = templates_terraform_generic.CraftGenericCloudTerraformTfVarsFile(cloud, nonExistentTemplateName, terraformVarsFileName, *awsEnv)
			Expect(err).To(HaveOccurred())
			Expect(resultingTerraformTfVarsPath).To(Equal(""))
		})
	})
})


var _ = Describe("GetAmiToBeUsed", func ()  {
	
})
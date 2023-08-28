package tests

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ed3899/kumo/manager/environment"
	"github.com/ed3899/kumo/utils/packer_manifest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("NewTerraformAwsEnvironment", func() {
	var (
		amiId = "ami:ami-12345678"
		tempManifestFilePath string
	)

	BeforeEach(func() {
		// Create a temporary manifest file for testing
		tempManifest := &packer_manifest.PackerManifest{
			Builds: []*packer_manifest.PackerBuild{
				{PackerRunUUID: "run_uuid_1", ArtifactId: fmt.Sprintf("ami:%s", amiId)},
				{PackerRunUUID: "run_uuid_2", ArtifactId: "ami:ami-98765432"},
			},
			LastRunUUID: "run_uuid_1",
		}

		cwd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		tempFile, err := os.CreateTemp(cwd, "packer_manifest.json")
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		jsonEncoder := json.NewEncoder(tempFile)
		Expect(jsonEncoder.Encode(tempManifest)).To(Succeed())

		tempManifestFilePath = tempFile.Name()
	})

	AfterEach(func() {
		// Clean up temporary files
		Expect(os.Remove(tempManifestFilePath)).To(Succeed())
	})

	var (
		AccessKeyId     = "fake-access-key-id"
		SecretAccessKey = "secret-access-key"
		Region          = "us-east-1"
		InstanceType    = "t2.micro"
		User            = "ubuntu"
		VolumeType      = "gp2"
		VolumeSize      = 8
	)

	JustBeforeEach(func() {
		viper.Set("AWS.AccessKeyId", AccessKeyId)
		viper.Set("AWS.SecretAccessKey", SecretAccessKey)
		viper.Set("AWS.Region", Region)
		viper.Set("AWS.EC2.Instance.Type", InstanceType)
		viper.Set("AMI.User", User)
		viper.Set("AWS.EC2.Volume.Type", VolumeType)
		viper.Set("AWS.EC2.Volume.Size", VolumeSize)
	})

	JustAfterEach(func() {
		viper.Reset()
	})

	It("should return a new TerraformAwsEnvironment", func() {
		_environment, err := environment.NewTerraformAwsEnvironment(tempManifestFilePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(_environment).NotTo(BeNil())
	})
})

package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/utils/packer_manifest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetLastAmiIdFromPackerManifest", func() {
	var (
		amiId = "ami:ami-12345678"

		tempManifestFilePath string
		fakeManifestFilePath string
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
		fakeManifestFilePath = filepath.Join(cwd, "packer_manifest_fake.json")
	})

	AfterEach(func() {
		// Clean up temporary files
		Expect(os.Remove(tempManifestFilePath)).To(Succeed())
	})

	Context("when the manifest file is absolute", func() {
		Context("when the manifest file exists", func() {
			It("should return the last built AMI ID from the manifest file", Label("unit"), func() {
				amiId, err := packer_manifest.GetLastBuiltAmiIdFromPackerManifest(tempManifestFilePath)
				Expect(err).NotTo(HaveOccurred())
				Expect(amiId).To(Equal(amiId))
			})
		})

		Context("when the manifest file doesn't exists", func() {
			It("should return an error for non-existent manifest file", Label("unit"), func() {
				amiId, err := packer_manifest.GetLastBuiltAmiIdFromPackerManifest(fakeManifestFilePath)
				Expect(err).To(HaveOccurred())
				Expect(amiId).To(BeEmpty())
				Expect(err).To(MatchError(ContainSubstring("Error occurred while opening packer manifest file")))
			})
		})
	})

	Context("when the manifest file is not absolute", func() {
		Context("when the manifest file exists", func() {
			It("should return an error for invalid manifest path", Label("unit"), func() {
				amiId, err := packer_manifest.GetLastBuiltAmiIdFromPackerManifest("./packer_manifest.json")
				Expect(err).To(HaveOccurred())
				Expect(amiId).To(BeEmpty())
				Expect(err).To(MatchError(ContainSubstring("path is not absolute")))
			})
		})

		Context("when the manifest file doesn't exists", func() {
			It("should return an error for invalid manifest path", Label("unit"), func() {
				amiId, err := packer_manifest.GetLastBuiltAmiIdFromPackerManifest("invalid_manifest_path.json")
				Expect(err).To(HaveOccurred())
				Expect(amiId).To(BeEmpty())
				Expect(err.Error()).To(ContainSubstring("path is not absolute"))
			})
		})
	})
})

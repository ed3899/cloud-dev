package packer_manifest

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PickAmiId", func() {
	Context("when lastBuildAmiId is empty", func() {
		Context("when amiIdFromConfig is empty", func() {
			It("returns an error", Label("unit"), func() {
				_, err := PickAmiId("", "")

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when amiIdFromConfig is not empty", func() {
			It("returns an error", Label("unit"), func() {
				_, err := PickAmiId("", "ami-1234567890")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("when amiIdFromConfig is not empty", func() {
		lastBuiltAmiId := "ami-1234567890"

		Context("when amiIdFromConfig is empty", func() {
			It("returns the last built ami id", Label("unit"), func() {
				amiIdToBeUsed, err := PickAmiId(lastBuiltAmiId, "")

				Expect(err).NotTo(HaveOccurred())
				Expect(amiIdToBeUsed).To(Equal(lastBuiltAmiId))
			})
		})

		Context("when amiIdFromConfig is not empty", func() {
			amiIdFromConfig := "ami-0987654321"

			It("returns the amiIdFromConfig", Label("unit"), func() {
				amiIdToBeUsed, err := PickAmiId(lastBuiltAmiId, amiIdFromConfig)

				Expect(err).NotTo(HaveOccurred())
				Expect(amiIdToBeUsed).To(Equal(amiIdFromConfig))
			})
		})
	})
})

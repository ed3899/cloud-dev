package tests

import (
	"github.com/ed3899/kumo/manager/environment"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("NewPackerBaseEnvironment", func() {
	BeforeEach(func() {
		viper.Set("Git.Username", "git_username")
		viper.Set("Git.Email", "git_email")
		viper.Set("AMI.Tools", []string{"ansible_tag_1", "ansible_tag_2"})
		viper.Set("GitHub.PersonalAccessTokenClassic", "git_hub_personal_access_token_classic")
	})

	AfterEach(func() {
		viper.Reset()
	})

	It("should return a new PackerBaseEnvironment", func() {
		_environment := environment.NewPackerBaseEnvironment()
		Expect(_environment).ToNot(BeNil())
		Expect(_environment.Required.GIT_USERNAME).To(Equal("git_username"))
		Expect(_environment.Required.GIT_EMAIL).To(Equal("git_email"))
		Expect(_environment.Required.ANSIBLE_TAGS).To(Equal([]string{"ansible_tag_1", "ansible_tag_2"}))
		Expect(_environment.Optional.GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC).To(Equal("git_hub_personal_access_token_classic"))
	})
})

package tests

import (
	"github.com/ed3899/kumo/manager/environment"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("NewPackerAwsEnvironment", func() {
	BeforeEach(func() {
		viper.Set("AWS.AccessKeyId", "access_key_id")
		viper.Set("AWS.SecretAccessKey", "secret_access_key")
		viper.Set("AWS.IamProfile", "iam_profile")
		viper.Set("AWS.UserIds", []string{"user_id_1", "user_id_2"})
		viper.Set("AMI.Name", "ami_name")
		viper.Set("AWS.EC2.Instance.Type", "instance_type")
		viper.Set("AWS.Region", "region")
		viper.Set("AMI.Base.Filter", "ami_base_filter")
		viper.Set("AMI.Base.RootDeviceType", "ami_base_root_device_type")
		viper.Set("AMI.Base.VirtualizationType", "ami_base_virtualization_type")
		viper.Set("AMI.Base.Owners", []string{"ami_base_owner_1", "ami_base_owner_2"})
		viper.Set("AMI.Base.User", "ami_base_user")
		viper.Set("AMI.User", "ami_user")
		viper.Set("AMI.Home", "ami_home")
		viper.Set("AMI.Password", "ami_password")
	})

	AfterEach(func() {
		viper.Reset()
	})

	It("should return a new PackerAwsEnvironment", func() {
		_environment := environment.NewPackerAwsEnvironment()
		Expect(_environment).ToNot(BeNil())
		Expect(_environment.Required.AWS_ACCESS_KEY).To(Equal("access_key_id"))
		Expect(_environment.Required.AWS_SECRET_KEY).To(Equal("secret_access_key"))
		Expect(_environment.Required.AWS_IAM_PROFILE).To(Equal("iam_profile"))
		Expect(_environment.Required.AWS_USER_IDS).To(Equal([]string{"user_id_1", "user_id_2"}))
		Expect(_environment.Required.AWS_AMI_NAME).To(Equal("ami_name"))
		Expect(_environment.Required.AWS_INSTANCE_TYPE).To(Equal("instance_type"))
		Expect(_environment.Required.AWS_REGION).To(Equal("region"))
		Expect(_environment.Required.AWS_EC2_AMI_NAME_FILTER).To(Equal("ami_base_filter"))
		Expect(_environment.Required.AWS_EC2_AMI_ROOT_DEVICE_TYPE).To(Equal("ami_base_root_device_type"))
		Expect(_environment.Required.AWS_EC2_AMI_VIRTUALIZATION_TYPE).To(Equal("ami_base_virtualization_type"))
		Expect(_environment.Required.AWS_EC2_AMI_OWNERS).To(Equal([]string{"ami_base_owner_1", "ami_base_owner_2"}))
		Expect(_environment.Required.AWS_EC2_SSH_USERNAME).To(Equal("ami_base_user"))
		Expect(_environment.Required.AWS_EC2_INSTANCE_USERNAME).To(Equal("ami_user"))
		Expect(_environment.Required.AWS_EC2_INSTANCE_USERNAME_HOME).To(Equal("ami_home"))
		Expect(_environment.Required.AWS_EC2_INSTANCE_USERNAME_PASSWORD).To(Equal("ami_password"))
	})
})

package tests

import (
	"github.com/ed3899/kumo/manager/environment"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("NewPackerAwsEnvironment", func() {
	var (
		AccessKeyId               = "access_key_id"
		SecretAccessKey           = "secret_access_key"
		IamProfile                = "iam_profile"
		UserIds                   = []string{"user_id_1", "user_id_2"}
		AmiName                   = "ami_name"
		InstanceType              = "instance_type"
		Region                    = "region"
		AmiBaseFilter             = "ami_base_filter"
		AmiBaseRootDeviceType     = "ami_base_root_device_type"
		AmiBaseVirtualizationType = "ami_base_virtualization_type"
		AmiBaseOwners             = []string{"ami_base_owner_1", "ami_base_owner_2"}
		AmiBaseUser               = "ami_base_user"
		AmiUser                   = "ami_user"
		AmiHome                   = "ami_home"
		AmiPassword               = "ami_password"
	)

	BeforeEach(func() {
		viper.Set("AWS.AccessKeyId", AccessKeyId)
		viper.Set("AWS.SecretAccessKey", SecretAccessKey)
		viper.Set("AWS.IamProfile", IamProfile)
		viper.Set("AWS.UserIds", UserIds)
		viper.Set("AMI.Name", AmiName)
		viper.Set("AWS.EC2.Instance.Type", InstanceType)
		viper.Set("AWS.Region", Region)
		viper.Set("AMI.Base.Filter", AmiBaseFilter)
		viper.Set("AMI.Base.RootDeviceType", AmiBaseRootDeviceType)
		viper.Set("AMI.Base.VirtualizationType", AmiBaseVirtualizationType)
		viper.Set("AMI.Base.Owners", AmiBaseOwners)
		viper.Set("AMI.Base.User", AmiBaseUser)
		viper.Set("AMI.User", AmiUser)
		viper.Set("AMI.Home", AmiHome)
		viper.Set("AMI.Password", AmiPassword)
	})

	AfterEach(func() {
		viper.Reset()
	})

	It("should return a new PackerAwsEnvironment", func() {
		_environment := environment.NewPackerAwsEnvironment()
		Expect(_environment).ToNot(BeNil())
		Expect(_environment.Required.AWS_ACCESS_KEY).To(Equal(AccessKeyId))
		Expect(_environment.Required.AWS_SECRET_KEY).To(Equal(SecretAccessKey))
		Expect(_environment.Required.AWS_IAM_PROFILE).To(Equal(IamProfile))
		Expect(_environment.Required.AWS_USER_IDS).To(Equal(UserIds))
		Expect(_environment.Required.AWS_AMI_NAME).To(Equal(AmiName))
		Expect(_environment.Required.AWS_INSTANCE_TYPE).To(Equal(InstanceType))
		Expect(_environment.Required.AWS_REGION).To(Equal(Region))
		Expect(_environment.Required.AWS_EC2_AMI_NAME_FILTER).To(Equal(AmiBaseFilter))
		Expect(_environment.Required.AWS_EC2_AMI_ROOT_DEVICE_TYPE).To(Equal(AmiBaseRootDeviceType))
		Expect(_environment.Required.AWS_EC2_AMI_VIRTUALIZATION_TYPE).To(Equal(AmiBaseVirtualizationType))
		Expect(_environment.Required.AWS_EC2_AMI_OWNERS).To(Equal(AmiBaseOwners))
		Expect(_environment.Required.AWS_EC2_SSH_USERNAME).To(Equal(AmiBaseUser))
		Expect(_environment.Required.AWS_EC2_INSTANCE_USERNAME).To(Equal(AmiUser))
		Expect(_environment.Required.AWS_EC2_INSTANCE_USERNAME_HOME).To(Equal(AmiHome))
		Expect(_environment.Required.AWS_EC2_INSTANCE_USERNAME_PASSWORD).To(Equal(AmiPassword))
	})
})

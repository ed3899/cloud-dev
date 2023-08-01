package cloud

type CloudSetupI interface {
	GetCloudName() (string)
	GetCloudType() (CloudType)
}
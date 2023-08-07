package interfaces

type PackerCloudEnvironmentI interface {
	IsPackerCloudEnvironment() bool
}

type TerraformCloudEnvironmentI interface {
	IsTerraformCloudEnvironment() bool
}

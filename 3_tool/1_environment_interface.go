package tool

type EnvironmentI interface {
	IsEnvironment() bool
}

type GeneralEnvironmentI interface {
	IsGeneralEnvironment() bool
}

type CloudEnvironmentI interface {
	IsCloudEnvironment() bool
}
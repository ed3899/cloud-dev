package environment

type GeneralEnvironmentI interface {
	IsGeneralEnvironment() bool
}

type CloudEnvironmentI interface {
	IsCloudEnvironment() bool
}


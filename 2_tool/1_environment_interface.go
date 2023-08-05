package tool

type EnvironmentI interface {
	IsNotValidEnvironment() bool
}
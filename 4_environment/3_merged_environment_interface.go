package environment

type MergedEnvironment[E EnvironmentI] struct {
	General E
	Cloud   E
}

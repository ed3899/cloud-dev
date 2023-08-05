package tool

type Merged[E EnvironmentI] struct {
	General E
	Cloud   E
}

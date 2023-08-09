package interfaces

type VersionGetter[V ~string] interface {
	Version() V
}
package interfaces

type NameGetter[N ~string] interface {
	Name() N
}

package interfaces

type KindGetter[K ~int] interface {
	Kind() K
}
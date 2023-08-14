package interfaces

type IClone[T any] interface {
	Clone() T
}
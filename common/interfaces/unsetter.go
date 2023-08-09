package interfaces

type Unsetter interface {
	Unset() (err error)
}

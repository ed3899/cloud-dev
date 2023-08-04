package interfaces

type Credentials interface {
	Set() (err error)
	Unset() (err error)
}
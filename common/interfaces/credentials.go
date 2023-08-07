package interfaces

type Setter interface {
	Set() (err error)
}

type Unsetter interface {
	Unset() (err error)
}

type Credentials interface {
	Setter
	Unsetter
}

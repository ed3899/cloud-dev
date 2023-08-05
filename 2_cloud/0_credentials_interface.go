package cloud

type CredentialsI interface {
	Set() (err error)
	Unset() (err error)
}

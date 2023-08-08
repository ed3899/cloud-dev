package interfaces

type SshConfigI interface {
	Create() (error)
	Remove() (error)
	GetAbsPath() (string)
}
package ssh

type SshConfigI interface {
	Create() (error)
	Remove() (error)
	GetAbsPath() (string)
}
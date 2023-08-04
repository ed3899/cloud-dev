package interfaces

type Template interface {
	ParentDirName() string
	AbsPath() string
	Environment() Environment
}
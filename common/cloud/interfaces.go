package cloud

type ConfigI interface {
	GetName() string
	GetType() Type
}

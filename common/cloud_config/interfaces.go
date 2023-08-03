package cloud_config

type CloudI interface {
	GetName() string
	GetType() Type
}

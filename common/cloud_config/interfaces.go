package cloud_config

type CloudI interface {
	Name() (cloudName string)
	Type() (cloudType Kind)
}

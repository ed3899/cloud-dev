package tool

type ToolSetupI interface {
	GetToolType() ToolType
	GetInitialDir() string
	GoInitialDir() error
	GoTargetDir() error
}

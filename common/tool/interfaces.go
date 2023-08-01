package tool

type ToolSetupI interface {
	GetToolType() ToolType
	GoInitialDir() error
	GoTargetDir() error
}

package block

type CameraBlock struct {
	SolidBase
}

func NewCameraBlock() *CameraBlock {
	return &CameraBlock{
		SolidBase: SolidBase{
			BlockID:       CAMERA,
			BlockName:     "Camera",
			BlockHardness: -1,
			BlockToolType: ToolTypeNone,
		},
	}
}

func (b *CameraBlock) GetDrops(toolType, toolTier int) []Drop {
	return nil
}

func init() {
	Registry.Register(NewCameraBlock())
}

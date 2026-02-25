package block

// CameraBlock 相机方块 (MCPE 独有教育版方块)
// MCPE 方块 ID 242
// 大部分服务器不使用此方块，但需要注册以避免未知方块

type CameraBlock struct {
	SolidBase
}

func NewCameraBlock() *CameraBlock {
	return &CameraBlock{
		SolidBase: SolidBase{
			BlockID:       CAMERA,
			BlockName:     "Camera",
			BlockHardness: -1, // 不可破坏
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

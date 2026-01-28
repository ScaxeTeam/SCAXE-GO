package adapter

type BlockState struct {
	ID   uint8
	Meta uint8
}

func NewBlockState(id uint8, meta uint8) BlockState {
	return BlockState{ID: id, Meta: meta}
}

func (s BlockState) Equals(other BlockState) bool {
	return s.ID == other.ID && s.Meta == other.Meta
}

func (s BlockState) GetBlockID() uint8 {
	return s.ID
}

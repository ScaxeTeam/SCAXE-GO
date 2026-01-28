package block

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/item"
)

type Block struct {
	ID   uint8
	Meta uint8
}

func NewBlock(id, meta uint8) Block {
	return Block{
		ID:   id,
		Meta: meta,
	}
}

func (b Block) GetID() int {
	return int(b.ID)
}

func (b Block) GetDamage() int {
	return int(b.Meta)
}

func (b Block) String() string {
	return fmt.Sprintf("Block{ID: %d, Meta: %d}", b.ID, b.Meta)
}

func (b Block) AsItem() item.Item {
	return item.NewItem(int(b.ID), int(b.Meta), 1)
}

func (b Block) ToBlockState() BlockState {
	return NewBlockState(b.ID, b.Meta)
}

func FromBlockState(bs BlockState) Block {
	return Block{ID: bs.ID, Meta: bs.Meta}
}

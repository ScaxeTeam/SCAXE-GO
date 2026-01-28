package nbt

func (t *EndTag) Clone() Tag {
	return NewEndTag()
}

func (t *ByteTag) Clone() Tag {
	return NewByteTag(t.name, t.value)
}

func (t *ShortTag) Clone() Tag {
	return NewShortTag(t.name, t.value)
}

func (t *IntTag) Clone() Tag {
	return NewIntTag(t.name, t.value)
}

func (t *LongTag) Clone() Tag {
	return NewLongTag(t.name, t.value)
}

func (t *FloatTag) Clone() Tag {
	return NewFloatTag(t.name, t.value)
}

func (t *DoubleTag) Clone() Tag {
	return NewDoubleTag(t.name, t.value)
}

func (t *ByteArrayTag) Clone() Tag {
	data := make([]byte, len(t.value))
	copy(data, t.value)
	return NewByteArrayTag(t.name, data)
}

func (t *StringTag) Clone() Tag {
	return NewStringTag(t.name, t.value)
}

func (t *ListTag) Clone() Tag {
	clone := NewListTag(t.name, t.tagType)
	for _, tag := range t.value {
		clone.Add(tag.Clone())
	}
	return clone
}

func (t *IntArrayTag) Clone() Tag {
	data := make([]int32, len(t.value))
	copy(data, t.value)
	return NewIntArrayTag(t.name, data)
}

func (t *CompoundTag) Clone() Tag {
	clone := NewCompoundTag(t.name)
	clone.order = make([]string, len(t.order))
	copy(clone.order, t.order)
	for name, tag := range t.value {
		clone.value[name] = tag.Clone()
	}
	return clone
}

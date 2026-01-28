package nbt

import (
	"bytes"
	"testing"
)

func TestByteTag(t *testing.T) {
	tag := NewByteTag("test", 42)
	if tag.Name() != "test" {
		t.Errorf("expected name 'test', got %q", tag.Name())
	}
	if tag.Value().(int8) != 42 {
		t.Errorf("expected value 42, got %v", tag.Value())
	}
	if tag.Type() != TagByte {
		t.Errorf("expected type %d, got %d", TagByte, tag.Type())
	}
}

func TestShortTag(t *testing.T) {
	tag := NewShortTag("test", 1000)
	if tag.Value().(int16) != 1000 {
		t.Errorf("expected value 1000, got %v", tag.Value())
	}
}

func TestIntTag(t *testing.T) {
	tag := NewIntTag("test", 100000)
	if tag.Value().(int32) != 100000 {
		t.Errorf("expected value 100000, got %v", tag.Value())
	}
}

func TestLongTag(t *testing.T) {
	tag := NewLongTag("test", 9999999999)
	if tag.Value().(int64) != 9999999999 {
		t.Errorf("expected value 9999999999, got %v", tag.Value())
	}
}

func TestFloatTag(t *testing.T) {
	tag := NewFloatTag("test", 3.14)
	val := tag.Value().(float32)
	if val < 3.13 || val > 3.15 {
		t.Errorf("expected value ~3.14, got %v", val)
	}
}

func TestDoubleTag(t *testing.T) {
	tag := NewDoubleTag("test", 3.141592653589793)
	val := tag.Value().(float64)
	if val < 3.14 || val > 3.15 {
		t.Errorf("expected value ~3.14, got %v", val)
	}
}

func TestStringTag(t *testing.T) {
	tag := NewStringTag("name", "Hello, NBT!")
	if tag.Value().(string) != "Hello, NBT!" {
		t.Errorf("expected 'Hello, NBT!', got %v", tag.Value())
	}
}

func TestByteArrayTag(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	tag := NewByteArrayTag("data", data)
	result := tag.Value().([]byte)
	if len(result) != 5 {
		t.Errorf("expected 5 bytes, got %d", len(result))
	}
}

func TestIntArrayTag(t *testing.T) {
	data := []int32{100, 200, 300}
	tag := NewIntArrayTag("ints", data)
	result := tag.Value().([]int32)
	if len(result) != 3 {
		t.Errorf("expected 3 ints, got %d", len(result))
	}
}

func TestCompoundTag(t *testing.T) {
	compound := NewCompoundTag("root")
	compound.Set(NewByteTag("byte", 1))
	compound.Set(NewStringTag("name", "test"))
	compound.Set(NewIntTag("count", 42))

	if compound.Count() != 3 {
		t.Errorf("expected 3 entries, got %d", compound.Count())
	}
	if compound.GetByte("byte") != 1 {
		t.Errorf("expected byte 1, got %d", compound.GetByte("byte"))
	}
	if compound.GetString("name") != "test" {
		t.Errorf("expected string 'test', got %s", compound.GetString("name"))
	}
	if compound.GetInt("count") != 42 {
		t.Errorf("expected int 42, got %d", compound.GetInt("count"))
	}
}

func TestListTag(t *testing.T) {
	list := NewListTag("items", TagInt)
	list.Add(NewIntTag("", 1))
	list.Add(NewIntTag("", 2))
	list.Add(NewIntTag("", 3))

	if list.Len() != 3 {
		t.Errorf("expected 3 items, got %d", list.Len())
	}
	if list.TagType() != TagInt {
		t.Errorf("expected tag type %d, got %d", TagInt, list.TagType())
	}
}

func TestReadWrite(t *testing.T) {

	root := NewCompoundTag("root")
	root.Set(NewByteTag("byte", 127))
	root.Set(NewShortTag("short", 32767))
	root.Set(NewIntTag("int", 2147483647))
	root.Set(NewLongTag("long", 9223372036854775807))
	root.Set(NewFloatTag("float", 3.14))
	root.Set(NewDoubleTag("double", 3.141592653589793))
	root.Set(NewStringTag("string", "Hello, World!"))
	root.Set(NewByteArrayTag("bytes", []byte{1, 2, 3, 4, 5}))
	root.Set(NewIntArrayTag("ints", []int32{100, 200, 300}))

	buf := new(bytes.Buffer)
	writer := NewWriter(buf, LittleEndian)
	if err := writer.WriteTag(root); err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	reader := NewReader(bytes.NewReader(buf.Bytes()), LittleEndian)
	tag, err := reader.ReadTag()
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	readRoot, ok := tag.(*CompoundTag)
	if !ok {
		t.Fatalf("expected CompoundTag, got %T", tag)
	}

	if readRoot.GetByte("byte") != 127 {
		t.Errorf("byte mismatch: expected 127, got %d", readRoot.GetByte("byte"))
	}
	if readRoot.GetShort("short") != 32767 {
		t.Errorf("short mismatch: expected 32767, got %d", readRoot.GetShort("short"))
	}
	if readRoot.GetInt("int") != 2147483647 {
		t.Errorf("int mismatch: expected 2147483647, got %d", readRoot.GetInt("int"))
	}
	if readRoot.GetLong("long") != 9223372036854775807 {
		t.Errorf("long mismatch: expected 9223372036854775807, got %d", readRoot.GetLong("long"))
	}
	if readRoot.GetString("string") != "Hello, World!" {
		t.Errorf("string mismatch: expected 'Hello, World!', got %s", readRoot.GetString("string"))
	}

	bytes := readRoot.GetByteArray("bytes")
	if len(bytes) != 5 || bytes[0] != 1 || bytes[4] != 5 {
		t.Errorf("byte array mismatch: %v", bytes)
	}

	ints := readRoot.GetIntArray("ints")
	if len(ints) != 3 || ints[0] != 100 || ints[2] != 300 {
		t.Errorf("int array mismatch: %v", ints)
	}
}

func TestNestedCompound(t *testing.T) {
	root := NewCompoundTag("root")

	pos := NewCompoundTag("pos")
	pos.Set(NewIntTag("x", 100))
	pos.Set(NewIntTag("y", 64))
	pos.Set(NewIntTag("z", 200))
	root.Set(pos)

	root.Set(NewStringTag("name", "Player"))

	buf := new(bytes.Buffer)
	writer := NewWriter(buf, LittleEndian)
	if err := writer.WriteTag(root); err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	reader := NewReader(bytes.NewReader(buf.Bytes()), LittleEndian)
	tag, err := reader.ReadTag()
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	readRoot := tag.(*CompoundTag)
	readPos := readRoot.GetCompound("pos")
	if readPos == nil {
		t.Fatal("nested compound 'pos' not found")
	}

	if readPos.GetInt("x") != 100 {
		t.Errorf("x mismatch: expected 100, got %d", readPos.GetInt("x"))
	}
	if readPos.GetInt("y") != 64 {
		t.Errorf("y mismatch: expected 64, got %d", readPos.GetInt("y"))
	}
	if readPos.GetInt("z") != 200 {
		t.Errorf("z mismatch: expected 200, got %d", readPos.GetInt("z"))
	}
}

func TestListOfCompounds(t *testing.T) {
	root := NewCompoundTag("root")

	items := NewListTag("items", TagCompound)

	item1 := NewCompoundTag("")
	item1.Set(NewStringTag("id", "minecraft:diamond"))
	item1.Set(NewByteTag("count", 64))
	items.Add(item1)

	item2 := NewCompoundTag("")
	item2.Set(NewStringTag("id", "minecraft:iron_ingot"))
	item2.Set(NewByteTag("count", 32))
	items.Add(item2)

	root.Set(items)

	buf := new(bytes.Buffer)
	writer := NewWriter(buf, LittleEndian)
	if err := writer.WriteTag(root); err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	reader := NewReader(bytes.NewReader(buf.Bytes()), LittleEndian)
	tag, err := reader.ReadTag()
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	readRoot := tag.(*CompoundTag)
	readItems := readRoot.GetList("items")
	if readItems == nil {
		t.Fatal("list 'items' not found")
	}

	if readItems.Len() != 2 {
		t.Errorf("expected 2 items, got %d", readItems.Len())
	}

	firstItem := readItems.Get(0).(*CompoundTag)
	if firstItem.GetString("id") != "minecraft:diamond" {
		t.Errorf("first item id mismatch: %s", firstItem.GetString("id"))
	}
	if firstItem.GetByte("count") != 64 {
		t.Errorf("first item count mismatch: %d", firstItem.GetByte("count"))
	}
}

func TestBigEndian(t *testing.T) {
	root := NewCompoundTag("test")
	root.Set(NewIntTag("value", 12345678))

	buf := new(bytes.Buffer)
	writer := NewWriter(buf, BigEndian)
	if err := writer.WriteTag(root); err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	reader := NewReader(bytes.NewReader(buf.Bytes()), BigEndian)
	tag, err := reader.ReadTag()
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	readRoot := tag.(*CompoundTag)
	if readRoot.GetInt("value") != 12345678 {
		t.Errorf("value mismatch: expected 12345678, got %d", readRoot.GetInt("value"))
	}
}

func TestNBTClass(t *testing.T) {
	nbt := New(LittleEndian)

	root := NewCompoundTag("")
	root.Set(NewStringTag("level", "world"))
	root.Set(NewIntTag("seed", 12345))
	nbt.SetData(root)

	data, err := nbt.Write()
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	nbt2 := New(LittleEndian)
	if err := nbt2.Read(data); err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if nbt2.GetData().GetString("level") != "world" {
		t.Errorf("level mismatch")
	}
	if nbt2.GetData().GetInt("seed") != 12345 {
		t.Errorf("seed mismatch")
	}
}

func TestCompression(t *testing.T) {
	nbt := New(LittleEndian)

	root := NewCompoundTag("")
	root.Set(NewStringTag("data", "This is some test data for compression"))
	root.Set(NewByteArrayTag("bytes", make([]byte, 1000)))
	nbt.SetData(root)

	gzipData, err := nbt.WriteGzip()
	if err != nil {
		t.Fatalf("failed to write gzip: %v", err)
	}

	nbt2 := New(LittleEndian)
	if err := nbt2.ReadCompressed(gzipData); err != nil {
		t.Fatalf("failed to read gzip: %v", err)
	}
	if nbt2.GetData().GetString("data") != "This is some test data for compression" {
		t.Errorf("gzip data mismatch")
	}

	zlibData, err := nbt.WriteZlib()
	if err != nil {
		t.Fatalf("failed to write zlib: %v", err)
	}

	nbt3 := New(LittleEndian)
	if err := nbt3.ReadCompressed(zlibData); err != nil {
		t.Fatalf("failed to read zlib: %v", err)
	}
	if nbt3.GetData().GetString("data") != "This is some test data for compression" {
		t.Errorf("zlib data mismatch")
	}
}

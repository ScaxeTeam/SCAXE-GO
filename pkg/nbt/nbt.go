package nbt

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"
)

const (
	TagEnd       byte = 0
	TagByte      byte = 1
	TagShort     byte = 2
	TagInt       byte = 3
	TagLong      byte = 4
	TagFloat     byte = 5
	TagDouble    byte = 6
	TagByteArray byte = 7
	TagString    byte = 8
	TagList      byte = 9
	TagCompound  byte = 10
	TagIntArray  byte = 11
)

type Endianness int

const (
	LittleEndian Endianness = iota
	BigEndian
)

type Tag interface {
	Type() byte
	Name() string
	SetName(name string)
	Value() interface{}
	SetValue(value interface{}) error
	Read(r *Reader) error
	Write(w *Writer) error
	String() string
	Clone() Tag
}

type baseTag struct {
	name string
}

func (t *baseTag) Name() string {
	return t.name
}

func (t *baseTag) SetName(name string) {
	t.name = name
}

type EndTag struct {
	baseTag
}

func NewEndTag() *EndTag {
	return &EndTag{}
}

func (t *EndTag) Type() byte                   { return TagEnd }
func (t *EndTag) Value() interface{}           { return nil }
func (t *EndTag) SetValue(v interface{}) error { return nil }
func (t *EndTag) Read(r *Reader) error         { return nil }
func (t *EndTag) Write(w *Writer) error        { return nil }
func (t *EndTag) String() string               { return "EndTag" }

type ByteTag struct {
	baseTag
	value int8
}

func NewByteTag(name string, value int8) *ByteTag {
	return &ByteTag{baseTag: baseTag{name: name}, value: value}
}

func (t *ByteTag) Type() byte         { return TagByte }
func (t *ByteTag) Value() interface{} { return t.value }
func (t *ByteTag) SetValue(v interface{}) error {
	switch val := v.(type) {
	case int8:
		t.value = val
	case int:
		t.value = int8(val)
	case byte:
		t.value = int8(val)
	default:
		return fmt.Errorf("invalid type for ByteTag: %T", v)
	}
	return nil
}
func (t *ByteTag) Read(r *Reader) error {
	v, err := r.ReadByte()
	if err != nil {
		return err
	}
	t.value = int8(v)
	return nil
}
func (t *ByteTag) Write(w *Writer) error {
	return w.WriteByte(byte(t.value))
}
func (t *ByteTag) String() string {
	return fmt.Sprintf("ByteTag(%q): %d", t.name, t.value)
}

type ShortTag struct {
	baseTag
	value int16
}

func NewShortTag(name string, value int16) *ShortTag {
	return &ShortTag{baseTag: baseTag{name: name}, value: value}
}

func (t *ShortTag) Type() byte         { return TagShort }
func (t *ShortTag) Value() interface{} { return t.value }
func (t *ShortTag) SetValue(v interface{}) error {
	switch val := v.(type) {
	case int16:
		t.value = val
	case int:
		t.value = int16(val)
	default:
		return fmt.Errorf("invalid type for ShortTag: %T", v)
	}
	return nil
}
func (t *ShortTag) Read(r *Reader) error {
	v, err := r.ReadInt16()
	if err != nil {
		return err
	}
	t.value = v
	return nil
}
func (t *ShortTag) Write(w *Writer) error {
	return w.WriteInt16(t.value)
}
func (t *ShortTag) String() string {
	return fmt.Sprintf("ShortTag(%q): %d", t.name, t.value)
}

type IntTag struct {
	baseTag
	value int32
}

func NewIntTag(name string, value int32) *IntTag {
	return &IntTag{baseTag: baseTag{name: name}, value: value}
}

func (t *IntTag) Type() byte         { return TagInt }
func (t *IntTag) Value() interface{} { return t.value }
func (t *IntTag) SetValue(v interface{}) error {
	switch val := v.(type) {
	case int32:
		t.value = val
	case int:
		t.value = int32(val)
	default:
		return fmt.Errorf("invalid type for IntTag: %T", v)
	}
	return nil
}
func (t *IntTag) Read(r *Reader) error {
	v, err := r.ReadInt32()
	if err != nil {
		return err
	}
	t.value = v
	return nil
}
func (t *IntTag) Write(w *Writer) error {
	return w.WriteInt32(t.value)
}
func (t *IntTag) String() string {
	return fmt.Sprintf("IntTag(%q): %d", t.name, t.value)
}

type LongTag struct {
	baseTag
	value int64
}

func NewLongTag(name string, value int64) *LongTag {
	return &LongTag{baseTag: baseTag{name: name}, value: value}
}

func (t *LongTag) Type() byte         { return TagLong }
func (t *LongTag) Value() interface{} { return t.value }
func (t *LongTag) SetValue(v interface{}) error {
	switch val := v.(type) {
	case int64:
		t.value = val
	case int:
		t.value = int64(val)
	default:
		return fmt.Errorf("invalid type for LongTag: %T", v)
	}
	return nil
}
func (t *LongTag) Read(r *Reader) error {
	v, err := r.ReadInt64()
	if err != nil {
		return err
	}
	t.value = v
	return nil
}
func (t *LongTag) Write(w *Writer) error {
	return w.WriteInt64(t.value)
}
func (t *LongTag) String() string {
	return fmt.Sprintf("LongTag(%q): %d", t.name, t.value)
}

type FloatTag struct {
	baseTag
	value float32
}

func NewFloatTag(name string, value float32) *FloatTag {
	return &FloatTag{baseTag: baseTag{name: name}, value: value}
}

func (t *FloatTag) Type() byte         { return TagFloat }
func (t *FloatTag) Value() interface{} { return t.value }
func (t *FloatTag) SetValue(v interface{}) error {
	switch val := v.(type) {
	case float32:
		t.value = val
	case float64:
		t.value = float32(val)
	default:
		return fmt.Errorf("invalid type for FloatTag: %T", v)
	}
	return nil
}
func (t *FloatTag) Read(r *Reader) error {
	v, err := r.ReadFloat32()
	if err != nil {
		return err
	}
	t.value = v
	return nil
}
func (t *FloatTag) Write(w *Writer) error {
	return w.WriteFloat32(t.value)
}
func (t *FloatTag) String() string {
	return fmt.Sprintf("FloatTag(%q): %f", t.name, t.value)
}

type DoubleTag struct {
	baseTag
	value float64
}

func NewDoubleTag(name string, value float64) *DoubleTag {
	return &DoubleTag{baseTag: baseTag{name: name}, value: value}
}

func (t *DoubleTag) Type() byte         { return TagDouble }
func (t *DoubleTag) Value() interface{} { return t.value }
func (t *DoubleTag) SetValue(v interface{}) error {
	switch val := v.(type) {
	case float64:
		t.value = val
	case float32:
		t.value = float64(val)
	default:
		return fmt.Errorf("invalid type for DoubleTag: %T", v)
	}
	return nil
}
func (t *DoubleTag) Read(r *Reader) error {
	v, err := r.ReadFloat64()
	if err != nil {
		return err
	}
	t.value = v
	return nil
}
func (t *DoubleTag) Write(w *Writer) error {
	return w.WriteFloat64(t.value)
}
func (t *DoubleTag) String() string {
	return fmt.Sprintf("DoubleTag(%q): %f", t.name, t.value)
}

type ByteArrayTag struct {
	baseTag
	value []byte
}

func NewByteArrayTag(name string, value []byte) *ByteArrayTag {
	return &ByteArrayTag{baseTag: baseTag{name: name}, value: value}
}

func (t *ByteArrayTag) Type() byte         { return TagByteArray }
func (t *ByteArrayTag) Value() interface{} { return t.value }
func (t *ByteArrayTag) SetValue(v interface{}) error {
	if val, ok := v.([]byte); ok {
		t.value = val
		return nil
	}
	return fmt.Errorf("invalid type for ByteArrayTag: %T", v)
}
func (t *ByteArrayTag) Read(r *Reader) error {
	length, err := r.ReadInt32()
	if err != nil {
		return err
	}
	t.value = make([]byte, length)
	_, err = io.ReadFull(r.r, t.value)
	return err
}
func (t *ByteArrayTag) Write(w *Writer) error {
	if err := w.WriteInt32(int32(len(t.value))); err != nil {
		return err
	}
	_, err := w.w.Write(t.value)
	return err
}
func (t *ByteArrayTag) String() string {
	return fmt.Sprintf("ByteArrayTag(%q): [%d bytes]", t.name, len(t.value))
}

type StringTag struct {
	baseTag
	value string
}

func NewStringTag(name string, value string) *StringTag {
	return &StringTag{baseTag: baseTag{name: name}, value: value}
}

func (t *StringTag) Type() byte         { return TagString }
func (t *StringTag) Value() interface{} { return t.value }
func (t *StringTag) SetValue(v interface{}) error {
	if val, ok := v.(string); ok {
		t.value = val
		return nil
	}
	return fmt.Errorf("invalid type for StringTag: %T", v)
}
func (t *StringTag) Read(r *Reader) error {
	v, err := r.ReadString()
	if err != nil {
		return err
	}
	t.value = v
	return nil
}
func (t *StringTag) Write(w *Writer) error {
	return w.WriteString(t.value)
}
func (t *StringTag) String() string {
	return fmt.Sprintf("StringTag(%q): %q", t.name, t.value)
}

type ListTag struct {
	baseTag
	tagType byte
	value   []Tag
}

func NewListTag(name string, tagType byte) *ListTag {
	return &ListTag{baseTag: baseTag{name: name}, tagType: tagType, value: []Tag{}}
}

func (t *ListTag) Type() byte              { return TagList }
func (t *ListTag) Value() interface{}      { return t.value }
func (t *ListTag) TagType() byte           { return t.tagType }
func (t *ListTag) SetTagType(tagType byte) { t.tagType = tagType }
func (t *ListTag) Len() int                { return len(t.value) }

func (t *ListTag) SetValue(v interface{}) error {
	if val, ok := v.([]Tag); ok {
		t.value = val
		return nil
	}
	return fmt.Errorf("invalid type for ListTag: %T", v)
}

func (t *ListTag) Add(tag Tag) {
	t.value = append(t.value, tag)
}

func (t *ListTag) Get(index int) Tag {
	if index < 0 || index >= len(t.value) {
		return nil
	}
	return t.value[index]
}

func (t *ListTag) Read(r *Reader) error {
	tagType, err := r.ReadByte()
	if err != nil {
		return err
	}
	t.tagType = tagType

	length, err := r.ReadInt32()
	if err != nil {
		return err
	}

	t.value = make([]Tag, 0, length)
	for i := int32(0); i < length; i++ {
		tag := CreateTag(t.tagType)
		if tag == nil {
			return fmt.Errorf("unknown tag type: %d", t.tagType)
		}
		if err := tag.Read(r); err != nil {
			return err
		}
		t.value = append(t.value, tag)
	}
	return nil
}

func (t *ListTag) Write(w *Writer) error {

	if t.tagType == TagEnd && len(t.value) > 0 {
		t.tagType = t.value[0].Type()
	}

	if err := w.WriteByte(t.tagType); err != nil {
		return err
	}
	if err := w.WriteInt32(int32(len(t.value))); err != nil {
		return err
	}
	for _, tag := range t.value {
		if err := tag.Write(w); err != nil {
			return err
		}
	}
	return nil
}

func (t *ListTag) String() string {
	return fmt.Sprintf("ListTag(%q): [%d entries of type %d]", t.name, len(t.value), t.tagType)
}

type IntArrayTag struct {
	baseTag
	value []int32
}

func NewIntArrayTag(name string, value []int32) *IntArrayTag {
	return &IntArrayTag{baseTag: baseTag{name: name}, value: value}
}

func (t *IntArrayTag) Type() byte         { return TagIntArray }
func (t *IntArrayTag) Value() interface{} { return t.value }
func (t *IntArrayTag) SetValue(v interface{}) error {
	if val, ok := v.([]int32); ok {
		t.value = val
		return nil
	}
	return fmt.Errorf("invalid type for IntArrayTag: %T", v)
}
func (t *IntArrayTag) Read(r *Reader) error {
	length, err := r.ReadInt32()
	if err != nil {
		return err
	}
	t.value = make([]int32, length)
	for i := int32(0); i < length; i++ {
		v, err := r.ReadInt32()
		if err != nil {
			return err
		}
		t.value[i] = v
	}
	return nil
}
func (t *IntArrayTag) Write(w *Writer) error {
	if err := w.WriteInt32(int32(len(t.value))); err != nil {
		return err
	}
	for _, v := range t.value {
		if err := w.WriteInt32(v); err != nil {
			return err
		}
	}
	return nil
}
func (t *IntArrayTag) String() string {
	return fmt.Sprintf("IntArrayTag(%q): [%d ints]", t.name, len(t.value))
}

type CompoundTag struct {
	baseTag
	value map[string]Tag
	order []string
}

func NewCompoundTag(name string) *CompoundTag {
	return &CompoundTag{
		baseTag: baseTag{name: name},
		value:   make(map[string]Tag),
		order:   []string{},
	}
}

func (t *CompoundTag) Type() byte         { return TagCompound }
func (t *CompoundTag) Value() interface{} { return t.value }
func (t *CompoundTag) SetValue(v interface{}) error {
	if val, ok := v.(map[string]Tag); ok {
		t.value = val
		return nil
	}
	return fmt.Errorf("invalid type for CompoundTag: %T", v)
}

func (t *CompoundTag) Set(tag Tag) {
	name := tag.Name()
	if _, exists := t.value[name]; !exists {
		t.order = append(t.order, name)
	}
	t.value[name] = tag
}

func (t *CompoundTag) Get(name string) Tag {
	return t.value[name]
}

func (t *CompoundTag) Has(name string) bool {
	_, exists := t.value[name]
	return exists
}

func (t *CompoundTag) Remove(name string) {
	delete(t.value, name)
	for i, n := range t.order {
		if n == name {
			t.order = append(t.order[:i], t.order[i+1:]...)
			break
		}
	}
}

func (t *CompoundTag) Count() int {
	return len(t.value)
}

func (t *CompoundTag) Tags() []Tag {
	tags := make([]Tag, 0, len(t.value))
	for _, name := range t.order {
		if tag, exists := t.value[name]; exists {
			tags = append(tags, tag)
		}
	}
	return tags
}

func (t *CompoundTag) GetByte(name string) int8 {
	if tag, ok := t.value[name].(*ByteTag); ok {
		return tag.value
	}
	return 0
}

func (t *CompoundTag) GetShort(name string) int16 {
	if tag, ok := t.value[name].(*ShortTag); ok {
		return tag.value
	}
	return 0
}

func (t *CompoundTag) GetInt(name string) int32 {
	if tag, ok := t.value[name].(*IntTag); ok {
		return tag.value
	}
	return 0
}

func (t *CompoundTag) GetLong(name string) int64 {
	if tag, ok := t.value[name].(*LongTag); ok {
		return tag.value
	}
	return 0
}

func (t *CompoundTag) GetFloat(name string) float32 {
	if tag, ok := t.value[name].(*FloatTag); ok {
		return tag.value
	}
	return 0
}

func (t *CompoundTag) GetDouble(name string) float64 {
	if tag, ok := t.value[name].(*DoubleTag); ok {
		return tag.value
	}
	return 0
}

func (t *CompoundTag) GetString(name string) string {
	if tag, ok := t.value[name].(*StringTag); ok {
		return tag.value
	}
	return ""
}

func (t *CompoundTag) GetByteArray(name string) []byte {
	if tag, ok := t.value[name].(*ByteArrayTag); ok {
		return tag.value
	}
	return nil
}

func (t *CompoundTag) GetIntArray(name string) []int32 {
	if tag, ok := t.value[name].(*IntArrayTag); ok {
		return tag.value
	}
	return nil
}

func (t *CompoundTag) GetList(name string) *ListTag {
	if tag, ok := t.value[name].(*ListTag); ok {
		return tag
	}
	return nil
}

func (t *CompoundTag) GetCompound(name string) *CompoundTag {
	if tag, ok := t.value[name].(*CompoundTag); ok {
		return tag
	}
	return nil
}

func (t *CompoundTag) Read(r *Reader) error {
	t.value = make(map[string]Tag)
	t.order = []string{}

	for {
		tag, err := r.ReadTag()
		if err != nil {
			return err
		}
		if tag.Type() == TagEnd {
			break
		}
		name := tag.Name()
		t.value[name] = tag
		t.order = append(t.order, name)
	}
	return nil
}

func (t *CompoundTag) Write(w *Writer) error {
	for _, name := range t.order {
		tag := t.value[name]
		if err := w.WriteTag(tag); err != nil {
			return err
		}
	}
	return w.WriteTag(NewEndTag())
}

func (t *CompoundTag) String() string {
	return fmt.Sprintf("CompoundTag(%q): {%d entries}", t.name, len(t.value))
}

func CreateTag(tagType byte) Tag {
	switch tagType {
	case TagEnd:
		return NewEndTag()
	case TagByte:
		return &ByteTag{}
	case TagShort:
		return &ShortTag{}
	case TagInt:
		return &IntTag{}
	case TagLong:
		return &LongTag{}
	case TagFloat:
		return &FloatTag{}
	case TagDouble:
		return &DoubleTag{}
	case TagByteArray:
		return &ByteArrayTag{}
	case TagString:
		return &StringTag{}
	case TagList:
		return &ListTag{value: []Tag{}}
	case TagCompound:
		return NewCompoundTag("")
	case TagIntArray:
		return &IntArrayTag{}
	default:
		return nil
	}
}

type Reader struct {
	r         io.Reader
	endian    Endianness
	byteOrder binary.ByteOrder
}

func NewReader(r io.Reader, endian Endianness) *Reader {
	var order binary.ByteOrder = binary.LittleEndian
	if endian == BigEndian {
		order = binary.BigEndian
	}
	return &Reader{r: r, endian: endian, byteOrder: order}
}

func (r *Reader) ReadByte() (byte, error) {
	var b [1]byte
	_, err := io.ReadFull(r.r, b[:])
	return b[0], err
}

func (r *Reader) ReadInt16() (int16, error) {
	var b [2]byte
	_, err := io.ReadFull(r.r, b[:])
	if err != nil {
		return 0, err
	}
	return int16(r.byteOrder.Uint16(b[:])), nil
}

func (r *Reader) ReadInt32() (int32, error) {
	var b [4]byte
	_, err := io.ReadFull(r.r, b[:])
	if err != nil {
		return 0, err
	}
	return int32(r.byteOrder.Uint32(b[:])), nil
}

func (r *Reader) ReadInt64() (int64, error) {
	var b [8]byte
	_, err := io.ReadFull(r.r, b[:])
	if err != nil {
		return 0, err
	}
	return int64(r.byteOrder.Uint64(b[:])), nil
}

func (r *Reader) ReadFloat32() (float32, error) {
	var b [4]byte
	_, err := io.ReadFull(r.r, b[:])
	if err != nil {
		return 0, err
	}
	bits := r.byteOrder.Uint32(b[:])
	return float32frombits(bits), nil
}

func (r *Reader) ReadFloat64() (float64, error) {
	var b [8]byte
	_, err := io.ReadFull(r.r, b[:])
	if err != nil {
		return 0, err
	}
	bits := r.byteOrder.Uint64(b[:])
	return float64frombits(bits), nil
}

func (r *Reader) ReadString() (string, error) {
	length, err := r.ReadInt16()
	if err != nil {
		return "", err
	}
	if length < 0 {
		return "", fmt.Errorf("negative string length: %d", length)
	}
	data := make([]byte, length)
	_, err = io.ReadFull(r.r, data)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *Reader) ReadTag() (Tag, error) {
	tagType, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if tagType == TagEnd {
		return NewEndTag(), nil
	}

	name, err := r.ReadString()
	if err != nil {
		return nil, err
	}

	tag := CreateTag(tagType)
	if tag == nil {
		return nil, fmt.Errorf("unknown tag type: %d", tagType)
	}

	tag.SetName(name)
	if err := tag.Read(r); err != nil {
		return nil, err
	}

	return tag, nil
}

type Writer struct {
	w         io.Writer
	endian    Endianness
	byteOrder binary.ByteOrder
}

func NewWriter(w io.Writer, endian Endianness) *Writer {
	var order binary.ByteOrder = binary.LittleEndian
	if endian == BigEndian {
		order = binary.BigEndian
	}
	return &Writer{w: w, endian: endian, byteOrder: order}
}

func (w *Writer) WriteByte(b byte) error {
	_, err := w.w.Write([]byte{b})
	return err
}

func (w *Writer) WriteInt16(v int16) error {
	var b [2]byte
	w.byteOrder.PutUint16(b[:], uint16(v))
	_, err := w.w.Write(b[:])
	return err
}

func (w *Writer) WriteInt32(v int32) error {
	var b [4]byte
	w.byteOrder.PutUint32(b[:], uint32(v))
	_, err := w.w.Write(b[:])
	return err
}

func (w *Writer) WriteInt64(v int64) error {
	var b [8]byte
	w.byteOrder.PutUint64(b[:], uint64(v))
	_, err := w.w.Write(b[:])
	return err
}

func (w *Writer) WriteFloat32(v float32) error {
	var b [4]byte
	w.byteOrder.PutUint32(b[:], float32bits(v))
	_, err := w.w.Write(b[:])
	return err
}

func (w *Writer) WriteFloat64(v float64) error {
	var b [8]byte
	w.byteOrder.PutUint64(b[:], float64bits(v))
	_, err := w.w.Write(b[:])
	return err
}

func (w *Writer) WriteString(s string) error {
	if len(s) > 32767 {
		return fmt.Errorf("string too long: %d bytes", len(s))
	}
	if err := w.WriteInt16(int16(len(s))); err != nil {
		return err
	}
	_, err := w.w.Write([]byte(s))
	return err
}

func (w *Writer) WriteTag(tag Tag) error {
	if err := w.WriteByte(tag.Type()); err != nil {
		return err
	}

	if tag.Type() == TagEnd {
		return nil
	}

	if err := w.WriteString(tag.Name()); err != nil {
		return err
	}

	return tag.Write(w)
}

func float32bits(f float32) uint32 {
	return *(*uint32)(unsafe.Pointer(&f))
}

func float32frombits(b uint32) float32 {
	return *(*float32)(unsafe.Pointer(&b))
}

func float64bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}

func float64frombits(b uint64) float64 {
	return *(*float64)(unsafe.Pointer(&b))
}

type NBT struct {
	endian Endianness
	data   *CompoundTag
}

func New(endian Endianness) *NBT {
	return &NBT{endian: endian}
}

func (n *NBT) Read(data []byte) error {
	r := NewReader(bytes.NewReader(data), n.endian)
	tag, err := r.ReadTag()
	if err != nil {
		return err
	}
	compound, ok := tag.(*CompoundTag)
	if !ok {
		return fmt.Errorf("root tag must be CompoundTag, got %T", tag)
	}
	n.data = compound
	return nil
}

func (n *NBT) ReadCompressed(data []byte) error {

	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err == nil {
		defer gr.Close()
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, gr)
		if err == nil {
			return n.Read(buf.Bytes())
		}
	}

	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decompress: not gzip or zlib")
	}
	defer zr.Close()
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, zr); err != nil {
		return err
	}
	return n.Read(buf.Bytes())
}

func (n *NBT) Write() ([]byte, error) {
	if n.data == nil {
		return nil, fmt.Errorf("no data to write")
	}
	buf := new(bytes.Buffer)
	w := NewWriter(buf, n.endian)
	if err := w.WriteTag(n.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (n *NBT) WriteGzip() ([]byte, error) {
	data, err := n.Write()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	gw := gzip.NewWriter(buf)
	if _, err := gw.Write(data); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (n *NBT) WriteZlib() ([]byte, error) {
	data, err := n.Write()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	zw := zlib.NewWriter(buf)
	if _, err := zw.Write(data); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (n *NBT) GetData() *CompoundTag {
	return n.data
}

func (n *NBT) SetData(data *CompoundTag) {
	n.data = data
}

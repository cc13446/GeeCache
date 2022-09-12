package GeeCache

// A ValueView holds an immutable view of bytes.
type ValueView struct {
	data []byte
}

// Len return the view's length
func (v ValueView) Len() int {
	return len(v.data)
}

// ByteSlice returns a copy of the data as a byte slice
func (v ValueView) ByteSlice() []byte {
	return cloneBytes(v.data)
}

// String returns the data as string
func (v ValueView) String() string {
	return string(v.data)
}

func cloneBytes(data []byte) []byte {
	temp := make([]byte, len(data))
	copy(temp, data)
	return temp
}

// A KeyView holds an immutable view of bytes.
type KeyView struct {
	key string
}

// Len return the view's length
func (k KeyView) Len() int {
	return len(k.key)
}

// String returns the data as string
func (k KeyView) String() string {
	return k.key
}

// FromString returns KeyView from string
func FromString(s string) KeyView {
	return KeyView{
		key: s,
	}
}

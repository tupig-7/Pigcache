package pigcache

// A ByteView holds an immutable view of bytes
// ByteView 用于缓存值
type ByteView struct {
	b []byte
}

// Len returns the view's length
// Len 返回视图的长度
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice returns a copy of the data as a byte slice
// ByteSlice 使用字节切片返回数据拷贝
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String returns the data as a string, making a copy if necessary
// String 将数据作为字符串返回，必要时制作副本
func (v ByteView) String() string {
	return string(v.b)
}

// cloneBytes returns a copy of the b as a byte slice
// cloneBytes返回传入值b的切片拷贝
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
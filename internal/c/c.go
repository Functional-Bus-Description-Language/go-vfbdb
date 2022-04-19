// The c package contains miscellaneous code common to all C targets.
package c

// WidthToCReadType returns type that is sufficient to represent data
// of given width in the C language for read functions.
func WidthToReadType(width int64) Type {
	if width > 64 {
		return ByteArray{}
	} else if width > 32 {
		return Uint64Ptr{}
	} else if width > 16 {
		return Uint32Ptr{}
	} else if width > 8 {
		return Uint16Ptr{}
	}
	return Uint8Ptr{}
}

// WidthToWriteType returns type that is sufficient to represent data
// of given width in the C language for write functions.
func WidthToWriteType(width int64) Type {
	if width > 64 {
		return ByteArray{}
	} else if width > 32 {
		return Uint64{}
	} else if width > 16 {
		return Uint32{}
	} else if width > 8 {
		return Uint16{}
	}
	return Uint8{}
}
